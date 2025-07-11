package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	// "crypto/x509"
	// "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	// "net"
	"net/http"
	"os"
	"strings"
	"sync"

	// utls "github.com/refraction-networking/utls"
	// "golang.org/x/net/http2"
	"waguri-joiner/modules"
)

type Config struct {
	Proxy   string `json:"proxy"`
	Threads int    `json:"threads"`
	Invite  string `json:"invite"`
}

func LoadConfig() Config {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Failed to open config.json: %v", err)
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		log.Fatalf("Failed to parse config.json: %v", err)
	}
	return cfg
}

func LoadTokens() []string {
	file, err := os.Open("tokens.txt")
	if err != nil {
		log.Fatalf("Failed to open tokens.txt: %v", err)
	}
	defer file.Close()

	var tokens []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		split := strings.Split(line, ":")
		token := line
		if len(split) >= 3 {
			token = split[2]
		}
		tokens = append(tokens, token)
	}
	return tokens
}

func JoinInvite(token, invite, proxy string) {
	clientConn, _, err := headers.DialTLS("discord.com", proxy)
	if err != nil {
		headers.LogFailure(token, 0)
		return
	}

	payload := map[string]string{"session_id": headers.GenSession()}
	jsonPayload, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", fmt.Sprintf("https://discord.com/api/v9/invites/%s", invite), bytes.NewBuffer(jsonPayload))
	if err != nil {
		headers.LogFailure(token, 0)
		return
	}

	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-encoding", "gzip, deflate, br, zstd")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("authorization", token)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("origin", "https://discord.com")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", "https://discord.com/channels/@me")
	req.Header.Set("sec-ch-ua", `"Not A;Brand";v="99", "Chromium";v="131", "Microsoft Edge";v="131"`)
	req.Header.Set("sec-ch-ua-platform", "Windows")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	req.Header.Set("x-debug-options", "bugReporterEnabled")
	req.Header.Set("x-discord-locale", "en-US")
	req.Header.Set("x-discord-timezone", "Asia/Katmandu")
	req.Header.Set("x-super-properties", headers.BuildXsup())
	req.Header.Set("cookie", headers.GetDiscordCookies())

	resp, err := clientConn.RoundTrip(req)
	if err != nil {
		headers.LogFailure(token, 0)
		return
	}
	defer resp.Body.Close()

	var reader io.ReadCloser
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			headers.LogFailure(token, resp.StatusCode)
			return
		}
		defer reader.Close()
	} else {
		reader = resp.Body
	}

	_, err = io.ReadAll(reader)
	if err != nil {
		headers.LogFailure(token, resp.StatusCode)
		return
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		headers.LogSuccess(token, resp.StatusCode)
	} else {
		headers.LogFailure(token, resp.StatusCode)
	}
}

func main() {
	config := LoadConfig()
	tokens := LoadTokens()

	var wg sync.WaitGroup
	sem := make(chan struct{}, config.Threads)

	for _, token := range tokens {
		wg.Add(1)
		sem <- struct{}{}
		go func(t string) {
			defer wg.Done()
			JoinInvite(t, config.Invite, config.Proxy)
			<-sem
		}(token)
	}

	wg.Wait()
}
