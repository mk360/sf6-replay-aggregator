package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sf6-replays/utils"
	"strconv"

	"github.com/joho/godotenv"
)

const PAGE_SIZE int = 24

var characters [23]string = [23]string{"Luke", "Jamie", "Manon", "Kimberly", "Marisa", "Lily", "JP", "Juri", "Dee Jay", "Cammy", "Ryu", "E.Honda", "Blanka", "Guile", "Ken", "Chun-Li", "Zangief", "Dhalsim", "Rashid", "A.K.I.", "Ed", "Akuma", "M. Bison"}

func corsMiddleware(next http.Handler) http.Handler {
	var corsMiddleware = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Access-Control-Allow-Origin", os.Getenv("CORS_DOMAIN"))
		writer.Header().Add("Access-Control-Allow-Methods", "GET")
		next.ServeHTTP(writer, request)
	})

	return corsMiddleware
}

func contains(arr []string, name string) bool {
	for _, char := range arr {
		if char == name {
			return true
		}
	}
	return false
}

func handleRequest(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writer.WriteHeader(405)
		return
	}

	var character = request.URL.Query().Get("character")
	// var opponent = request.URL.Query().Get("opponent")

	var page = request.URL.Query().Get("page")

	if character != "" && !contains(characters[:], character) {
		writer.WriteHeader(400)
		writer.Write([]byte("Unknown character: " + character + "\n"))
		return
	}

	file, _ := os.ReadFile("./playlists/" + character + ".json")

	var jsonVideos []utils.JSONVideo = []utils.JSONVideo{}
	json.Unmarshal(file, &jsonVideos)

	pagination, e := strconv.Atoi(page)

	if e != nil {
		writer.WriteHeader(400)
		writer.Write([]byte("Invalid page param\n"))
		return
	}

	var minOffset = (pagination - 1) * PAGE_SIZE
	var maxOffset = pagination * PAGE_SIZE

	var subset = jsonVideos[minOffset:maxOffset]

	utils.RenderToHTML(writer, subset)
}

func main() {
	godotenv.Load("./env")
	mux := http.NewServeMux()

	var mainHandler = http.HandlerFunc(handleRequest)

	mux.Handle("/replays", corsMiddleware(mainHandler))

	http.ListenAndServe(":4444", mux)
}
