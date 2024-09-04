package ui

import (
	"embed"
	"io/fs"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html/v2"
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
