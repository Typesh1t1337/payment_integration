package http_transport

import (
	"log"
	"net/http"
	"payment_integration/internal/domain/uow"
)

type Container struct {
	uow uow.UoW
}

func NewContainer(uow uow.UoW) Container {
	return Container{
		uow: uow,
	}
}

type Runner struct {
	container *Container
}

func NewRunner(container *Container) *Runner {
	return &Runner{container: container}
}

func (r *Runner) Run() {
	mux := http.NewServeMux()
	log.Fatal(http.ListenAndServe(":8080", mux))
}
