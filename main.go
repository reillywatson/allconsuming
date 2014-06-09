package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"github.com/go-martini/martini"
)

type BookInfo struct {
	Id string `json::id`
	VolumeInfo struct {
		Title string `json::Title`
		Subtitle string `json::subtitle`
		Authors []string `json::authors`
		Description string `json::description`
		ImageLinks struct {
			SmallThumbnail string `json::smallThumbnail`
			Thumbnail string `json::thumbnail`
		} `json::imageLinks`
		CanonicalVolumeLink string `json::canonicalVolumeLink`
	} `json::volumeInfo`
}

type SearchResponse struct {
	Items []BookInfo `json::items`
}

func search(req *http.Request) string {
	apiUrl, _ := url.Parse("https://www.googleapis.com/books/v1/volumes")
	params := url.Values{}
	params.Set("q", req.URL.Query().Get("q"))
	apiUrl.RawQuery = params.Encode()
	println(apiUrl.String())
	resp, err := http.Get(apiUrl.String())
	if err != nil {
		return "google failed!"
	}
	var m SearchResponse
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &m)
	println(len(m.Items))
	if err != nil {
		panic("whaaaaaaaaa!")
	}
	var response = "<html>"
	for _, element := range m.Items {
		title := element.VolumeInfo.Title
		authors := strings.Join(element.VolumeInfo.Authors, ", ")
		response = response + fmt.Sprintf("<p>title: %s<br/>", title)
		response = response + fmt.Sprintf("author: %s</p>", authors)
	}
	response = response + "</html>"
	return response
}

func main() {
	m := martini.Classic()
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Get("/search", search)
	m.Run()
}
