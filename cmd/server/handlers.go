package main

import (
	"context"
	"io"
	"log/slog"
	"sync"

	"connectrpc.com/connect"
	renderv1 "github.com/peter-mghendi/talaria/internal/gen/render/v1"
	service "github.com/peter-mghendi/talaria/internal/render"
)

type RenderServer struct{}

func (s *RenderServer) Render(
	ctx context.Context,
	req *connect.Request[renderv1.RenderRequest],
) (*connect.Response[renderv1.RenderResponse], error) {
	slog.Info("Handling unary Render request.")

	response := service.Render(req.Msg)
	res := connect.NewResponse(&response)

	slog.Info("Render request completed successfully.")
	return res, nil
}

func (s *RenderServer) RenderStream(
	ctx context.Context,
	req *connect.BidiStream[renderv1.RenderStreamRequest, renderv1.RenderStreamResponse],
) error {
	slog.Info("Starting RenderStream")

	responses := make(chan *renderv1.RenderStreamResponse)
	var wg sync.WaitGroup

	go func() {
		for response := range responses {
			if err := req.Send(response); err != nil {
				slog.Warn("Failed to send response", "error", err)
			}
		}
	}()

	for {
		item, err := req.Receive()
		if err == io.EOF {
			break
		}

		if err != nil {
			slog.Error("Error receiving request", "error", err)
			break
		}

		slog.Debug("Received RenderStream request", "identifier", item.Identifier)

		wg.Add(1)
		go func(request *renderv1.RenderStreamRequest) {
			defer wg.Done()

			response := service.Render(request.Request)
			responses <- &renderv1.RenderStreamResponse{Identifier: item.Identifier, Response: &response}
		}(item)
	}

	wg.Wait()
	close(responses)

	slog.Info("RenderStream completed")
	return nil
}
