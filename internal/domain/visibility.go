package domain

type Visibility string

const (
	VisibilityPublic  Visibility = "PUBLIC"
	VisibilityPrivate Visibility = "PRIVATE"
)

func (v Visibility) IsValid() bool {
	return v == VisibilityPublic || v == VisibilityPrivate
}
