package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/marstr/quickstart/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
