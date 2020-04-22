package giphy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

type GiphyPicker struct {
	ApiKey string
}

type GiphyDownSizedLarge struct {
	URL string `json:"url"`
}

type GiphyImageEntry struct {
	Image GiphyDownSizedLarge `json:"downsized_large"`
}

type GiphyDataItem struct {
	Images GiphyImageEntry `json:"images"`
}

type GiphyAPIResponse struct {
	Data []GiphyDataItem `json:"data"`
}

func NewPicker() (*GiphyPicker, error) {
	giphy_api_key := os.Getenv("GIPHY_API_KEY")
	if giphy_api_key == "" {
		return nil, fmt.Errorf("Missing Giphy API Key")
	}

	return &GiphyPicker{
		ApiKey: giphy_api_key,
	}, nil
}

func (p *GiphyPicker) GetRandomGiphy(s string) (string, error) {
	roffset := rand.Intn(3)
	url := fmt.Sprintf("https://api.giphy.com/v1/gifs/search?api_key=%s&q=%s&limit=1&offset=%d&rating=G&lang=en", p.ApiKey, s, roffset)

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

	var apiResponse GiphyAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		fmt.Println("error unmarshalling", err)
		return "", err
	}

	resp := apiResponse.Data
	imageUrl := resp[0].Images.Image.URL
	return imageUrl, nil
}
