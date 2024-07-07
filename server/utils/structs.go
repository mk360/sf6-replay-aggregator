package utils

type ThumbnailsStruct struct {
	Medium struct {
		Url string `json:"url"`
	} `json:"medium"`
}

type SnippetStruct struct {
	Snippet struct {
		Thumbnails ThumbnailsStruct `json:"thumbnails"`
		Title      string           `json:"title"`
		ResourceId struct {
			VideoId string `json:"videoId"`
		} `json:"resourceId"`
	} `json:"snippet"`
}

type VideosResponse struct {
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
