package movies

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"

	"github.com/sridharavinash/guessygif/lib/giphy"
	"github.com/sridharavinash/guessygif/lib/moviedb"
)

type Movie struct {
	Name   string
	GifURL string
}

type MovieGenerator struct {
	giphyPicker  *giphy.GiphyPicker
	movieFetcher *moviedb.MovieFetcher
	movieList    []string
}

func NewGenerator() (*MovieGenerator, error) {
	picker, err := giphy.NewPicker()
	if err != nil {
		return nil, err
	}

	fetcher, err := moviedb.NewFetcher()
	if err != nil {
		return nil, err
	}

	return &MovieGenerator{
		giphyPicker:  picker,
		movieFetcher: fetcher,
		movieList:    getMovieNamesFromFile(),
	}, nil
}

func (mg *MovieGenerator) GetRandomMovie() (*Movie, error) {
	rint := rand.Intn(len(mg.movieList))
	var randomMovie string
	var err error
	if mg.movieFetcher.CanFetch {
		randomMovie, err = mg.movieFetcher.GetRandomMovieTitle()
		if err != nil {
			return nil, err
		}
	} else {
		randomMovie = mg.movieList[rint]
	}

	imageUrl, _ := mg.giphyPicker.GetRandomGiphy(strings.ReplaceAll(randomMovie, " ", "+"))

	movie := &Movie{
		Name:   randomMovie,
		GifURL: imageUrl,
	}

	return movie, nil
}

func getMovieNamesFromFile() []string {
	fdata, err := ioutil.ReadFile("movies.txt")
	if err != nil {
		fmt.Println("error reading movie txt", err)
		return []string{"The+Departed"}
	}

	return strings.Split(string(fdata), "\n")
}
