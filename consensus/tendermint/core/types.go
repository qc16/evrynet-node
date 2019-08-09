package core

type Engine interface {
	Start() error
	Stop() error
}
