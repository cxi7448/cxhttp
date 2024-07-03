package faceswap

const (
	TYPE_AKOOL = "akool" // https://docs.akool.io/ai-tools-suite/faceswap
)

type Api interface {
	FaceSwap(src, face Img) (string, string, error)
	FaceSwapVideo(src, face []Img, video_url string) (string, string, error)
	CheckResult(id string) (uint32, error)
	Ty()
}

type Img struct {
	Image string // 图片路径，绝对路径
	Opts  string // 脸部信息
}

func NewImg(src, opts string) Img {
	return Img{
		Image: src,
		Opts:  opts,
	}
}
