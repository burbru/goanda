package api

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var (
	mutex   = sync.Mutex{}
	client  = &http.Client{}
	request = &http.Request{
		Header: http.Header{},
	}
	lastRequestTime time.Time
	rateLimit       = 1 * time.Millisecond
)

func SetHeader(key string, value string) {
	mutex.Lock()
	defer mutex.Unlock()
	request.Header.Set(key, value)
}

func SetToken(token string) {
	mutex.Lock()
	defer mutex.Unlock()
	request.Header.Set("Authorization", "Bearer "+token)
}

func GetHeaderAsString() string {
	mutex.Lock()
	defer mutex.Unlock()

	return headersToString(request.Header.Clone())
}

func SetRateLimit(limit time.Duration) {
	mutex.Lock()
	defer mutex.Unlock()
	rateLimit = limit
}

func headersToString(header http.Header) string {
	var headerString string
	for key, values := range header {
		for _, value := range values {
			headerString += key + ": " + value + "\r\n"
		}
	}
	return headerString
}

func ResetHeaders() {
	mutex.Lock()
	defer mutex.Unlock()
	request.Header = http.Header{}
}

func SendRequest(reqMethod string, reqUrl string, reqBody []byte) ([]byte, error) {
	mutex.Lock()
	defer mutex.Unlock()

	// Calculate the time to wait based on the last request time.
	timeSinceLastRequest := time.Since(lastRequestTime)
	if timeSinceLastRequest < rateLimit {
		timeToWait := rateLimit - timeSinceLastRequest
		time.Sleep(timeToWait)
	}

	request.Method = reqMethod
	parsedURL, err := url.Parse(reqUrl)
	if err != nil {
		PrintWithColor("Error parsing url %s: %s\n", Red, reqUrl, err)
		return nil, err
	}
	request.URL = parsedURL
	request.Body = io.NopCloser(bytes.NewReader(reqBody))

	// Update the last request time
	lastRequestTime = time.Now()

	// Send the request
	resp, err := client.Do(request)

	if err != nil {
		PrintWithColor("Error sending request to url %s: %s", Red, reqUrl, err)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		LogInColor("Error reading response body: %s\n%s", Red, err, string(respBody))
		return nil, err
	}

	return respBody, nil
}

type Color int

const (
	Black Color = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

func ColorizeText(text string, color Color) string {
	return fmt.Sprintf("\033[1;%dm%s\033[0m", 30+int(color), text)
}

func PrintWithColor(format string, color Color, args ...interface{}) {
	for i, v := range args {
		if str, ok := v.(string); ok {
			args[i] = escapePercent(str)
		}
	}
	message := fmt.Sprintf(format, args...)
	coloredMessage := ColorizeText(message, color)
	log.Print(coloredMessage)
}

func escapePercent(s string) string {
	return strings.ReplaceAll(s, "%", "%%")
}

func LogInColor(format string, color Color, args ...interface{}) {
	for i, v := range args {
		if str, ok := v.(string); ok {
			args[i] = escapePercent(str)
		}
	}
	message := fmt.Sprintf(format, args...)

	coloredMessage := ColorizeText(message, color)
	fmt.Printf("%s %s\n", time.Now().Format("2001-02-03 13:14:15"), coloredMessage)
}
