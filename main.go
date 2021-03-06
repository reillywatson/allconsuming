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
		Title string
		Subtitle string
		Authors []string
		Description string
		ImageLinks struct {
			SmallThumbnail string
			Thumbnail string
		}
		CanonicalVolumeLink string
	}
}

type SearchResponse struct {
	Items []BookInfo
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
		response += fmt.Sprintf("<a href=\"%s\">", element.VolumeInfo.CanonicalVolumeLink)
		response += fmt.Sprintf("<img src=\"%s\"/>", element.VolumeInfo.ImageLinks.SmallThumbnail)
		response += "</a>"
		response += fmt.Sprintf("<p>%s<br/>", element.VolumeInfo.Title)
		authors := strings.Join(element.VolumeInfo.Authors, ", ")
		response += fmt.Sprintf("%s</p>", authors)
	}
	response += "</html>"
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
