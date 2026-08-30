package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	start := time.Now()
	if len(os.Args) < 2 {
		fmt.Printf("URL寄越しやがれください")
		return
	}
	url := os.Args[1]
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	elapsed := time.Since(start)
	fmt.Printf("Response: %s (%v)\n", resp.Status, elapsed)
}
