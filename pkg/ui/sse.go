package ui

import (
	"bufio"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func (s *Server) sse(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		for {
			e, ok := <-s.events
			if !ok {
				return
			}
			_, err := fmt.Fprintf(w, "event:restart\ndata: %s\n\n", e)
			if err != nil {
				s.logger.Warn("writing sse data failed", "err", err)
			}
			err = w.Flush()
			if err != nil {
				s.logger.Warn("sse flushing failed", "err", err)
			}
			s.logger.Debug("sse sent", "e", e)
		}
	}))

	return nil
}
