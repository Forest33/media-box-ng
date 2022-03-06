package entity

type MpvRequest struct {
	Payload []string `json:"command"`
}

type MpvResponse struct {
	Data      interface{} `json:"data"`
	RequestID int         `json:"request_id"`
	Error     string      `json:"error"`
}

type Channel struct {
	Title string `json:"title"`
	Img   string `json:"img"`
}
