package sanitizing

// LinkSanitizing ...
type LinkSanitizing int

// ...
const (
	DoNotSanitizeLink LinkSanitizing = iota
	SanitizeLink
)
