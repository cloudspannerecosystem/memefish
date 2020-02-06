package analyzer

type Placeholder struct {
	*PlaceholderType
}

func NewPlaceholder() *Placeholder {
	return &Placeholder{
		PlaceholderType: &PlaceholderType{},
	}
}
