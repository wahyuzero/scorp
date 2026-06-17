package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

// GeoIP cache (bounded)
const geoCacheMaxSize = 500

var (
	geoCache   = make(map[string]GeoInfo)
	geoCacheMu sync.RWMutex
)

type GeoInfo struct {
	Country string
	City    string
	ISP     string
	Org     string
	CC      string
	AS      string
}

var flagMap = map[string]string{
	"US": "🇺🇸", "CN": "🇨🇳", "RU": "🇷🇺", "DE": "🇩🇪", "NL": "🇳🇱",
	"FR": "🇫🇷", "GB": "🇬🇧", "JP": "🇯🇵", "KR": "🇰🇷", "BR": "🇧🇷",
	"IN": "🇮🇳", "ID": "🇮🇩", "VN": "🇻🇳", "TW": "🇹🇼", "TH": "🇹🇭",
	"SG": "🇸🇬", "AU": "🇦🇺", "CA": "🇨🇦", "IT": "🇮🇹", "ES": "🇪🇸",
	"PL": "🇵🇱", "UA": "🇺🇦", "TR": "🇹🇷", "AR": "🇦🇷", "MX": "🇲🇽",
	"SE": "🇸🇪", "HK": "🇭🇰", "MY": "🇲🇾", "PH": "🇵🇭", "BD": "🇧🇩",
	"RO": "🇷🇴", "BG": "🇧🇬", "ZA": "🇿🇦", "PK": "🇵🇰", "CO": "🇨🇴",
	"IR": "🇮🇷", "EG": "🇪🇬", "SA": "🇸🇦", "AE": "🇦🇪", "CZ": "🇨🇿",
}

func lookupIP(ip string) GeoInfo {
	geoCacheMu.RLock()
	if info, ok := geoCache[ip]; ok {
		geoCacheMu.RUnlock()
		return info
	}
	geoCacheMu.RUnlock()

	info := GeoInfo{Country: "?", City: "?", ISP: "?", Org: "?", AS: "?"}

	client := httpShort
	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,country,countryCode,city,isp,org,as", ip)
	resp, err := client.Get(url)
	if err != nil {
		return info
	}
	defer resp.Body.Close()

	var data struct {
		Status      string `json:"status"`
		Country     string `json:"country"`
		CountryCode string `json:"countryCode"`
		City        string `json:"city"`
		ISP         string `json:"isp"`
		Org         string `json:"org"`
		AS          string `json:"as"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err == nil && data.Status == "success" {
		info = GeoInfo{
			Country: data.Country,
			City:    data.City,
			ISP:     data.ISP,
			Org:     data.Org,
			CC:      data.CountryCode,
			AS:      data.AS,
		}
	}

	geoCacheMu.Lock()
	if len(geoCache) >= geoCacheMaxSize {
		// Evict ~25% of entries
		count := 0
		for k := range geoCache {
			delete(geoCache, k)
			count++
			if count >= geoCacheMaxSize/4 {
				break
			}
		}
	}
	geoCache[ip] = info
	geoCacheMu.Unlock()

	return info
}

func flag(cc string) string {
	if f, ok := flagMap[cc]; ok {
		return f
	}
	return "🏴"
}

func enrichWithGeo(event SecurityEvent) {
	if event.IP == "" {
		return
	}

	geo := lookupIP(event.IP)
	locLine := fmt.Sprintf("%s %s, %s", flag(geo.CC), geo.Country, geo.City)
	ispLine := fmt.Sprintf("🏢 %s", geo.ISP)

	msg := fmt.Sprintf("📍 <b>IP Info</b> — <code>%s</code>\n%s\n%s", event.IP, locLine, ispLine)

	if geo.Org != "" && geo.Org != geo.ISP && geo.Org != "?" {
		msg += fmt.Sprintf("\n🔖 Org: %s", geo.Org)
	}

	sendMessageSmart(msg, nil)
	log.Printf("[geo] Enriched %s: %s, %s", event.IP, geo.Country, geo.City)
}
