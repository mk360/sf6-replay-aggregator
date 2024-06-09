package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
)

// store playlist ids in file, sort them by channel

type ThumbnailsStruct struct {
	MaxRes struct {
		Url    string `json:"url"`
		Width  int16  `json:"width"`
		Height int16  `json:"height"`
	} `json:"maxres"`
}

type YouTubePlaylistResponse struct {
	Kind     string `json:"kind"`
	PageInfo struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []struct {
		ID      string `json:"id"`
		Snippet struct {
			Thumbnails ThumbnailsStruct `json:"thumbnails"`
			Title      string           `json:"title"`
		} `json:"snippet"`
	} `json:"items"`
}

var API_KEY = os.Getenv("API_KEY")

func sendYoutubeApiRequest(channelId string, httpClient *http.Client) YouTubePlaylistResponse {
	query := url.Values{
		"key":        {API_KEY},
		"part":       {"snippet"},
		"channelId":  {channelId},
		"maxResults": {"50"},
	}
	var playlistResponse YouTubePlaylistResponse
	request, _ := http.NewRequest("GET", "https://www.googleapis.com/youtube/v3/playlists?"+query.Encode(), nil)
	request.Header.Add("User-Agent", "axios/1.6.8")
	request.Header.Add("Host", "www.googleapis.com")
	request.Header.Add("Accept", "application/json, text/plain, */*")

	var response http.Response
	for {
		tempResponse, e := httpClient.Do(request)
		if e == nil {
			response = *tempResponse
			defer response.Body.Close()
			break
		}
	}

	responseBody, _ := io.ReadAll(response.Body)
	json.Unmarshal(responseBody, &playlistResponse)

	return playlistResponse
}

type ChannelMapStruct struct {
	Name                   string
	CharacterNameConverter func(string) string
}

// var playlistByteData map[string][]byte

// var channelNames []string = []string{"SF6 High Level Replays", "The FGC Place"}

func main() {
	mux := http.NewServeMux()

	ChannelMap := make(map[string]ChannelMapStruct)
	ChannelMap["UCi5rlUH3C4BzDB5-fRJ8hHg"] = ChannelMapStruct{Name: "SF6 High Level Replays", CharacterNameConverter: func(name string) string {
		return name + " ▰ high level gameplay [SF6]"
	}}
	ChannelMap["UCx2dkBZglt1xlVMbzb63uCQ"] = ChannelMapStruct{Name: "The FGC Place", CharacterNameConverter: func(name string) string {
		return "SF6 🔥 " + name
	}}

	httpClient := &http.Client{}
	sendYoutubeApiRequest("UCx2dkBZglt1xlVMbzb63uCQ", httpClient)
	var playlists map[string]string

	var a = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var character = request.PathValue("character")
		var channelId = request.URL.Query().Get("channelId")
		channelName, ok := ChannelMap[channelId]
		var transformed = channelName.CharacterNameConverter(character)
		var characterPlaylistId string = ""

		if ok {
			fileData, _ := os.ReadFile(channelName.Name + ".json")
			json.Unmarshal(fileData, &playlists)
			for characterPlaylist, playlistId := range playlists {
				if characterPlaylist == transformed {
					characterPlaylistId = playlistId
					break
				}
			}

			if characterPlaylistId != "" {
				var videos = getVideos(characterPlaylistId, *httpClient)
				videoData, _ := json.Marshal(videos.Items)
				writer.Write(videoData)
			}
		} else {
			writer.Write(nil)
		}
		// rassembler toutes les vidéos de toutes les chaînes. Ajouter un filtre par nom de chaîne
	})

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Hello world\n"))
	})

	http.ListenAndServe(":4444", nil)
}

func getVideos(playlistId string, httpClient http.Client) YouTubePlaylistResponse {
	var videosResponse YouTubePlaylistResponse
	query := url.Values{
		"key":        {API_KEY},
		"part":       {"snippet"},
		"playlistId": {playlistId},
		"maxResults": {"50"},
	}
	var a = "https://www.googleapis.com/youtube/v3/playlistItems?" + query.Encode()
	request, _ := http.NewRequest("GET", a, nil)
	response, _ := httpClient.Do(request)
	data, _ := io.ReadAll(response.Body)
	defer response.Body.Close()
	json.Unmarshal(data, &videosResponse)
	return videosResponse
}
