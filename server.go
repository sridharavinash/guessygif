package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"math/rand"
	"time"

	"github.com/labstack/echo"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type DownSizedLarge struct {
	URL string `json:"url"`
}

type ImageEntry struct {
	Image DownSizedLarge `json:"downsized_large"`
}

type DataItem struct {
	Images ImageEntry `json:"images"`
}
type APIResponse struct {
	Data []DataItem `json:"data"`
}

//called only once to setup movie names
var movies = getMovieNames()

func main() {
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
	rint := rand.Intn(len(movies))
	randomMovie := movies[rint]
	imageUrl, _ := getRandomGiphy(strings.ReplaceAll(randomMovie, " ", "+"))
	data := map[string]string{
		"imageUrl": imageUrl,
		"randomMovie" : randomMovie,
	}
	return c.Render(http.StatusOK, "index.html", data)
}

func getMovieNames() []string{
	fdata, err := ioutil.ReadFile("movies.txt")
	if err !=  nil{
		fmt.Println("error reading movie txt", err)
		return []string{"The+Departed"}
	}

	return strings.Split(string(fdata), "\n")
}

func getRandomGiphy(s string) (string, error) {

	apiKey := os.Getenv("GIPHY_API_KEY")
	if apiKey == "" {
		fmt.Println("No Giphy API set!")
		os.Exit(1)
	}
	roffset := rand.Intn(5)
	url := fmt.Sprintf("https://api.giphy.com/v1/gifs/search?api_key=%s&q=%s&limit=1&offset=%d&rating=G&lang=en", apiKey, s, roffset)

	payload := strings.NewReader("{}")

	req, err := http.NewRequest("GET", url, payload)

	if err != nil {
		fmt.Println("new request error", err)
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("client call error", err)
		return "", err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var apiResponse APIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		fmt.Println("error unmarshalling", err)
		return "", err
	}

	resp := apiResponse.Data
	imageUrl := resp[0].Images.Image.URL
	return imageUrl, nil
}
