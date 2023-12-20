package xffmpeg

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	ffmpeg := New()
	input := "video.mp4"
	output := "test.jpg"
	ffmpeg.Input(input).Output(output)
	out, err := ffmpeg.GetPreviewImage()
	fmt.Println(err)
	fmt.Println(string(out))
	//"ffmpeg", "-i", input, "-vf", fmt.Sprintf("fps=1/%v:round=zero:start_time=0.99,scale=160:-1,tile=%vx1", per_time, count), output,
}
