package model

type RequestData struct {
	IsPost bool `json:"isPost"`
	Uris []UrisData `json:"uris"`
	Params []ParamsData `json:"params"`
}

type UrisData struct {
	ID int `json:"id"`
	Uri string `json:"uri"`
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