package transcoder

import (
	"fmt"
	"context"
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

func (f *ffmpeg) Transcode(ctx context.Context, job queue.JobQueue) error {
	fmt.Println("Entering into transcoding using ffmpeg and hls format")
	output_dir:= filepath.Join(f.working_dir,job.VideoId)
	err:=os.MkdirAll(output_dir,0755)
	if err!=nil{
		return err
	}
	output_file_locn:=filepath.Join(output_dir,"index.m3u8")
	//use ffmpeg docs
	//This command encodes a video with good quality, using slower preset to achieve better compression:
// ffmpeg -i input -c:v libx264 -preset slow -crf 22 -c:a copy output.mkv
	//fmpeg -i vide_file_name
	//so raw video to HTTP LIve streaming format
	//H.264 format Lossless H.264 ¶
	//https://trac.ffmpeg.org/wiki/Encode/H.264
	//crf constant rate factpr(0-51)
	//audio codec //hls time basically segements cut into tjjhis length 6
	//vod video on demand static not live stream
	cmd:=exec.CommandContext(ctx,f.binlocn,
		"-y", "-i", job.RawPath,
		"-c:v", "libx264", "-preset", "veryfast", "-crf", "23",
		"-c:a", "aac",
		"-hls_time", "6", "-hls_playlist_type", "vod", output_file_locn,
	)
	err=cmd.Run()
	if err!=nil{
		return err
	}
	return nil
}