package domain

type Status string

const (
	StatusPending   Status = "PENDING"
	StatusAvailable Status = "AVAILABLE"
)

func (s Status) IsValid() bool {
	return s == StatusPending || s == StatusAvailable
}
