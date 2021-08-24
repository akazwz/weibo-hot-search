package influx

import (
	"context"
	"fmt"
	"github.com/akazwz/weibo-hot-search/global"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/http"
	"log"
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

func Query() {
	client := influxdb2.NewClient(global.CFG.URL, global.CFG.Token)
	// always close client at the end
	defer client.Close()
	stop := time.Now().Format(time.RFC3339)
	timeStart, err := time.Parse("2006-01-02 15:04:05", "2021-08-24 00:00:00")
	if err != nil {
		log.Fatalln("解析时间失败")
	}
	start := timeStart.Format(time.RFC3339)
	query := fmt.Sprintf("from(bucket:\"%v\")|> range(start: %v, stop: %v) |> filter(fn: (r) => r[\"_measurement\"] == \"hot-search\")",
		global.CFG.Bucket, start, stop)
	// Get query client
	queryAPI := client.QueryAPI(global.CFG.Org)
	// get QueryTableResult
	result, err := queryAPI.Query(context.Background(), query)

	fmt.Println(query)
	fmt.Println(result.Record())
	if err == nil {
		// Iterate over query response
		for result.Next() {
			// Notice when group key has changed
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			// Access data
			fmt.Printf("value: %v\n", result.Record().Value())
		}
		// check for an error
		if result.Err() != nil {
			fmt.Printf("query parsing error: %v\n", result.Err().Error())
		}
	} else {
		panic(err)
	}
}
