package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/killa-beez/gopkgs/sets/builtins"
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
	addr := os.Getenv("TCP_ADDRESS")
	if addr == "" {
		addr = ":" + port
	}
	e.Logger.Fatal(e.Start(addr))
}

func indexRender(c echo.Context) error {
	options, answer := getOptions(4, rand.Int63())
	imageUrl, _ := getRandomGiphy(strings.ReplaceAll(options[answer], " ", "+"))
	data := viewData{
		ImageURL: imageUrl,
		Correct:  answer,
		Choices:  options,
	}
	return c.Render(http.StatusOK, "index.html", data)
}

type viewData struct {
	RandomMovie string
	ImageURL    string
	Choices     []string
	Correct     int
}

//returns count titles and the index for the correct answer
func getOptions(count int, seed int64) ([]string, int) {
	rnd := rand.New(rand.NewSource(seed))
	set := builtins.NewStringSet(count)
	for set.Len() < count {
		set.Add(movies[rnd.Intn(len(movies))])
	}
	options := set.Values()
	rnd.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})
	return options, rnd.Intn(count)
}

func getMovieNames() []string {
	fdata, err := ioutil.ReadFile("movies.txt")
	if err != nil {
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
