package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// store playlist ids in file, sort them by channel

type ThumbnailsStruct struct {
	MaxRes struct {
		Url    string `json:"url"`
		Width  int16  `json:"width"`
		Height int16  `json:"height"`
	} `json:"maxres"`
}

type SnippetStruct struct {
	Snippet struct {
		Thumbnails ThumbnailsStruct `json:"thumbnails"`
		Title      string           `json:"title"`
	} `json:"snippet"`
}

type VideoResponse struct {
	NextPageToken string          `json:"nextPageToken"`
	Items         []SnippetStruct `json:"items"`
}

type YouTubePlaylistResponse struct {
	Kind     string `json:"kind"`
	PageInfo struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	NextPageToken string          `json:"nextPageToken"`
	Items         []SnippetStruct `json:"items"`
}

var httpClient = &http.Client{}

// func getPlaylists(channelId string, API_KEY string, pageToken string, playlistResponses map[string]string) map[string]string {
// 	query := url.Values{
// 		"key":        {API_KEY},
// 		"part":       {"snippet"},
// 		"channelId":  {channelId},
// 		"maxResults": {"50"},
// 	}
// 	if pageToken != "" {
// 		query.Add("pageToken", pageToken)
// 	}
// 	var playlistResponse YouTubePlaylistResponse
// 	request, _ := http.NewRequest("GET", "https://www.googleapis.com/youtube/v3/playlists?"+query.Encode(), nil)
// 	request.Header.Add("User-Agent", "axios/1.6.8")
// 	request.Header.Add("Host", "www.googleapis.com")
// 	request.Header.Add("Accept", "application/json, text/plain, */*")

// 	responseBody := sendRequest(request)

// 	json.Unmarshal(responseBody, &playlistResponse)

// 	for _, snippet := range playlistResponse.Items {
// 		if strings.Contains(snippet.Snippet.Title, "SF6") {
// 			playlistResponses[snippet.Snippet.Title] = snippet.ID
// 		}
// 	}

// 	if playlistResponse.NextPageToken == "" {
// 		return playlistResponses
// 	}

// 	return getPlaylists(channelId, API_KEY, playlistResponse.NextPageToken, playlistResponses)
// }

var characters [23]string = [23]string{"Luke", "Jamie", "Manon", "Kimberly", "Marisa", "Lily", "JP", "Juri", "Dee Jay", "Cammy", "Ryu", "E. Honda", "Blanka", "Guile", "Ken", "Chun-Li", "Zangief", "Dhalsim", "Rashid", "A.K.I", "Ed", "Akuma", "Bison"}

type ChannelMapStruct struct {
	Name                   string
	CharacterNameConverter func(string) string
}

func sendRequest(request *http.Request) []byte {
	var response http.Response
	for {
		tempResponse, e := httpClient.Do(request)
		if e == nil {
			response = *tempResponse
			defer response.Body.Close()
			break
		}
	}

	data, _ := io.ReadAll(response.Body)

	return data
}

func corsMiddleware(next http.Handler) http.Handler {
	var corsMiddleware = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Access-Control-Allow-Origin", "*")
		writer.Header().Add("Access-Control-Allow-Methods", "GET")
		next.ServeHTTP(writer, request)
	})

	return corsMiddleware
}

// when a character is requested
// create a go routine for each channel
// in each goroutine
// find character playlist
// store next page pagination key in a map
// then find videos
// send the videos to a channel

func contains(arr []string, name string) bool {
	for _, char := range arr {
		if char == name {
			return true
		}
	}
	return false
}

func main() {
	godotenv.Load()
	mux := http.NewServeMux()

	ChannelMap := make(map[string]ChannelMapStruct)
	ChannelMap["UCi5rlUH3C4BzDB5-fRJ8hHg"] = ChannelMapStruct{Name: "SF6 High Level Replays", CharacterNameConverter: func(name string) string {
		return name + " ▰ high level gameplay [SF6]"
	}}

	ChannelMap["UCx2dkBZglt1xlVMbzb63uCQ"] = ChannelMapStruct{Name: "The FGC Place", CharacterNameConverter: func(name string) string {
		return "SF6 🔥 " + name
	}}

	var playlists map[string]map[string]string = make(map[string]map[string]string)
	dir, _ := os.ReadDir("./playlists")
	for _, fileName := range dir {
		var channelPlaylists map[string]string = make(map[string]string)
		var filename = fileName.Name()
		channel, _ := os.ReadFile("./playlists/" + filename)
		var withoutJSON = strings.Replace(filename, ".json", "", 1)
		json.Unmarshal(channel, &channelPlaylists)
		playlists[withoutJSON] = channelPlaylists
	}

	var mainHandler = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writer.WriteHeader(405)
			return
		}

		var character = request.URL.Query().Get("character")
		// var opponent = request.URL.Query().Get("opponent")
		var channelId = request.URL.Query().Get("channelId")
		// var page = request.URL.Query().Get("page")

		if character != "" && !contains(characters[:], character) {
			writer.WriteHeader(400)
			writer.Write([]byte("Unknown character: " + character + "\n"))
			return
		}
		// if opponent != "" && !contains(characters[:], opponent) {
		// 	writer.WriteHeader(400)
		// 	return
		// }

		channelName, ok := ChannelMap[channelId]
		if channelId != "" && !ok {
			writer.WriteHeader(400)
			writer.Write([]byte("Unknown channel: " + channelId + "\n"))
			return
		}

		var API_KEY = os.Getenv("API_KEY")
		var query = url.Values{
			"part":       {"snippet"},
			"maxResults": {"20"},
			"key":        {API_KEY},
		}

		if character != "" && channelId != "" {
			var transformed = channelName.CharacterNameConverter(character)
			var playlistName = channelName.Name
			playlistContent, ok := playlists[playlistName]

			if !ok {
				// shouldn't be there
				log.Fatalln(transformed, playlistName)
			}

			characterPlaylistId, ok := playlistContent[transformed]

			if !ok {
				// shouldn't reach this point
				log.Fatalln("can't find playlist name for character " + transformed)
			}

			query.Add("playlistId", characterPlaylistId)
		} else if character != "" {
			// get from both channels
			fmt.Println("bojnour")
		}

		var resp = getVideos(query)
		marshaled, _ := json.Marshal(resp)
		writer.Write(marshaled)
		// rassembler toutes les vidéos de toutes les chaînes. Ajouter un filtre par nom de chaîne
	})

	mux.Handle("/replays", corsMiddleware(mainHandler))

	http.ListenAndServe(":4444", mux)
}

func getVideos(queryParams url.Values) VideoResponse {
	var url = "https://www.googleapis.com/youtube/v3/playlistItems?" + queryParams.Encode()
	var videosResponse VideoResponse
	request, _ := http.NewRequest("GET", url, nil)
	data := sendRequest(request)
	json.Unmarshal(data, &videosResponse)
	return videosResponse
}
