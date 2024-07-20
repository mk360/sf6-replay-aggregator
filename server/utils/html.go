package utils

import (
	"net/http"
	"os"
	stdTemplate "text/template"
)

var templateData, _ = os.ReadFile("./utils/video-template.html")
var template, _ = stdTemplate.New("video-template").Parse(string(templateData))

func RenderToHTML(responseWriter http.ResponseWriter, data []JSONVideo) {
	for _, item := range data {
		template.Execute(responseWriter, item)
	}
}
