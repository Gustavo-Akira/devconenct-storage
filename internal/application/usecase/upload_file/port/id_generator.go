package port

type IdGenerator interface {
	Generate() string
}
