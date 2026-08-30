package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("URL寄越しやがれください")
		return
	}
	url := os.Args[1]

	start := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	CLresponse := "unknown"
	if resp.ContentLength >= 0 {
		CLresponse = fmt.Sprintf("%d bytes", resp.ContentLength)
	}

	bodySize, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	elapsed := time.Since(start)

	fmt.Printf("Response: %s\n", resp.Status)
	fmt.Printf("Content type: %s\n", resp.Header.Get("Content-Type"))
	fmt.Printf("Content length: %s\n", CLresponse)
	fmt.Printf("Body size: %d bytes\n", bodySize)
	fmt.Printf("Time: %s\n", elapsed)
}
