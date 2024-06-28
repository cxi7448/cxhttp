package faceswap

const (
	TYPE_AKOOL = "akool" // https://docs.akool.io/ai-tools-suite/faceswap
)

type Api interface {
	FaceSwap(src, face Img) error
	Video()
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
