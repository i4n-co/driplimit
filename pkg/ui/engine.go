package ui

import (
	"bufio"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"

	"github.com/Masterminds/sprig"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html/v2"
	"github.com/valyala/fasthttp"
)

//go:embed views/*
var views embed.FS

//go:embed statics/*
var statics embed.FS

func loadEngine() *html.Engine {
	sub, err := fs.Sub(views, "views")
	if err != nil {
		panic(err)
	}
	engine := html.NewFileSystem(http.FS(sub), ".html")
	if os.Getenv("VIEWS_DIR") != "" {
		engine = html.New(os.Getenv("VIEWS_DIR"), ".html")
		engine.Reload(true)
	}

	engine.AddFuncMap(sprig.FuncMap())
	engine.AddFunc("icon", func(icon string) template.HTML {
		svg, err := statics.ReadFile(fmt.Sprintf("statics/feather/%s.svg", icon))
		if err != nil {
			panic(err)
		}
		return template.HTML(string(svg))
	})
	return engine
}

func loadStaticsMiddleware() func(*fiber.Ctx) error {
	fs := http.FS(statics)
	pathPrefix := "statics"
	if os.Getenv("STATICS_DIR") != "" {
		fs = http.FS(os.DirFS(os.Getenv("STATICS_DIR")))
		pathPrefix = ""
	}
	return filesystem.New(filesystem.Config{
		Root:       fs,
		PathPrefix: pathPrefix,
		Browse:     true,
	})
}

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
