package service

import (
	"blossom/internal/entity"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
)

const (
	Ytdlp = "yt-dlp"
	//Best  = "bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best"
	Best          = "best[ext=mp4]/best"
	FilesizeLimit = 50 * 1024 * 1024
)

type Downloader interface {
	Download(ctx context.Context, url string) (media entity.Media, data []byte, thumb []byte, err error)
	IsDownloaderEnabled(ctx context.Context, chatID int64) (bool, error)
}

func (svc *service) Download(ctx context.Context, url string) (media entity.Media, data []byte, thumb []byte, err error) {
	svc.log.Info("got url", "url", url)

	ytdlpCmd := exec.Command(Ytdlp, "-f", Best, "-j", url)

	ytdlpOut, err := ytdlpCmd.Output()
	if err != nil {
		return media, nil, nil, err
	}

	if err := json.Unmarshal(ytdlpOut, &media); err != nil {
		return media, nil, nil, err
	}

	if media.Filesize > FilesizeLimit || media.FilesizeApprox > FilesizeLimit {
		return media, nil, nil, fmt.Errorf("file too big")
	}

	svc.log.Debug("media", "media", fmt.Sprintf("%+v", media))

	data, err = svc.get(ctx, media.URL)
	if err != nil {
		return media, nil, nil, err
	}

	thumb, err = svc.get(ctx, media.Thumbnail)
	if err != nil {
		return media, data, nil, err
	}

	return media, data, thumb, nil
}

func (svc *service) IsDownloaderEnabled(ctx context.Context, chatID int64) (bool, error) {
	return svc.repo.IsDownloaderEnabled(ctx, chatID)
}

func (svc *service) get(ctx context.Context, url string) (data []byte, err error) {
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
