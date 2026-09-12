package thumbnail

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"video_processing_pipeline/internal/queue"
)

type ffmpeg struct {
	binlocn     string
	working_dir string
}

func NewFFMPEG(locn string, dir string) *ffmpeg {
	return &ffmpeg{
		binlocn:     locn,
		working_dir: dir,
	}
}

func (f *ffmpeg) GenerateThumbnail(ctx context.Context, job *queue.JobQueue) error {
	fmt.Println("Entering into generating thubnail from FFMPEG output raw hls segments")
    src:=filepath.Join(job.OutPutDir,"index.m3u8")
	dst:=filepath.Join(job.OutPutDir,"ani.jpg")
	//basically after 1 frame peocessing generate image in dst
	//fast seeking to 1 seconf mark
	cmd:=exec.CommandContext(ctx,f.binlocn,
	"-y",
		"-ss", "00:00:01",
		"-i", src,
		"-frames:v", "1",
		"-q:v", "3",
		dst,
	)
	cmd.Stderr = os.Stderr
	err:=cmd.Run()
	if err!=nil{
		return err
	}
	return nil
}