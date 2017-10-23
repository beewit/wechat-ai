package enum

// Virtual key codes
const (
	VK_V = 86
)

type Offset struct {
	Left float32 `json:"left"`
	Top  float32 `json:"top"`
}
type Size struct {
	Width, Height int32
}
