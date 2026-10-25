package caption

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)


type WhisperManager struct {
	BINLOCN     string
	WORKING_DIR string
	MODEL_PATH string
}

func NewWhisperManager(locn string, dir string, mdlpth string) *WhisperManager {
	return &WhisperManager{
		BINLOCN:     locn,
		WORKING_DIR: dir,
		MODEL_PATH: mdlpth,
	}
}

func(W *WhisperManager)GenerateCaption(ctx context.Context, wavpath string, vttpath string) error {
	fmt.Println("Entering into generating caption by wav file generated from ffmpeg")
    cmd:=exec.CommandContext(ctx,W.BINLOCN,
		"-m", W.MODEL_PATH,
		"-f", wavpath,
		"-ovtt",
		"-of", strings.TrimSuffix(vttpath, ".vtt"),
	)
    cmd.Stderr=os.Stderr
	return cmd.Run()
}


