package main

import (
	"context"
	"embed"
	"encoding/json"
	"os"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/i4n-co/driplimit/pkg/api"
	"github.com/i4n-co/driplimit/pkg/config"
)

//go:embed *.template.md
var templates embed.FS

func main() {
	cfg, err := config.FromEnv(context.Background())
	if err != nil {
		panic(err)
	}
	docs, err := api.New(cfg, nil).GenerateDocs()
	if err != nil {
		panic(err)
	}

	t, err := template.New("").Funcs(funcMap()).ParseFS(templates, "*.md")
	if err != nil {
		panic(err)
	}

	err = t.ExecuteTemplate(os.Stdout, "rpc-api-v1.template.md", docs)
	if err != nil {
		panic(err)
	}
}

func funcMap() template.FuncMap {
	fmap := sprig.FuncMap()
	fmap["json"] = func(s any) string {
		jsn, err := json.MarshalIndent(s, "", "  ")
		if err != nil {
			panic(err)
		}
		return string(jsn)
	}
	return fmap
}
