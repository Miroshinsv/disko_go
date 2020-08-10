package web_server

import (
	"context"
	"github.com/gorilla/mux"
)

type IWebServer interface {
	RegisterRoutes(router *mux.Router)
	ListenAndServe(ctx context.Context) error
	Stop(ctx context.Context) error
	IsRunning() bool
}
