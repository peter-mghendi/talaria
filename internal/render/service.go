package service

import renderv1 "github.com/peter-mghendi/talaria/internal/gen/render/v1"

func Render(request *renderv1.RenderRequest) renderv1.RenderResponse {
	return renderv1.RenderResponse{
		Html: "<p>Hey there!</p>",
		Text: "Hey there!",
	}
}
