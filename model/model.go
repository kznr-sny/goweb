package model

type RequestData struct {
	IsPost bool `json:"isPost"`
	Urls []UrlsData `json:"urls"`
	Params []ParamsData `json:"params"`
}

type UrlsData struct {
	ID int `json:"id"`
	Url string `json:"url"`
}

type ParamsData struct {
	ID int `json:"id"`
	Key string `json:"key"`
	Value string `json:"value"`
}

type ParamData struct {
	Key string
	Value string
}