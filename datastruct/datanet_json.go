package datastruct

//DatanetRequest is a request param json
type DatanetRequest struct {
	Filename []string `json:"filename"`
}

//DatanetResponse is a response json
type DatanetResponse struct {
	Status string                 `json:"status"`
	Data   map[string]interface{} `json:"data"`
}
