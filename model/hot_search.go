package model

import "time"

type HotSearch struct {
	Time     time.Time         `json:"time"`
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
