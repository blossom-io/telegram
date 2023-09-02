package ffmpeg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
)

const (
	FFprobe = "ffprobe"
	FFmpeg  = "ffmpeg"
)

var (
	FFprobeMediaInfoArgs = []string{"-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height,duration", "-of", "json", "-"}                                    // Get width, height and duration of the input media
	FFmpegFirstFrameArgs = []string{"-i", "-", "-vf", "select=eq(n\\,0),scale='if(gt(iw,ih),min(320,iw),-1)':'if(gt(iw,ih),-1,min(320,ih))'", "-q:v", "1", "-f", "image2", "pipe:1"} // Downscale to 320px each side infopecting the aspect ratio
)

type Info struct {
	Programs []any `json:"programs"`
	Streams  []struct {
		Width    int    `json:"width"`
		Height   int    `json:"height"`
		Duration string `json:"duration"`
	} `json:"streams"`
}

// GetFileInfo returns width, height and duration of the input media data
func GetFileInfo(ctx context.Context, data *bytes.Reader) (info Info, err error) {
	FFprobeCmd := exec.CommandContext(ctx, FFprobe, FFprobeMediaInfoArgs...)
	FFprobeCmd.Stdin = data

	FFprobeOut, err := FFprobeCmd.Output()
	if err != nil {
		return info, err
	}
	if err := json.Unmarshal(FFprobeOut, &info); err != nil {
		return info, err
	}

	if len(info.Streams) == 0 {
		return info, fmt.Errorf("ffprobe: no streams")
	}

	return info, nil
}

func GetFirstFrame(ctx context.Context, data *bytes.Reader) (frame []byte, err error) {
	FFmpegCmd := exec.CommandContext(ctx, FFmpeg, FFmpegFirstFrameArgs...)
	FFmpegCmd.Stdin = data

	frame, err = FFmpegCmd.Output()
	if err != nil {
		return frame, err
	}

	return frame, nil
}
