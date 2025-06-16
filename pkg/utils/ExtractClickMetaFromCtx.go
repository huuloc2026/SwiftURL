package utils

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/huuloc2026/SwiftURL/internal/entity"
	"github.com/mileusna/useragent"
)

// ExtractClickMetaFromCtx reads info from Fiber context to prepare a ClickLog
func ExtractClickMetaFromCtx(c *fiber.Ctx, shortCode string) *entity.ClickLog {
	// Parse User-Agent string
	uaParsed := useragent.Parse(c.Get("User-Agent"))

	deviceType := "desktop"
	if uaParsed.Mobile {
		deviceType = "mobile"
	} else if uaParsed.Tablet {
		deviceType = "tablet"
	}

	return &entity.ClickLog{
		ShortCode:  shortCode,
		ClickedAt:  time.Now(),
		Referrer:   c.Get("Referer"),
		UserAgent:  c.Get("User-Agent"),
		DeviceType: deviceType,
		OS:         uaParsed.OS,
		Browser:    uaParsed.Name,
		IPAddress:  c.IP(),
	}
}

type GeoIP struct {
	Country string `json:"country"`
	City    string `json:"city"`
}

func LookupIP(ip string) (GeoIP, error) {
	var result GeoIP

	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://ip-api.com/json/" + ip)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return result, nil // không fatal
	}

	_ = json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}
