package entity

type Media struct {
	Title          string  `json:"title"`
	Filesize       int     `json:"filesize"`
	FilesizeApprox int     `json:"filesize_approx"`
	Duration       float32 `json:"duration"`
	URL            string  `json:"url"`
	WebpageURL     string  `json:"webpage_url"`
	Thumbnail      string  `json:"thumbnail"`
}
