package api

type InformationResponse struct {
	Category string `json:"category"`
	Info     string `json:"Information"`
}

type InformationCategoryResponse struct {
	Name string `json:"category"`
	URL  string `json:"url"`
}
