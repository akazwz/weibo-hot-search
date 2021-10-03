package influx

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/akazwz/weibo-hot-search/global"
	"github.com/akazwz/weibo-hot-search/model"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// GetCurrentHotSearch 获取当前热搜
func GetCurrentHotSearch() (model.HotSearch, error) {
	client := influxdb2.NewClient(global.CFG.URL, global.CFG.Token)
	defer client.Close()
	stop := time.Now().Format(time.RFC3339)
	// 获取两分钟之内的热搜，整分情况下有一条，其他情况为两条，最新的一条覆盖掉前一条
	start := time.Now().Add(-2 * time.Minute).Format(time.RFC3339)

	query := `import "influxdata/influxdb/schema"
    from(bucket: "weibo")
    |> range(start: ` + start + `, stop: ` + stop + `)
	|> filter(fn: (r) => r["_measurement"] == "new-hot")
    |> schema.fieldsAsCols()
    |> timeShift(duration: 8h, columns: ["_start", "_stop", "_time"])
	|> group(columns: ["_time"])`

	queryAPI := client.QueryAPI(global.CFG.Org)
	result, err := queryAPI.Query(context.Background(), query)

	hotSearch := model.HotSearch{}
	searches := make([]model.SingleHotSearch, 0)
	if err == nil {
		tableIndex := 0
		for result.Next() {
			if result.Record().Measurement() == "new-hot" {
				values := result.Record().Values()
				table := values["table"] // table,table相同时为同一条热搜
				tableStr := fmt.Sprintf("%v", table)
				tableInt, err := strconv.Atoi(tableStr)
				if err != nil {
					log.Println("table conv error")
				}
				if tableInt != tableIndex {
					// 有一次以上热搜时,重置searches
					searches = make([]model.SingleHotSearch, 0)
					tableIndex = tableInt
				}
				timeStr := fmt.Sprintf("%v", values["_time"])
				timeStr = timeStr[:19]
				rankStr := fmt.Sprintf("%v", values["rank"])
				// 日期在热搜第一条获取即可
				if rankStr == "01" {
					hotSearch.Time = timeStr
				}
				rankInt, err := strconv.Atoi(rankStr)
				if err != nil {
					log.Println("rank conv error")
				}
				contentStr := fmt.Sprintf("%v", values["content"])
				LinkStr := fmt.Sprintf("%v", values["link"])
				hotStr := fmt.Sprintf("%v", values["hot"])
				hotInt, err := strconv.Atoi(hotStr)
				if err != nil {
					log.Println("hot conv error")
				}

				// 为空 置零
				tagStr := fmt.Sprintf("%v", values["tag"])
				if values["tag"] == nil {
					tagStr = ""
				}
				iconStr := fmt.Sprintf("%v", values["icon"])
				if values["icon"] == nil {
					iconStr = ""
				}

				singleHotSearch := model.SingleHotSearch{
					Rank:    rankInt,
					Content: contentStr,
					Link:    LinkStr,
					Hot:     hotInt,
					Tag:     tagStr,
					Icon:    iconStr,
				}
				searches = append(searches, singleHotSearch)
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

/*func GetDurationHotSearch(start, stop string) ([]model.HotSearch, error) {
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

	hotSearches := make([]model.HotSearch, 0)
	searches := make([]model.SingleHotSearch, 0)
	hotSearch := model.HotSearch{}
	if err == nil {
		tableIndex := 0
		index := 0
		for result.Next() {
			index++
			// 获取包含热搜数据的map
			values := result.Record().Values()
			fmt.Println("values:", values)

			table := values["table"] // table,table相同时为同一条热搜
			tableStr := fmt.Sprintf("%v", table)
			tableInt, err := strconv.Atoi(tableStr)
			if err != nil {
				log.Println("table conv error")
				return hotSearches, err
			}
			if tableInt != tableIndex {
				// 热搜切换时,重置searches
				searches = make([]model.SingleHotSearch, 0)
				tableIndex = tableInt
			}
			rank := values["rank"]
			content := values["content"]
			hot := values["hot"]
			link := values["link"]
			topicLead := values["topic_lead"]
			tagStr := ""

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
				hotArr := strings.Split(hotStr, " ")
				if len(hotArr) > 1 {
					hotStr = hotArr[1]
					tagStr = hotArr[0]
				}
				linkStr := fmt.Sprintf("%v", link)
				topicLeadStr := fmt.Sprintf("%v", topicLead)
				if topicLead == nil {
					topicLeadStr = ""
				}

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
				singleHotSearch.Tag = tagStr
				singleHotSearch.Link = linkStr
				singleHotSearch.TopicLead = topicLeadStr
				searches = append(searches, singleHotSearch)
			}
			hotSearch.Searches = searches
			//hotSearches[tableIndex] = hotSearch
		}
		fmt.Println(index)
		if result.Err() != nil {
			fmt.Printf("query parsing error: %s\n", result.Err().Error())
		}
	} else {
		panic(err)
	}
	return hotSearches, nil
}*/

func GetHotSearchesByContent(content, start, stop string) ([]model.HotSearch, error) {
	client := influxdb2.NewClient(global.CFG.URL, global.CFG.Token)
	defer client.Close()
	if start == "" || stop == "" {
		stop = time.Now().Format(time.RFC3339)
		start = time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
	}
	query := `import "influxdata/influxdb/schema"
    from(bucket: "weibo")
    |> range(start: ` + start + `, stop: ` + stop + `)
    |> schema.fieldsAsCols()
    |> group(columns: ["time"], mode:"by")
    |> sort(columns: ["_time"])
    |> timeShift(duration: 8h, columns: ["_start", "_stop", "_time"])
    |> filter(fn: (r) => r.content == "` + content + `")`
	queryAPI := client.QueryAPI(global.CFG.Org)
	result, err := queryAPI.Query(context.Background(), query)
	var hotSearches []model.HotSearch
	if err == nil {
		for result.Next() {
			// table changed
			if result.TableChanged() {
			}
			values := result.Record().Values()

			timeStr := fmt.Sprintf("%v", values["_time"])
			timeStr = timeStr[:19]
			rankStr := fmt.Sprintf("%v", values["rank"])
			rankInt, err := strconv.Atoi(rankStr)
			if err != nil {
				log.Println("rank conv error")
			}
			contentStr := fmt.Sprintf("%v", values["content"])
			LinkStr := fmt.Sprintf("%v", values["link"])
			hotStr := fmt.Sprintf("%v", values["hot"])
			hotInt, err := strconv.Atoi(hotStr)
			if err != nil {
				log.Println("hot conv error")
			}

			// 为空 置零
			tagStr := fmt.Sprintf("%v", values["tag"])
			if values["tag"] == nil {
				tagStr = ""
			}
			iconStr := fmt.Sprintf("%v", values["icon"])
			if values["icon"] == nil {
				iconStr = ""
			}
			singleHotSearch := model.SingleHotSearch{
				Rank:    rankInt,
				Content: contentStr,
				Link:    LinkStr,
				Hot:     hotInt,
				Tag:     tagStr,
				Icon:    iconStr,
			}

			hotSearches = append(hotSearches, model.HotSearch{
				Time:     timeStr,
				Searches: []model.SingleHotSearch{singleHotSearch},
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
