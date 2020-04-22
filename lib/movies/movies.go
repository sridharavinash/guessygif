package movies

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	"github.com/sridharavinash/guessygif/lib/giphy"
	"github.com/sridharavinash/guessygif/lib/moviedb"
)

type MovieGenerator struct {
	giphyPicker  *giphy.GiphyPicker
	movieFetcher *moviedb.MovieFetcher
	movieList    []string
	rnd          *rand.Rand
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
		rnd:          rand.New(rand.NewSource(time.Now().UnixNano())),
		giphyPicker:  picker,
		movieFetcher: fetcher,
		movieList:    getMovieNamesFromFile(),
	}, nil
}

func (mg *MovieGenerator) Intn(n int) int {
	return mg.rnd.Intn(n)
}

func (mg *MovieGenerator) GetRandomMovie() string {
	var randomMovie string
	var err error
	if mg.movieFetcher.CanFetch {
		randomMovie, err = mg.movieFetcher.GetRandomMovieTitle()
		if err != nil {
			randomMovie = ""
			fmt.Printf("\nAn error occurred: %+v\n", err)
		}
	}
	if randomMovie == "" {
		randomMovie = mg.movieList[mg.Intn(len(mg.movieList))]
	}
	return randomMovie
}

func (mg *MovieGenerator) GetMovieGif(title string) (string, error) {
	return mg.giphyPicker.GetRandomGiphy(strings.ReplaceAll(title, " ", "+"))
}

func getMovieNamesFromFile() []string {
	fdata, err := ioutil.ReadFile("movies.txt")
	if err != nil {
		fmt.Println("error reading movie txt", err)
		return []string{"The+Departed"}
	}

	return strings.Split(string(fdata), "\n")
}
