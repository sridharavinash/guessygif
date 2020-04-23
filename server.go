package main

import (
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/killa-beez/gopkgs/sets/builtins"
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
	e.GET("/refresh", refreshGif)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	addr := os.Getenv("TCP_ADDRESS")
	if addr == "" {
		addr = ":" + port
	}
	e.Logger.Fatal(e.Start(addr))
}

func indexRender(c echo.Context) error {
	options, answer := getOptions(4)
	title := options[answer]
	gif, err := generator.GetMovieGif(title, 3)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return c.Render(http.StatusOK, "index.html", viewData{
		Title:    title,
		Choices:  options,
		Correct:  answer,
		ImageURL: gif,
	})
}

func refreshGif(c echo.Context) error {
	title := c.QueryParam("title")
	gif, err := generator.GetMovieGif(title, 10)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.String(http.StatusOK, gif)
}

type viewData struct {
	Title    string
	ImageURL string
	Choices  []string
	Correct  int
}

//returns count titles and the index for the correct answer
func getOptions(count int) ([]string, int) {
	set := builtins.NewStringSet(count)
	for set.Len() < count {
		set.Add(generator.GetRandomMovie())
	}
	return set.Values(), generator.Intn(set.Len())
}
