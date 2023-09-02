package ytdlp

import (
	"context"
	"encoding/json"
	"os/exec"
)

const (
	YTdlp = "yt-dlp"
	//Best  = "bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best"
	Best          = "best[ext=mp4]/best"
	FilesizeLimit = 50 * 1024 * 1024
)

type Media struct {
	Body           []byte  `json:"-"`
	Preview        []byte  `json:"-"`
	Title          string  `json:"title"`
	Extractor      string  `json:"extractor"`
	Filesize       int     `json:"filesize"`
	FilesizeApprox int     `json:"filesize_approx"`
	Duration       float64 `json:"duration"`
	URL            string  `json:"url"`
	WebpageURL     string  `json:"webpage_url"`
	Thumbnail      string  `json:"thumbnail"`
	Height         int     `json:"height"`
	Width          int     `json:"width"`
}

// FetchMedia fetches best quality mp4 format from the given url
func FetchMedia(ctx context.Context, url string) (media Media, err error) {
	ytdlpCmd := exec.CommandContext(ctx, YTdlp, "-f", Best, "-j", url)

	ytdlpOut, err := ytdlpCmd.Output()
	if err != nil {
		return media, err
	}

	if err := json.Unmarshal(ytdlpOut, &media); err != nil {
		return media, err
	}

	return media, nil
}
