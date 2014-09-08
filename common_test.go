package sgf

func fnPanics(fn func()) (b bool) {
	defer func() {
		if r := recover(); r != nil {
			b = true
		}
	}()

	fn()
	return b
}
