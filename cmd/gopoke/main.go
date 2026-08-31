package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/harunya0/gopoke"
)

func main() {
	pingFlag := flag.Bool("ping", false, "ping実行")
	trackingFlag := flag.Bool("tracking", false, "リダイレクト追跡実行")
	timeoutFlag := flag.Int("timeout", 10, "タイムアウト時間")
	maxRedirectsFlag := flag.Int("max-redirects", 10, "リダイレクト追跡回数の上限(0で無制限)")
	jsonFlag := flag.Bool("json", false, "JSON形式で出力")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("URL寄越しやがれください")
		return
	}
	url := args[0]

	if *pingFlag {
		pingRun(url)
	} else if *trackingFlag {
		trackRun(url, time.Duration(*timeoutFlag)*time.Second, *maxRedirectsFlag)
	} else {
		pokeRun(url, time.Duration(*timeoutFlag)*time.Second, *jsonFlag)
	}
}

func pokeRun(url string, timeout time.Duration, jsonOutput bool) {
	start := time.Now()
	result, err := gopoke.Poke(url, timeout)
	elapsed := time.Since(start)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	CLresponse := "unknown"
	if result.ContentLength >= 0 {
		CLresponse = fmt.Sprintf("%d bytes", result.ContentLength)
	}
	if jsonOutput {
		out := struct {
			gopoke.PokeResult
			Elapsed string `json:"elapsed"`
		}{
			PokeResult: result,
			Elapsed:    elapsed.String(),
		}
		b, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(b))
		return
	}
	fmt.Printf("Response: %s\n", result.Status)
	fmt.Printf("Content type: %s\n", result.ContentType)
	fmt.Printf("Content length: %s\n", CLresponse)
	fmt.Printf("Body size: %d bytes\n", result.BodySize)
	fmt.Printf("Time: %s\n", elapsed)
}

func pingRun(url string) {
	host, err := gopoke.ExtractHost(url)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	result, err := gopoke.Ping(host)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("Ping result: %d bytes, Type: %d\n", result.Bytes, result.Type)
}

func trackRun(url string, timeout time.Duration, maxRedirects int) {
	result, err := gopoke.Track(url, timeout, maxRedirects)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	if len(result.Chain) == 0 {
		fmt.Println("リダイレクトはありませんでした")
		return
	}
	fmt.Println("Redirect chain:")
	for i, u := range result.Chain {
		fmt.Printf("%d: %s\n", i+1, u)
	}
	fmt.Printf("Final HTTP status code: %d\n", result.FinalCode)
}
