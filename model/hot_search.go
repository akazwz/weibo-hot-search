package model

type KeywordHotSearch struct {
	Keyword     string      `json:"keyword"`
	HotSearches []HotSearch `json:"hot_searches"`
}

type HotSearch struct {
	Time     string            `json:"time"`
	Searches []SingleHotSearch `json:"searches"`
}

type SingleHotSearch struct {
	Rank    int    `json:"rank"`
	Content string `json:"content"`
	Link    string `json:"link"`
	Hot     int    `json:"hot"`
	Tag     string `json:"tag"`
	Icon    string `json:"icon"`
}
