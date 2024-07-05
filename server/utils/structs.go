package utils

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
