package urlutils

// LinkTrimming ...
type LinkTrimming int

// ...
const (
	DoNotTrimLink LinkTrimming = iota
	TrimLinkLeft
	TrimLinkRight
	TrimLink
)
