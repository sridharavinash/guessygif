package moviedb

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// Result represents the results returned inside from ResponseJSON
type Result struct {
	Title         string `json:"title"`
	PosterPath    string `json:"poster_path"`
	OriginalTitle string `json:"original_title"`
	ReleaseDate   string `json:"release_date"`
}

// ResponseJSON represents a response returned from themoviedb
type ResponseJSON struct {
	Page         int      `json:"page"`
	TotalResults int      `json:"total_results"`
	TotalPages   int      `json:"total_pages"`
	Results      []Result `json:"results"`
}

type MovieFetcher struct {
	ApiKey   string
	CanFetch bool
	client   *http.Client
}

func NewFetcher() (*MovieFetcher, error) {
	fetcher := &MovieFetcher{
		CanFetch: true,
		client:   &http.Client{Timeout: 10 * time.Second},
	}

	//Don't care if this is empty, will fallback to using list
	moviedb_api_key := os.Getenv("MOVIEDB_API_KEY")

	if moviedb_api_key == "" {
		fetcher.CanFetch = false
	}

	fetcher.ApiKey = moviedb_api_key

	return fetcher, nil
}

func (f *MovieFetcher) GetRandomMovieTitle() (string, error) {
	if !f.CanFetch {
		return "", fmt.Errorf("Cannot fetch movies from MovieDB")
	}

	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/top_rated?api_key=%s&region=US", f.ApiKey)

	resp := new(ResponseJSON)
	err := f.fetchFromApi(url, resp)
	if err != nil {
		fmt.Println("Occurred an error: ", err)
		return "", err
	}

	numMovie, pageNumber := randomMovieIndex(resp.TotalResults, len(resp.Results))
	newUrl := fmt.Sprintf("%s&page=%d", url, pageNumber)

	newResp := new(ResponseJSON)
	f.fetchFromApi(newUrl, newResp)

	result := newResp.Results[numMovie]

	return result.Title, nil
}

func randomMovieIndex(total, perPage int) (int, int) {
	numMovie := rand.Intn(total)
	pageNumber := numMovie / perPage
	numMovie = int(numMovie % perPage)
	return numMovie, pageNumber
}

func (f *MovieFetcher) fetchFromApi(url string, target interface{}) error {
	r, err := f.client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
