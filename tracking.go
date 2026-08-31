package main

import (
	"fmt"
	"net/http"
	"time"
)

type TrackResult struct {
	Chain     []string
	FinalCode int
}

func Track(url string, timeout time.Duration, maxRedirects int) (TrackResult, error) {
	var chain []string

	client := &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			chain = append(chain, req.URL.String())
			if maxRedirects > 0 && len(via) >= maxRedirects {
				return fmt.Errorf("リダイレクトの最大回数を超えました: %d", maxRedirects)
			}
			return nil
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return TrackResult{}, fmt.Errorf("HTTPリクエストに失敗しました: %w", err)
	}
	defer resp.Body.Close()

	return TrackResult{
		Chain:     chain,
		FinalCode: resp.StatusCode,
	}, nil
}
