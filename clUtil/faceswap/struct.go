package faceswap

const (
	TYPE_AKOOL = "akool" // https://docs.akool.io/ai-tools-suite/faceswap
)

type Api interface {
	FaceSwap(src, face string) error
	Video()
}
