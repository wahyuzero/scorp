package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var coolifyContainers = map[string]bool{
	"coolify": true, "coolify-db": true, "coolify-redis": true,
	"coolify-realtime": true, "coolify-sentinel": true, "coolify-proxy": true,
}

// DockerData holds container information.
type DockerData struct {
	Containers []ContainerInfo
	Summary    DockerSummary
}

type ContainerInfo struct {
	Name       string
	Status     string
	Health     string
	CPUPercent float64
	MemoryMB   float64
	Type       string
	Image      string
}

type DockerSummary struct {
	Total     int
	Running   int
	Unhealthy int
	Stopped   int
}

// Background stats cache
var (
	dockerStatsCache = make(map[string]containerStats)
	dockerStatsMu    sync.RWMutex
)

type containerStats struct {
	CPU float64
	Mem float64
}

// Docker API via unix socket (no SDK dependency needed)
var dockerHTTPClient *http.Client

func initDockerClient() {
	dockerHTTPClient = &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", "/var/run/docker.sock")
			},
		},
		Timeout: 10 * time.Second,
	}
}

func dockerGet(path string) ([]byte, error) {
	if dockerHTTPClient == nil {
		initDockerClient()
	}
	resp, err := dockerHTTPClient.Get("http://localhost" + path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func startDockerStatsSampler() {
	go func() {
		for {
			body, err := dockerGet("/containers/json")
			if err != nil {
				time.Sleep(30 * time.Second)
				continue
			}

			var containers []struct {
				ID    string   `json:"Id"`
				Names []string `json:"Names"`
			}
			json.Unmarshal(body, &containers)

			// Fetch stats concurrently with bounded parallelism
			type namedStats struct {
				name  string
				stats containerStats
			}
			results := make(chan namedStats, len(containers))
			sem := make(chan struct{}, 5) // max 5 concurrent

			var wg sync.WaitGroup
			for _, c := range containers {
				wg.Add(1)
				go func(id string, names []string) {
					defer wg.Done()
					sem <- struct{}{}
					defer func() { <-sem }()

					name := containerName(names)
					statsBody, err := dockerGet("/containers/" + id + "/stats?stream=false")
					if err != nil {
						return
					}

					var v struct {
						CPUStats struct {
							CPUUsage struct {
								TotalUsage uint64 `json:"total_usage"`
							} `json:"cpu_usage"`
							SystemUsage uint64 `json:"system_cpu_usage"`
							OnlineCPUs  int    `json:"online_cpus"`
						} `json:"cpu_stats"`
						PreCPUStats struct {
							CPUUsage struct {
								TotalUsage uint64 `json:"total_usage"`
							} `json:"cpu_usage"`
							SystemUsage uint64 `json:"system_cpu_usage"`
						} `json:"precpu_stats"`
						MemoryStats struct {
							Usage uint64 `json:"usage"`
						} `json:"memory_stats"`
					}
					if json.Unmarshal(statsBody, &v) != nil {
						return
					}

					cpuDelta := float64(v.CPUStats.CPUUsage.TotalUsage - v.PreCPUStats.CPUUsage.TotalUsage)
					sysDelta := float64(v.CPUStats.SystemUsage - v.PreCPUStats.SystemUsage)
					numCPUs := float64(v.CPUStats.OnlineCPUs)
					if numCPUs == 0 {
						numCPUs = 1
					}
					cpuPct := 0.0
					if sysDelta > 0 {
						cpuPct = round1((cpuDelta / sysDelta) * numCPUs * 100)
					}
					memMB := round1(float64(v.MemoryStats.Usage) / (1024 * 1024))

					results <- namedStats{name, containerStats{CPU: cpuPct, Mem: memMB}}
				}(c.ID, c.Names)
			}

			go func() {
				wg.Wait()
				close(results)
			}()

			newStats := make(map[string]containerStats)
			for r := range results {
				newStats[r.name] = r.stats
			}

			dockerStatsMu.Lock()
			dockerStatsCache = newStats
			dockerStatsMu.Unlock()

			time.Sleep(30 * time.Second)
		}
	}()
}

func collectDocker() DockerData {
	var d DockerData

	body, err := dockerGet("/containers/json?all=true")
	if err != nil {
		log.Printf("[docker] list error: %v", err)
		return d
	}

	var containers []struct {
		ID     string   `json:"Id"`
		Names  []string `json:"Names"`
		Image  string   `json:"Image"`
		State  string   `json:"State"`
		Status string   `json:"Status"`
	}
	if json.Unmarshal(body, &containers) != nil {
		return d
	}

	for _, c := range containers {
		name := containerName(c.Names)
		d.Summary.Total++
		status := c.State
		health := "N/A"

		if status == "running" {
			d.Summary.Running++
			// Get health from inspect
			inspectBody, err := dockerGet("/containers/" + c.ID + "/json")
			if err == nil {
				var inspect struct {
					State struct {
						Health *struct {
							Status string `json:"Status"`
						} `json:"Health"`
					} `json:"State"`
				}
				if json.Unmarshal(inspectBody, &inspect) == nil && inspect.State.Health != nil {
					health = inspect.State.Health.Status
				} else {
					health = "no-healthcheck"
				}
			}
			if health == "unhealthy" {
				d.Summary.Unhealthy++
			}
		} else {
			d.Summary.Stopped++
		}

		dockerStatsMu.RLock()
		cached := dockerStatsCache[name]
		dockerStatsMu.RUnlock()

		cType := "app"
		if coolifyContainers[name] {
			cType = "coolify"
		}

		image := c.Image
		if len(image) > 40 {
			image = image[:12]
		}

		d.Containers = append(d.Containers, ContainerInfo{
			Name:       name,
			Status:     status,
			Health:     health,
			CPUPercent: cached.CPU,
			MemoryMB:   cached.Mem,
			Type:       cType,
			Image:      image,
		})
	}

	return d
}

// collectDockerFallback uses docker CLI as fallback
func collectDockerFallback() DockerData {
	var d DockerData
	out, err := exec.Command("docker", "ps", "-a", "--format", "{{.Names}}\t{{.State}}\t{{.Image}}").Output()
	if err != nil {
		return d
	}
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		parts := strings.Split(line, "\t")
		if len(parts) < 3 {
			continue
		}
		d.Summary.Total++
		status := parts[1]
		if status == "running" {
			d.Summary.Running++
		} else {
			d.Summary.Stopped++
		}
		cType := "app"
		if coolifyContainers[parts[0]] {
			cType = "coolify"
		}
		d.Containers = append(d.Containers, ContainerInfo{
			Name:   parts[0],
			Status: status,
			Health: "N/A",
			Type:   cType,
			Image:  parts[2],
		})
	}
	return d
}

func containerName(names []string) string {
	if len(names) == 0 {
		return "unknown"
	}
	name := names[0]
	if len(name) > 0 && name[0] == '/' {
		name = name[1:]
	}
	return name
}
