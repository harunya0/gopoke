package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	pingflag := flag.Bool("ping", false, "ping実行")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("URL寄越しやがれください")
		return
	}
	url := args[0]

	if *pingflag {
		pingRun(url)
	} else {
		pokeRun(url)
	}
}

func pokeRun(url string) {
	start := time.Now()
	result, err := poke(url)
	elapsed := time.Since(start)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	CLresponse := "unknown"
	if result.ContentLength >= 0 {
		CLresponse = fmt.Sprintf("%d bytes", result.ContentLength)
	}
	fmt.Printf("Response: %s\n", result.Status)
	fmt.Printf("Content type: %s\n", result.ContentType)
	fmt.Printf("Content length: %s\n", CLresponse)
	fmt.Printf("Body size: %d bytes\n", result.BodySize)
	fmt.Printf("Time: %s\n", elapsed)
}

func pingRun(url string) {
	host, err := extractHost(url)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	result, err := ping(host)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("Ping result: %d bytes, Type: %d\n", result.Bytes, result.Type)
}

type PokeResult struct {
	Status        string
	ContentType   string
	ContentLength int64
	BodySize      int64
}

func poke(url string) (PokeResult, error) {
	resp, err := http.Get(url)
	if err != nil {
		return PokeResult{}, fmt.Errorf("HTTPリクエストに失敗しました: %w", err)
	}
	defer resp.Body.Close()

	bodySize, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		return PokeResult{}, fmt.Errorf("レスポンスボディの読み込みに失敗しました: %w", err)
	}

	return PokeResult{
		Status:        resp.Status,
		ContentType:   resp.Header.Get("Content-Type"),
		ContentLength: resp.ContentLength,
		BodySize:      bodySize,
	}, nil
}
