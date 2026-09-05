package skills

import (
	"os"
	"path/filepath"
	"scorp-agent/config"
)

// EnsureBuiltinSkills creates default foundational skills if none exist
func EnsureBuiltinSkills() {
	globalDir := config.SkillsDirPath()
	_ = os.MkdirAll(globalDir, 0755)

	builtinList := []struct {
		folder string
		md     string
	}{
		{
			folder: "docker-ops",
			md: `---
name: docker-ops
description: Expert playbook for Docker, Dockerfile, docker compose, container logs, inspection, and network troubleshooting.
emoji: 🐳
category: devops
keywords: [docker, container, compose, dockerfile, containerize]
---
# Docker Operations Playbook

When handling Docker tasks:
1. Always check existing containers and images before modifying: 'docker ps -a' or 'docker images'.
2. Use non-destructive commands first. Never run 'docker rm -f' or 'docker system prune -a' without explicit confirmation.
3. For Dockerfiles: prioritize minimal base images (Alpine or Distroless), multi-stage builds, and non-root users.
4. For Docker Compose: verify YAML syntax and port collisions before executing 'docker compose up -d'.
5. Always inspect container logs with 'docker logs --tail 50 <container>' when diagnosing failures.
`,
		},
		{
			folder: "vps-devops",
			md: `---
name: vps-devops
description: Linux server administration, systemd service management, firewall, performance tuning, and reverse proxy setup.
emoji: 🐧
category: devops
keywords: [vps, linux, systemd, service, nginx, caddy, firewall, ufw, debian, ubuntu]
---
# Linux VPS DevOps Playbook

When administering Linux VPS:
1. Inspect running services with 'systemctl status <service>' before attempting restart.
2. Check resource bottlenecks with 'free -h', 'df -h /', and 'uptime' (load averages).
3. For Nginx/Caddy configs: always test syntax ('nginx -t') before reloading.
4. Do not delete production system logs blindly; prefer logrotate or journalctl vacuuming.
5. Follow least-privilege security practices. Never chmod 777 sensitive paths.
`,
		},
		{
			folder: "git-workflow",
			md: `---
name: git-workflow
description: Professional Git workflows, rebase, clean atomic commits, conflict resolution, and branch management.
emoji: 📦
category: dev
keywords: [git, commit, branch, rebase, merge, conflict, stash, diff]
---
# Git Workflow Playbook

When managing Git operations:
1. Always run 'git status --short' and 'git diff' to review exact changes before staging.
2. Keep commits atomic, descriptive, and follow Conventional Commits (feat, fix, refactor, chore, test).
3. Do not perform forced pushes ('git push --force') on shared or master branches unless explicitly requested.
4. When resolving conflicts, inspect both HEAD and incoming changes carefully before staging.
`,
		},
		{
			folder: "golang-pro",
			md: `---
name: golang-pro
description: Go (Golang) idioms, concurrency safety, data-race prevention, zero-alloc patterns, and idiomatic error handling.
emoji: 🐹
category: dev
keywords: [golang, go, goroutine, channel, mutex, race, benchmark, pprof]
---
# Go Professional Playbook

When writing and reviewing Go code:
1. Concurrency: always run 'go test -race ./...' to guarantee zero data races. Protect shared memory with sync.Mutex or RWMutex.
2. Error Handling: wrap errors with context using fmt.Errorf("action description: %w", err). Never ignore returned errors.
3. Decoupling: prefer interfaces at consumption site, avoid circular package imports, and leverage dependency inversion callbacks.
4. Goroutine Leaks: ensure every goroutine has a termination condition via context.Context or done channel.
`,
		},
	}

	for _, b := range builtinList {
		skillDir := filepath.Join(globalDir, b.folder)
		skillFile := filepath.Join(skillDir, "SKILL.md")
		if _, err := os.Stat(skillFile); os.IsNotExist(err) {
			_ = os.MkdirAll(skillDir, 0755)
			_ = os.WriteFile(skillFile, []byte(b.md), 0644)
		}
	}
}
