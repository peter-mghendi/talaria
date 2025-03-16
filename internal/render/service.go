package service

import (
	"fmt"
	"strings"

	"github.com/go-hermes/hermes/v2"
	renderv1 "github.com/peter-mghendi/talaria/internal/gen/render/v1"
)

func Render(request *renderv1.RenderRequest) renderv1.RenderResponse {
	hermes := ConvertProtoToHermes(request.Hermes)
	email := ConvertProtoToEmail(request.Email)

	html, _ := hermes.GenerateHTML(email)
	text, _ := hermes.GeneratePlainText(email)

	return renderv1.RenderResponse{Html: html, Text: text}
}

// ConvertProtoToHermes converts renderv1.Hermes to hermes.Hermes
func ConvertProtoToHermes(h *renderv1.Hermes) hermes.Hermes {
	return hermes.Hermes{
		Theme:              GetTheme(&h.Theme),
		TextDirection:      hermes.TextDirection(h.TextDirection),
		DisableCSSInlining: h.DisableCssInlining,
		Product: hermes.Product{
			Name:        h.Product.Name,
			Link:        h.Product.Link,
			Logo:        h.Product.Logo,
			Copyright:   h.Product.Copyright,
			TroubleText: h.Product.TroubleText,
		},
	}
}

// ConvertProtoToEmail converts renderv1.Email to hermes.Email
func ConvertProtoToEmail(e *renderv1.Email) hermes.Email {
	var tables []hermes.Table
	for _, t := range e.Body.Tables {
		var tableData [][]hermes.Entry
		for _, entry := range t.Data {
			tableData = append(tableData, []hermes.Entry{
				{
					Key:   entry.Key,
					Value: entry.Value,
				},
			})
		}
		tables = append(tables, hermes.Table{
			Title: t.Title,
			Data:  tableData,
		})
	}

	var actions []hermes.Action
	for _, a := range e.Body.Actions {
		actions = append(actions, hermes.Action{
			Instructions: a.Instructions,
			InviteCode:   a.InviteCode,
			Button: hermes.Button{
				Color:     a.Button.Color,
				TextColor: a.Button.TextColor,
				Text:      a.Button.Text,
				Link:      a.Button.Link,
			},
		})
	}

	return hermes.Email{
		Body: hermes.Body{
			Name:      e.Body.Name,
			Intros:    e.Body.Intros,
			Outros:    e.Body.Outros,
			Greeting:  e.Body.Greeting,
			Signature: e.Body.Signature,
			Title:     e.Body.Title,
			Dictionary: func() []hermes.Entry {
				var entries []hermes.Entry
				for _, entry := range e.Body.Dictionary {
					entries = append(entries, hermes.Entry{
						Key:   entry.Key,
						Value: entry.Value,
					})
				}
				return entries
			}(),
			Tables:  tables,
			Actions: actions,
		},
	}
}

func GetTheme(name *string) hermes.Theme {
	if name == nil {
		return new(hermes.Default)
	}

	switch theme := strings.ToLower(strings.TrimSpace(*name)); theme {
	case "default":
		return new(hermes.Default)
	case "flat":
		return new(hermes.Flat)
	default:
		panic(fmt.Errorf("unknown theme: %s", theme))
	}
}
