package clFile

import (
	"testing"
)

func TestGetFileMD5(t *testing.T) {
	localPath := "1.png"
	link := "https://d2qf6ukcym4kn9.cloudfront.net/final_6f4000e4e94940f9bdb4834180dca9d7-0cedf7ea-3d21-49f4-959c-e47bdcb2f113-5354.png"
	Download(link, localPath)
}
