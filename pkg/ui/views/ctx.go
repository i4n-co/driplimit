package views

import (
	"log/slog"

	"github.com/i4n-co/driplimit"
)

type Ctx struct {
	Service driplimit.Service
	Logger  *slog.Logger
}
