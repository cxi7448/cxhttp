package xffmpeg

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/xoss/s3"
	"io/ioutil"
	"testing"
)

func init() {
}

func TestNew3(t *testing.T) {
	xs3 := s3.NewVideo()
	files, _ := ioutil.ReadDir("./tmp")
	for _, file := range files {
		filename := fmt.Sprintf("./tmp/%v", file.Name())
		objectName := fmt.Sprintf("books/mp3/%v", file.Name())
		xs3.UploadFile(filename, objectName, 3)
	}

	// https://mt.qkmkpz.com/books/mp3/playlist.m3u8
}

func TestNew(t *testing.T) {
	ffmpeg := New()
	input := "video.mp4"
	output := "test.jpg"
	ffmpeg.Input(input).Output(output)
	out, err := ffmpeg.GetPreviewImage()
	fmt.Println(err)
	fmt.Println(string(out))
	// ffmpeg -i input.mp3 -codec:a libmp3lame -b:a 128k -ar 44100 -ac 2 -hls_time 10 -hls_list_size 0 -hls_wrap 10 -f hls -hls_key_info_file encryption.keyinfo output.m3u8
}

// ffmpeg -i input.mp3 -codec:a libmp3lame -q:a 4 output.mp3
func TestNew2(t *testing.T) {
	ffmpeg := New()
	input := "xxx.mp3"
	folder := "tmp"
	ts_format := folder + "/output%d.mp3"
	m2u8_path := folder + "/playlist.m3u8"
	//mp3 := "output.mp3" // "-map", "0",
	//ffmpeg.AddArgs("-c:a", "libmp3lame", "-aq", "0", "-f", "segment", "-segment_time", "10", "-segment_list", m2u8_path)
	ffmpeg.AddArgs("-c:a", "libmp3lame", "-aq", "0", "-map", "0", "-f", "segment", "-segment_time", "10", "-segment_list", m2u8_path)
	//ffmpeg.AddArgs("-c:a", "libmp3lame", "-q:a", "0")
	//ffmpeg -i "%%~sa" -y -acodec libmp3lame -aq 0 -map 0 -f segment -segment_time 180 -write_xing 0 "result\%%~na-%%03d.mp3"

	ffmpeg.Input(input).Output(ts_format)
	//ffmpeg.Input(input).Output(m2u8_path)
	out, err := ffmpeg.Run()
	fmt.Println(err)
	fmt.Println(string(out))

	//fmt.Println(Mp4ToM3u8Encrypt(input, folder))

}

//
//# 假设原始文件是original.mp3
//
//# 步骤1: 切割MP3文件
//ffmpeg -i original.mp3 -codec:copy libmp3lame -map 0 -f segment -segment_time 10 output%03d.mp3
//
//# 步骤2: 生成密钥
//for file in output*.mp3; do
//ffmpeg -i "$file" -c copy -bsf:a aac_adtstoasc -movflags +faststart "$file.enc"
//ffmpeg -i "$file.enc" -c copy -bsf:a aes_setiv -iv 0x00000000000000000000000000000001 -f mpegts -encryption_format_id 5 -encryption_kid 00000000-0000-0000-0000-000000000000 -encryption_key 00000000000000000000000000000000 -movflags +faststart "$file.enc"
//done
//
//# 步骤3: 创建加密的M3U8播放列表
//echo "#EXTM3U" > encrypted.m3u8
//echo "#EXT-X-VERSION:3" >> encrypted.m3u8
//echo "#EXT-X-MEDIA-SEQUENCE:0" >> encrypted.m3u8
//echo "#EXT-X-KEY:METHOD=AES-128,URI=\"https://example.com/key.php?r=123\"" >> encrypted.m3u8
//
//for file in output*.enc.m3u8; do
//echo "#EXTINF:10.000000,$file" >> encrypted.m3u8
//done
