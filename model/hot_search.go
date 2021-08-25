package model

type HotSearch struct {
	Time      string
	ImageFile string
	PdfFile   string
	Searches  []SingleHotSearch
}

type SingleHotSearch struct {
	Rank      int
	Content   string
	Hot       int
	Link      string
	TopicLead string
}
