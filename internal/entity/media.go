package entity

import (
	"blossom/pkg/ffmpeg"
	"blossom/pkg/ytdlp"
	"strconv"
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

func (m *Media) SetMedia(media ytdlp.Media) {
	m.Title = media.Title
	m.Extractor = media.Extractor
	m.Filesize = media.Filesize
	m.FilesizeApprox = media.FilesizeApprox
	m.Duration = media.Duration
	m.URL = media.URL
	m.WebpageURL = media.WebpageURL
	m.Thumbnail = media.Thumbnail
	m.Height = media.Height
	m.Width = media.Width
}

func (m *Media) SetInfo(info ffmpeg.Info) {
	if len(info.Streams) == 0 {
		return
	}

	dur, _ := strconv.ParseFloat(info.Streams[0].Duration, 32)
	m.Duration = dur
	m.Width = info.Streams[0].Width
	m.Height = info.Streams[0].Height
}
