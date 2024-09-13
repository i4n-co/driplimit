package ui

import (
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
)

//go:embed views/*
var viewsfs embed.FS

//go:embed statics/*
var staticsfs embed.FS

func NewTemplateEngine() *html.Engine {
	sub, err := fs.Sub(viewsfs, "views")
	if err != nil {
		panic(err)
	}
	engine := html.NewFileSystem(http.FS(sub), ".html")
	if os.Getenv("VIEWS_DIR") != "" {
		engine = html.New(os.Getenv("VIEWS_DIR"), ".html")
	}

	engine.AddFuncMap(sprig.FuncMap())
	engine.AddFunc("icon", func(icon string) template.HTML {
		svg, err := staticsfs.ReadFile(fmt.Sprintf("statics/feather/%s.svg", icon))
		if err != nil {
			panic(err)
		}
		return template.HTML(string(svg))
	})
	engine.AddFunc("raw", func(html template.HTML) string {
		return string(html)
	})
	return engine
}

func StaticsMiddleware() func(*fiber.Ctx) error {
	fs := http.FS(staticsfs)
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
