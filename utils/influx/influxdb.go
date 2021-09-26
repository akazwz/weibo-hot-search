package influx

import (
	"context"
	"fmt"
	"github.com/akazwz/weibo-hot-search/global"
	"github.com/akazwz/weibo-hot-search/model"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"log"
	"time"
)

// GetCurrentHotSearch 获取当前热搜
func GetCurrentHotSearch() (model.HotSearch, error) {
	client := influxdb2.NewClient(global.CFG.URL, global.CFG.Token)
	defer client.Close()
	stop := time.Now().Format(time.RFC3339)
	start := time.Now().Add(-1 * time.Minute).Format(time.RFC3339)
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
		location, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			log.Fatal("时区加载失败")
		}
		for result.Next() {
			if result.TableChanged() {
				//fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			if result.Record().Measurement() == "new-hot" {
				values := result.Record().Values()
				//fmt.Println(values)
				timeStr := fmt.Sprintf("%v", values["_time"])
				timeStr = timeStr[:19]
				timeInLocation, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, location)
				formatTime := timeInLocation.Format("2006-01-02-15-04-05")
				if err != nil {
					log.Println("时间转换失败")
				}
				fmt.Println(formatTime)
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

/*func GetHotSearchesByContent(content, start, stop string) ([]model.HotSearch, error) {
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

			timeInterface := values["_time"]
			timeStr := fmt.Sprintf("%v", timeInterface)
			timeStr = timeStr[:19]
			rank := values["rank"]
			contentInterface := values["content"]
			hot := values["hot"]
			tagStr := ""
			link := values["link"]
			topicLead := values["topic_lead"]

			rankStr := fmt.Sprintf("%v", rank)
			contentStr := fmt.Sprintf("%v", contentInterface)
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
*/
