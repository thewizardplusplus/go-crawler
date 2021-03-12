package registers

// LinkGenerator ...
type LinkGenerator interface {
	GenerateLinks(baseLink string) ([]string, error)
}
