package controller

//Paging api response object
type Paging struct {
	Page      int `json:"page"`
	TotalPage int `json:"total_page"`
}

//Response standard API response object
type Response struct {
	Paging  interface{} `json:"paging,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}
