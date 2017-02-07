package models

// Flash deifines a flash message
type Flash struct {
	Type    string
	Title   string
	Message string
}

// NewWarningFlash ..
func NewWarningFlash(title, message string) *Flash {
	return NewFlash("warning", title, message)
}

// NewFlash ..
func NewFlash(flashType, flashTitle, flashMessage string) *Flash {
	return &Flash{
		Type:    flashType,
		Title:   flashTitle,
		Message: flashMessage,
	}
}
