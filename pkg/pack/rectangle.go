package pack

type Rectangle struct {
	X, Y, Width, Height int
}

func (r Rectangle) isLarger(other Rectangle) bool {
	return r.Width > other.Width || r.Height > other.Height
}

func (r Rectangle) top() int {
	return r.Y
}

func (r Rectangle) bottom() int {
	return r.Y + r.Height
}

func (r Rectangle) left() int {
	return r.X
}

func (r Rectangle) right() int {
	return r.X + r.Width
}

func (r Rectangle) overlaps(other Rectangle) bool {
	return r.left() < other.right() && r.right() > other.left() && r.top() < other.bottom() && r.bottom() > other.top()
}

func (r Rectangle) cut(cut Rectangle) []Rectangle {
	top := Rectangle{r.left(), r.top(), r.Width, cut.top() - r.top()}
	bottom := Rectangle{r.left(), cut.bottom(), r.Width, r.bottom() - cut.bottom()}
	left := Rectangle{r.left(), r.top(), cut.left() - r.left(), r.Height}
	right := Rectangle{cut.right(), r.top(), r.right() - cut.right(), r.Height}

	var rects []Rectangle

	for _, r := range []Rectangle{top, bottom, left, right} {
		if r.Width > 0 && r.Height > 0 {
			rects = append(rects, r)
		}
	}
	return rects
}
