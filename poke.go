package gopoke

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type PokeResult struct {
	Status        string `json:"status"`
	ContentType   string `json:"content_type"`
	ContentLength int64  `json:"content_length"`
	BodySize      int64  `json:"body_size"`
}

func Poke(url string, timeout time.Duration) (PokeResult, error) {
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(url)
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
