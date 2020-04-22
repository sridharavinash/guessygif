package movies

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"

	"github.com/sridharavinash/guessygif/lib/giphy"
)

type Movie struct {
	Name   string
	GifURL string
}

type MovieGenerator struct {
	giphyPicker   *giphy.GiphyPicker
	MovieDBAPIKey string
	movieList     []string
}

func NewGenerator() (*MovieGenerator, error) {
	//Don't care if this is empty, will fallback to using list
	moviedb_api_key := os.Getenv("MOVIEDB_API_KEY")

	picker, err := giphy.NewPicker()
	if err != nil {
		return nil, err
	}

	return &MovieGenerator{
		giphyPicker:   picker,
		MovieDBAPIKey: moviedb_api_key,
		movieList:     GetMovieNamesFromFile(),
	}, nil
}

func (mg *MovieGenerator) GetRandomMovie() (*Movie, error) {
	rint := rand.Intn(len(mg.movieList))
	randomMovie := mg.movieList[rint]

	imageUrl, _ := mg.giphyPicker.GetRandomGiphy(strings.ReplaceAll(randomMovie, " ", "+"))

	movie := &Movie{
		Name:   randomMovie,
		GifURL: imageUrl,
	}

	return movie, nil
}

func GetMovieNamesFromFile() []string {
	fdata, err := ioutil.ReadFile("movies.txt")
	if err != nil {
		fmt.Println("error reading movie txt", err)
		return []string{"The+Departed"}
	}

	return strings.Split(string(fdata), "\n")
}
