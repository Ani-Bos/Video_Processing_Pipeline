package caption

import "context"

type CaptionManager interface {
	GenerateCaption(ctx context.Context,wavpath string, vttpath string)error
}