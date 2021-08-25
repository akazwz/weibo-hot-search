package influx

import (
	"context"
	"fmt"
	"github.com/akazwz/weibo-hot-search/global"
	"github.com/akazwz/weibo-hot-search/model"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"log"
	"strconv"
	"time"
)

func GetCurrentHotSearch() (model.HotSearch, error) {
	client := influxdb2.NewClient(global.CFG.URL, global.CFG.Token)
	defer client.Close()
	stop := time.Now().Format(time.RFC3339)
	start := time.Now().Add(-15 * time.Minute).Format(time.RFC3339)
	query := `import "influxdata/influxdb/schema"
    from(bucket: "weibo")
    |> range(start: ` + start + `, stop: ` + stop + `)
    |> schema.fieldsAsCols()
    |> timeShift(duration: 8h, columns: ["_start", "_stop", "_time"])`
	queryAPI := client.QueryAPI(global.CFG.Org)
	result, err := queryAPI.Query(context.Background(), query)

	hotSearch := model.HotSearch{}
	searches := make([]model.SingleHotSearch, 0)

	if err == nil {
		for result.Next() {
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			if result.Record().Measurement() == "hot_search" {
				values := result.Record().Values()

				rank := values["rank"]
				content := values["content"]
				hot := values["hot"]
				link := values["link"]
				topicLead := values["topic_lead"]

				if rank == "00" {
					imageFile := values["image_file"]
					pdfFile := values["pdf_file"]
					timeInterface := values["_time"]
					timeStr := fmt.Sprintf("%v", timeInterface)
					timeStr = timeStr[:19]
					img := fmt.Sprintf("%v", imageFile)
					pdf := fmt.Sprintf("%v", pdfFile)
					hotSearch.Time = timeStr
					hotSearch.ImageFile = img
					hotSearch.PdfFile = pdf
				} else {
					rankStr := fmt.Sprintf("%v", rank)
					contentStr := fmt.Sprintf("%v", content)
					hotStr := fmt.Sprintf("%v", hot)
					linkStr := fmt.Sprintf("%v", link)
					topicLeadStr := fmt.Sprintf("%v", topicLead)

					rankInt, err := strconv.Atoi(rankStr)
					if err != nil {
						log.Println("rank conv error")
						return hotSearch, err
					}

					hotInt, err := strconv.Atoi(hotStr)

					if err != nil {
						log.Println("hot conv error")
						return hotSearch, err
					}

					singleHotSearch := model.SingleHotSearch{}
					singleHotSearch.Rank = rankInt
					singleHotSearch.Content = contentStr
					singleHotSearch.Hot = hotInt
					singleHotSearch.Link = linkStr
					singleHotSearch.TopicLead = topicLeadStr
					searches = append(searches, singleHotSearch)
				}
			}
		}
		hotSearch.Searches = searches
		if result.Err() != nil {
			fmt.Printf("query parsing error: %v\n", result.Err().Error())
			return hotSearch, result.Err()
		}
	} else {
		panic(err)
	}
	return hotSearch, nil
}

