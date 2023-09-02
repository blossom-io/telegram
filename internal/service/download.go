package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"blossom/internal/entity"
	"blossom/pkg/ffmpeg"
	"blossom/pkg/ytdlp"
)

const (
	YTdlp   = "yt-dlp"
	FFprobe = "ffprobe"
	FFmpeg  = "ffmpeg"
	//Best  = "bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best"
	Best          = "best[ext=mp4]/best"
	FilesizeLimit = 50 * 1024 * 1024
)

var (
	FFprobeMediaInfoArgs = []string{"-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height,duration", "-of", "json", "-"}                                    // Get width, height and duration of the input media
	FFmpegFirstFrameArgs = []string{"-i", "-", "-vf", "select=eq(n\\,0),scale='if(gt(iw,ih),min(320,iw),-1)':'if(gt(iw,ih),-1,min(320,ih))'", "-q:v", "1", "-f", "image2", "pipe:1"} // Downscale to 320px each side respecting the aspect ratio
)

type FFprobeStdout struct {
	Programs []any `json:"programs"`
	Streams  []struct {
		Width    int    `json:"width"`
		Height   int    `json:"height"`
		Duration string `json:"duration"`
	} `json:"streams"`
}

type Downloader interface {
	Fetch(ctx context.Context, url string) (media *entity.Media, err error)
	IsDownloaderEnabled(ctx context.Context, chatID int64) (bool, error)
}

// Fetch fetches best quality mp4 format from the given url
func (svc *service) Fetch(ctx context.Context, url string) (*entity.Media, error) {
	media := &entity.Media{}

	svc.log.Info("got url", "url", url)

	res, err := ytdlp.FetchMedia(ctx, url)
	if err != nil {
		return media, err
	}

	media.SetMedia(res)

	if max(media.Filesize, media.FilesizeApprox) > FilesizeLimit {
		return media, fmt.Errorf("file too big")
	}

	svc.log.Debug("media", "media", fmt.Sprintf("%+v", media))

	media.Body, err = svc.Get(ctx, media.URL)
	if err != nil {
		return media, err
	}

	info, err := ffmpeg.GetFileInfo(ctx, bytes.NewReader(media.Body))
	if err != nil {
		return media, err
	}

	media.SetInfo(info)

	media.Preview, err = ffmpeg.GetFirstFrame(ctx, bytes.NewReader(media.Body))
	if err != nil {
		return media, err
	}

	return media, nil
}

func (svc *service) IsDownloaderEnabled(ctx context.Context, chatID int64) (bool, error) {
	return svc.repo.IsDownloaderEnabled(ctx, chatID)
}

func (svc *service) Get(ctx context.Context, url string) (data []byte, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return data, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return data, err
	}
	defer res.Body.Close()

	data, err = io.ReadAll(res.Body)
	if err != nil {
		return data, err
	}

	return data, nil
}
