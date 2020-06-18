package api

type Controlable interface {
	Start() string
	Stop() string
	GetStatus() string
}
