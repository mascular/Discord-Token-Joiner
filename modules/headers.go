package headers

import (
    "encoding/base64"
    "encoding/json"
	"fmt"
	"net/http"
	"strings"
    fhttp "github.com/bogdanfinn/fhttp"
    "crypto/rand"
	"encoding/hex"
)

func BuildXsup() string {
    buildNumber := 417266 // fallback
	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
    props := map[string]interface{}{
        "os":                  "Windows",
        "browser":             "Chrome",
        "device":              "",
        "system_locale":       "en-US",
        "browser_user_agent":  userAgent,
        "browser_version":     "133.0.0.0",
        "os_version":          "10",
        "referrer":            "https://discord.com/channels/@me",
        "referring_domain":    "discord.com",
        "release_channel":     "stable",
        "client_build_number": buildNumber,
    }

    jsonBytes, err := json.Marshal(props)
    if err != nil {
        return ""
    }

    encoded := base64.StdEncoding.EncodeToString(jsonBytes)
    return encoded
}
func GetDiscordCookies() string {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://discord.com/api/v9/experiments", nil)
	if err != nil {
		fmt.Println("Request error:", err)
		return ""
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("HTTP error:", err)
		return ""
	}
	defer resp.Body.Close()

	cookies := []string{}
	for _, c := range resp.Cookies() {
		cookies = append(cookies, fmt.Sprintf("%s=%s", c.Name, c.Value))
	}

	return strings.Join(cookies, "; ")
}

func GenSession() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func GetHeaders(tokenasw string) fhttp.Header {
    headers := fhttp.Header{
        "accept":              {"*/*"},
        "accept-encoding":     {"gzip, deflate, br, zstd"},
        "accept-language":     {"en-US,en;q=0.9"},
        "authorization":       {tokenasw},
        "content-type":        {"application/json"},
        "origin":              {"https://discord.com"},
        "priority":            {"u=1, i"},
        "referer":             {"https://discord.com/channels/@me"},
        "sec-ch-ua":           {`"Not)A;Brand";v="99", "Microsoft Edge";v="133", "Chromium";v="133"`},
        "sec-ch-ua-platform":  {"Windows"},
        "sec-fetch-dest":      {"empty"},
        "sec-fetch-mode":      {"cors"},
        "user-agent":          {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"},
        "x-debug-options":     {"bugReporterEnabled"},
        "x-discord-locale":    {"en-US"},
        "x-discord-timezone":  {"Asia/Katmandu"},
        "x-super-properties":  {BuildXsup()},
        "cookie":              {GetDiscordCookies()},
    }

    return headers
}