package influx

import (
	"context"
	"fmt"
	"github.com/akazwz/weibo-hot-search/global"
	"github.com/akazwz/weibo-hot-search/model"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/http"
	"log"
	"strconv"
	"time"
)

func Write(measurement string, tags map[string]string, fields map[string]interface{}) (err error) {
	client := influxdb2.NewClient(global.CFG.URL, global.CFG.Token)
	// always close client at the end
	defer client.Close()
	p := influxdb2.NewPoint(measurement, tags, fields, time.Now())
	writeApi := client.WriteAPI(global.CFG.Org, global.CFG.Bucket)
	writeApi.WritePoint(p)
	writeApi.Flush()
	writeApi.SetWriteFailedCallback(func(batch string, error http.Error, retryAttempts uint) bool {
		err = &error
		return false
	})
	return
}

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
					//timeTime, err := time.Parse(time., timeStr)
					/*location, err := time.LoadLocation("Asia/Shanghai")
					if err != nil {
						log.Println("time location load error")
						return hotSearch, err
					}
					timeTime, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, location)
					if err != nil {
						log.Println("time parse error")
						return hotSearch, err
					}*/
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

func GetDurationHotSearch() {
	client := influxdb2.NewClient(global.CFG.URL, global.CFG.Token)
	defer client.Close()
	stop := time.Now().Format(time.RFC3339)
	start := time.Now().Add(-90 * time.Minute).Format(time.RFC3339)
	query := `import "influxdata/influxdb/schema"
    from(bucket: "weibo")
    |> range(start: ` + start + `, stop: ` + stop + `)
    |> schema.fieldsAsCols()
    |> timeShift(duration: 8h, columns: ["_start", "_stop", "_time"])`
	queryAPI := client.QueryAPI(global.CFG.Org)
	result, err := queryAPI.Query(context.Background(), query)

	if err == nil {
		// Iterate over query response
		for result.Next() {
			// Notice when group key has changed
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			values := result.Record().Values()
			fmt.Println(values)
		}
		// check for an error
		if result.Err() != nil {
			fmt.Printf("query parsing error: %s\n", result.Err().Error())
		}
	} else {
		panic(err)
	}

}
