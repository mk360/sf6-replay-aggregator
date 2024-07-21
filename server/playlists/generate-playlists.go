package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type MultiPlaylistsMap struct {
	Data          map[string]string
	TransformFunc func(name string) string
}

type YouTubeItem struct {
	Title   string `json:"title"`
	Snippet struct {
		ResourceId struct {
			Id string `json:"videoId"`
		} `json:"resourceId"`
		Thumbnails struct {
			Medium struct {
				Url string `json:"url"`
			} `json:"medium"`
		} `json:"thumbnails"`
		Title string `json:"title"`
	} `json:"snippet"`
}

type YouTubeResponse struct {
	Items         []YouTubeItem `json:"items"`
	NextPageToken string        `json:"nextPageToken"`
}

type StoredVideo struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Thumbnail string `json:"img"`
}

func getChannelPlaylists(channel string) map[string]string {
	var playlists = make(map[string]string)
	fileData, _ := os.ReadFile("../channels/" + channel + ".json")
	json.Unmarshal(fileData, &playlists)

	return playlists
}

func sendRequest(playlistId string, pageToken string) YouTubeResponse {
	var queryParams url.Values = url.Values{}
	queryParams.Add("part", "snippet,status")
	queryParams.Add("maxResults", "50")
	queryParams.Add("key", os.Getenv("API_KEY"))
	queryParams.Add("playlistId", playlistId)

	if pageToken != "" {
		queryParams.Add("pageToken", pageToken)
	}

	var url = `https://www.googleapis.com/youtube/v3/playlistItems?` + queryParams.Encode()
	req, _ := http.NewRequest("GET", url, nil)
	var response *http.Response
	for {
		freshResponse, e := httpClient.Do(req)
		if e == nil {
			response = freshResponse
			break
		}
	}
	defer response.Body.Close()
	responseData, _ := io.ReadAll(response.Body)
	var playlistResponse YouTubeResponse = YouTubeResponse{}
	json.Unmarshal(responseData, &playlistResponse)

	return playlistResponse
}

var httpClient = &http.Client{}

func scrapeCharacterPlaylists(playlistId string, waitGroup *sync.WaitGroup, characterName string) {
	defer waitGroup.Done()
	var nextPageToken string = ""
	var results []StoredVideo = []StoredVideo{}
	for {
		var data = sendRequest(playlistId, nextPageToken)
		for _, item := range data.Items {
			var storedVideo StoredVideo = StoredVideo{
				Id:        item.Snippet.ResourceId.Id,
				Title:     item.Snippet.Title,
				Thumbnail: item.Snippet.Thumbnails.Medium.Url,
			}

			results = append(results, storedVideo)
		}
		nextPageToken = data.NextPageToken
		if nextPageToken == "" {
			break
		}
	}

	jsonified, _ := json.Marshal(results)

	os.WriteFile(characterName+".json", jsonified, 0644)
}

func main() {
	var start = time.Now()
	channelMap := make(map[string]MultiPlaylistsMap)
	godotenv.Load("./../.env")
	channelMap["The FGC Place"] = MultiPlaylistsMap{
		Data: getChannelPlaylists("The FGC Place"),
		TransformFunc: func(name string) string {
			return "SF6 🔥 " + name
		},
	}

	channelMap["SF6 High Level Replays"] = MultiPlaylistsMap{
		Data: getChannelPlaylists("SF6 High Level Replays"),
		TransformFunc: func(name string) string {
			return name + " ▰ high level gameplay [SF6]"
		},
	}

	var characters [23]string = [23]string{"Ken", "Ryu", "Chun-Li", "Ed", "A.K.I.", "Dee Jay", "M. Bison", "Jamie", "Guile", "JP", "Juri", "Akuma", "Blanka", "Cammy", "Dhalsim", "Kimberly", "Lily", "Luke", "Manon", "Marisa", "Rashid", "Zangief", "E.Honda"}

	var waitGroup = sync.WaitGroup{}
	for _, character := range characters {
		var playlistIds [2]string = [2]string{"", ""}
		var i int = 0
		for channel := range channelMap {
			var characterPlaylistName = channelMap[channel].TransformFunc(character)
			var playlistId = channelMap[channel].Data[characterPlaylistName]
			playlistIds[i] = playlistId
			i++
		}

		for _, playlistId := range playlistIds {
			waitGroup.Add(1)
			go scrapeCharacterPlaylists(playlistId, &waitGroup, character)
		}
	}

	waitGroup.Wait()

	fmt.Println(time.Since(start))
}