func GetDurationHotSearch(start, stop string) ([]model.HotSearch, error) {
	client := influxdb2.NewClient(global.CFG.URL, global.CFG.Token)
	defer client.Close()
	query := `import "influxdata/influxdb/schema"
    from(bucket: "weibo")
    |> range(start: ` + start + `, stop: ` + stop + `)
    |> schema.fieldsAsCols()
    |> timeShift(duration: 8h, columns: ["_start", "_stop", "_time"])
    |> group(columns: ["_time"], mode:"by")`
	queryAPI := client.QueryAPI(global.CFG.Org)
	result, err := queryAPI.Query(context.Background(), query)

	var hotSearches []model.HotSearch
	searches := make([]model.SingleHotSearch, 0)
	hotSearch := model.HotSearch{}
	if err == nil {
		for result.Next() {
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			values := result.Record().Values()
			table := values["table"]
			tableStr := fmt.Sprintf("%v", table)

			tableInt, err := strconv.Atoi(tableStr)
			if err != nil {
				log.Println("table conv error")
				return hotSearches, err
			}

			rank := values["rank"]
			content := values["content"]
			hot := values["hot"]
			link := values["link"]
			topicLead := values["topic_lead"]

			if rank == "00" {
				hotSearches = append(hotSearches, model.HotSearch{})
				imageFile := values["image_file"]
				pdfFile := values["pdf_file"]
				timeInterface := values["_time"]
				timeStr := fmt.Sprintf("%v", timeInterface)
				timeStr = timeStr[:19]
				img := fmt.Sprintf("%v", imageFile)
				pdf := fmt.Sprintf("%v", pdfFile)
				hotSearch.Time = timeStr
				hotSearch.ImageFile = img
				hotSearch.PdfFile = pdf
			} else {
				rankStr := fmt.Sprintf("%v", rank)
				contentStr := fmt.Sprintf("%v", content)
				hotStr := fmt.Sprintf("%v", hot)
				linkStr := fmt.Sprintf("%v", link)
				topicLeadStr := fmt.Sprintf("%v", topicLead)

				rankInt, err := strconv.Atoi(rankStr)
				if err != nil {
					log.Println("rank conv error")
					return hotSearches, err
				}

				hotInt, err := strconv.Atoi(hotStr)

				if err != nil {
					log.Println("hot conv error")
					return hotSearches, err
				}

				singleHotSearch := model.SingleHotSearch{}
				singleHotSearch.Rank = rankInt
				singleHotSearch.Content = contentStr
				singleHotSearch.Hot = hotInt
				singleHotSearch.Link = linkStr
				singleHotSearch.TopicLead = topicLeadStr
				searches = append(searches, singleHotSearch)
			}
			hotSearch.Searches = searches
			hotSearches[tableInt] = hotSearch
		}
		if result.Err() != nil {
			fmt.Printf("query parsing error: %s\n", result.Err().Error())
		}
	} else {
		panic(err)
	}
	return hotSearches, nil
}

func GetHotSearchesByContent(content, start, stop string) ([]model.HotSearch, error) {
	client := influxdb2.NewClient(global.CFG.URL, global.CFG.Token)
	defer client.Close()
	if start == "" || stop == "" {
		stop = time.Now().Format(time.RFC3339)
		start = time.Now().Add(-6 * time.Hour).Format(time.RFC3339)
	}
	query := `import "influxdata/influxdb/schema"
    from(bucket: "weibo")
    |> range(start: ` + start + `, stop: ` + stop + `)
    |> schema.fieldsAsCols()
    |> group(columns: ["time"], mode:"by")
    |> sort(columns: ["_time"])
    |> timeShift(duration: 8h, columns: ["_start", "_stop", "_time"])
    |> filter(fn: (r) => r["content"] == "` + content + `")`
	queryAPI := client.QueryAPI(global.CFG.Org)
	result, err := queryAPI.Query(context.Background(), query)
	var hotSearches []model.HotSearch
	if err == nil {
		for result.Next() {
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			values := result.Record().Values()

			timeInterface := values["_time"]
			timeStr := fmt.Sprintf("%v", timeInterface)
			timeStr = timeStr[:19]
			rank := values["rank"]
			contentInterface := values["content"]
			hot := values["hot"]
			link := values["link"]
			topicLead := values["topic_lead"]

			rankStr := fmt.Sprintf("%v", rank)
			contentStr := fmt.Sprintf("%v", contentInterface)
			hotStr := fmt.Sprintf("%v", hot)
			linkStr := fmt.Sprintf("%v", link)
			topicLeadStr := fmt.Sprintf("%v", topicLead)

			rankInt, err := strconv.Atoi(rankStr)
			if err != nil {
				log.Println("rank conv error")
				return hotSearches, err
			}

			hotInt, err := strconv.Atoi(hotStr)

			if err != nil {
				log.Println("hot conv error")
				return hotSearches, err
			}
			singleHotSearch := model.SingleHotSearch{}
			singleHotSearch.Rank = rankInt
			singleHotSearch.Content = contentStr
			singleHotSearch.Hot = hotInt
			singleHotSearch.Link = linkStr
			singleHotSearch.TopicLead = topicLeadStr

			hotSearches = append(hotSearches, model.HotSearch{
				Time:      timeStr,
				ImageFile: "",
				PdfFile:   "",
				Searches:  []model.SingleHotSearch{singleHotSearch},
			})
		}
		if result.Err() != nil {
			fmt.Printf("query parsing error: %s\n", result.Err().Error())
			return hotSearches, result.Err()
		}
	} else {
		panic(err)
		return hotSearches, err
	}
	return hotSearches, nil
}
