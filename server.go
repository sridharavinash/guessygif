package main

import (
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/sridharavinash/guessygif/lib/movies"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

var generator *movies.MovieGenerator

func main() {
	generator, _ = movies.NewGenerator()
	rand.Seed(time.Now().UTC().UnixNano())
	e := echo.New()
	e.Static("/static", "assets")

	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}

	e.Renderer = t

	e.GET("/", indexRender)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	e.Logger.Fatal(e.Start(":" + port))

}

func indexRender(c echo.Context) error {
	movie, _ := generator.GetRandomMovie()
	data := map[string]string{
		"imageUrl":    movie.GifURL,
		"randomMovie": movie.Name,
	}
	return c.Render(http.StatusOK, "index.html", data)
}
