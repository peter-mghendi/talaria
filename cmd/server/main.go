package main

import (
	"log/slog"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/peter-mghendi/talaria/internal/gen/render/v1/renderv1connect"
)

func main() {
	address := "0.0.0.0:9999"
	path, handler := renderv1connect.NewRenderServiceHandler(&RenderServer{})
	mux := http.NewServeMux()
	mux.Handle(path, handler)

	slog.Info("Talaria server is starting", "address", address, "transport", "h2c")
	err := http.ListenAndServe(address, h2c.NewHandler(mux, &http2.Server{}))
	if err != nil {
		slog.Error("Server failed to start", "error", err.Error())
	}
}
