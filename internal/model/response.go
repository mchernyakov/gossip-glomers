package model

type SimpleResp struct {
	Type string `json:"type"`
}

type KafkaSendResp struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
}
