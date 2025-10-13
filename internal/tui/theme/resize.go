package theme

type Size struct {
	Width  int
	Height int
}

type Resizable interface {
	SetSize(Size)
}
