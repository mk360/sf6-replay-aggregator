package utils

import (
	"net/http"
	"os"
	stdTemplate "text/template"
)

var templateData, _ = os.ReadFile("./utils/video-template.html")
var template, _ = stdTemplate.New("video-template").Parse(string(templateData))

func RenderToHTML(responseWriter http.ResponseWriter, data VideosResponse) {
	for _, item := range data.Items {
		template.Execute(responseWriter, item)
	}
}
