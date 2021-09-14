package model

type HotSearch struct {
	Time      string            `json:"time"`
	ImageFile string            `json:"image_file"`
	PdfFile   string            `json:"pdf_file"`
	Searches  []SingleHotSearch `json:"searches"`
}

type SingleHotSearch struct {
	Rank      int    `json:"rank"`
	Content   string `json:"content"`
	Tag       string `json:"tag"`
	Hot       int    `json:"hot"`
	Link      string `json:"link"`
	TopicLead string `json:"topic_lead"`
}
