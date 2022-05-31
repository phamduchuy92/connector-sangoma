package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/olivere/elastic/v7"
)

func main() {
	var (
		elasticsearch_url   = *flag.String("elasticsearch_url", "http://localhost:9200", "Elasticsearch URL")
		elasticsearch_sniff = *flag.Bool("elasticsearch_sniff", true, "Default is enabled")
		// datetime            = *flag.String("datetime", time.Now().Format("2006-01-02"), "Datetime to calculate report")
		index = *flag.String("index", "call_log*", "Datetime to calculate report")
		// path_log            = *flag.String("path_log", "/var/log/spoofing-report/report", "Location for logs")
	)

	flag.Parse()
	// datetime = "2021-11-01"
	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	// logFileName := path_log + "-" + time.Now().Format("2006.01.02")
	// f, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error open file: " + logFileName)
	// } else {
	// 	defer f.Close()
	// 	log.SetOutput(f)
	// }

	// Create an Elasticsearch client
	elasticClient, err := elastic.NewClient(elastic.SetURL(strings.Split(elasticsearch_url, ",")...), elastic.SetSniff(elasticsearch_sniff))
	if err != nil {
		log.Fatalf("error connect to elasticsearch %v", err)
	}

	// Search with a term query
	// from, err := time.Parse("2006-01-02", datetime)
	// if err != nil {
	// 	log.Fatalf("error parse datetime %v", err)
	// }
	// to := from.Add(24 * time.Hour)
	// timeRange := elastic.NewRangeQuery("@timestamp").Gte(from.Format("2006-01-02T15:04:05")).Lte(to.Format("2006-01-02T15:04:05"))
	timeRange := elastic.NewRangeQuery("@timestamp").Gte("2021-09-06T00:00:00").Lte("2021-10-05T00:00:00")

	boolFilter := elastic.NewBoolQuery().Must(timeRange).MustNot(elastic.NewExistsQuery("jitter"))

	searchTotal, err := elasticClient.Search().
		Index(index).             // search in index "tweets"
		Query(boolFilter).        // specify the query
		Sort("@timestamp", true). // sort by "user" field, ascending
		From(0).Size(10000).      // take documents 0-9
		TrackTotalHits(true).
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		log.Fatalf("error get call_log - %v", err)
	}
	log.Printf("total %v", searchTotal.TotalHits())

	// 100
	term100 := elastic.NewTermQuery("error_code", 100)
	search100, err := elasticClient.Search().
		Index(index).                                                                                          // search in index "tweets"
		Query(elastic.NewBoolQuery().Must(timeRange).Must(term100).MustNot(elastic.NewExistsQuery("jitter"))). // specify the query
		Sort("@timestamp", true).                                                                              // sort by "user" field, ascending
		From(0).Size(10000).                                                                                   // take documents 0-9
		TrackTotalHits(true).
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		log.Fatalf("error get call_log - %v", err)
	}
	log.Printf("100 %v", search100.TotalHits())

	// // 200
	// term200 := elastic.NewTermQuery("error_code", 200)
	// search200, err := elasticClient.Search().
	// 	Index(index).                                                                                          // search in index "tweets"
	// 	Query(elastic.NewBoolQuery().Must(timeRange).Must(term200).MustNot(elastic.NewExistsQuery("jitter"))). // specify the query
	// 	Sort("@timestamp", true).                                                                              // sort by "user" field, ascending
	// 	From(0).Size(10000).                                                                                   // take documents 0-9
	// 	TrackTotalHits(true).
	// 	Pretty(true).            // pretty print request and response JSON
	// 	Do(context.Background()) // execute
	// if err != nil {
	// 	log.Fatalf("error get call_log - %v", err)
	// }
	// log.Printf("200 %v", search200.TotalHits())

	// 201
	term201 := elastic.NewTermQuery("error_code", 201)
	search201, err := elasticClient.Search().
		Index(index).                                                                                          // search in index "tweets"
		Query(elastic.NewBoolQuery().Must(timeRange).Must(term201).MustNot(elastic.NewExistsQuery("jitter"))). // specify the query
		Sort("@timestamp", true).                                                                              // sort by "user" field, ascending
		From(0).Size(10000).                                                                                   // take documents 0-9
		TrackTotalHits(true).
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		log.Fatalf("error get call_log - %v", err)
	}
	log.Printf("201 %v", search201.TotalHits())

	// 202
	term202 := elastic.NewTermQuery("error_code", 202)
	search202, err := elasticClient.Search().
		Index(index).                                                                                          // search in index "tweets"
		Query(elastic.NewBoolQuery().Must(timeRange).Must(term202).MustNot(elastic.NewExistsQuery("jitter"))). // specify the query
		Sort("@timestamp", true).                                                                              // sort by "user" field, ascending
		From(0).Size(10000).                                                                                   // take documents 0-9
		TrackTotalHits(true).
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		log.Fatalf("error get call_log - %v", err)
	}
	log.Printf("202 %v", search202.TotalHits())

	// 400
	term400 := elastic.NewTermQuery("error_code", 400)
	search400, err := elasticClient.Search().
		Index(index).                                                                                          // search in index "tweets"
		Query(elastic.NewBoolQuery().Must(timeRange).Must(term400).MustNot(elastic.NewExistsQuery("jitter"))). // specify the query
		Sort("@timestamp", true).                                                                              // sort by "user" field, ascending
		From(0).Size(10000).                                                                                   // take documents 0-9
		TrackTotalHits(true).
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		log.Fatalf("error get call_log - %v", err)
	}
	log.Printf("400 %v", search400.TotalHits())

	// 401
	term401 := elastic.NewTermQuery("error_code", 401)
	search401, err := elasticClient.Search().
		Index(index).                                                                                          // search in index "tweets"
		Query(elastic.NewBoolQuery().Must(timeRange).Must(term401).MustNot(elastic.NewExistsQuery("jitter"))). // specify the query
		Sort("@timestamp", true).                                                                              // sort by "user" field, ascending
		From(0).Size(10000).                                                                                   // take documents 0-9
		TrackTotalHits(true).
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		log.Fatalf("error get call_log - %v", err)
	}
	log.Printf("401 %v", search401.TotalHits())

	// 403
	term403 := elastic.NewTermQuery("error_code", 403)
	search403, err := elasticClient.Search().
		Index(index).                                                                                          // search in index "tweets"
		Query(elastic.NewBoolQuery().Must(timeRange).Must(term403).MustNot(elastic.NewExistsQuery("jitter"))). // specify the query
		Sort("@timestamp", true).                                                                              // sort by "user" field, ascending
		From(0).Size(10000).                                                                                   // take documents 0-9
		TrackTotalHits(true).
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		log.Fatalf("error get call_log - %v", err)
	}
	log.Printf("403 %v", search403.TotalHits())

	// 409
	term409 := elastic.NewTermQuery("error_code", 409)
	search409, err := elasticClient.Search().
		Index(index).                                                                                          // search in index "tweets"
		Query(elastic.NewBoolQuery().Must(timeRange).Must(term409).MustNot(elastic.NewExistsQuery("jitter"))). // specify the query
		Sort("@timestamp", true).                                                                              // sort by "user" field, ascending
		From(0).Size(10000).                                                                                   // take documents 0-9
		TrackTotalHits(true).
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		log.Fatalf("error get call_log - %v", err)
	}
	log.Printf("409 %v", search409.TotalHits())

	// 422
	term422 := elastic.NewTermQuery("error_code", 422)
	search422, err := elasticClient.Search().
		Index(index).                                                                                          // search in index "tweets"
		Query(elastic.NewBoolQuery().Must(timeRange).Must(term422).MustNot(elastic.NewExistsQuery("jitter"))). // specify the query
		Sort("@timestamp", true).                                                                              // sort by "user" field, ascending
		From(0).Size(10000).                                                                                   // take documents 0-9
		TrackTotalHits(true).
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		log.Fatalf("error get call_log - %v", err)
	}
	log.Printf("422 %v", search422.TotalHits())

	// 429
	term429 := elastic.NewTermQuery("error_code", 429)
	search429, err := elasticClient.Search().
		Index(index).                                                                                          // search in index "tweets"
		Query(elastic.NewBoolQuery().Must(timeRange).Must(term429).MustNot(elastic.NewExistsQuery("jitter"))). // specify the query
		Sort("@timestamp", true).                                                                              // sort by "user" field, ascending
		From(0).Size(10000).                                                                                   // take documents 0-9
		TrackTotalHits(true).
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		log.Fatalf("error get call_log - %v", err)
	}
	log.Printf("429 %v", search429.TotalHits())

	// 500
	term500 := elastic.NewTermQuery("error_code", 500)
	search500, err := elasticClient.Search().
		Index(index).                                                                                          // search in index "tweets"
		Query(elastic.NewBoolQuery().Must(timeRange).Must(term500).MustNot(elastic.NewExistsQuery("jitter"))). // specify the query
		Sort("@timestamp", true).                                                                              // sort by "user" field, ascending
		From(0).Size(10000).                                                                                   // take documents 0-9
		TrackTotalHits(true).
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		log.Fatalf("error get call_log - %v", err)
	}
	log.Printf("500 %v", search500.TotalHits())

	// 502
	term502 := elastic.NewTermQuery("error_code", 502)
	search502, err := elasticClient.Search().
		Index(index).                                                                                          // search in index "tweets"
		Query(elastic.NewBoolQuery().Must(timeRange).Must(term502).MustNot(elastic.NewExistsQuery("jitter"))). // specify the query
		Sort("@timestamp", true).                                                                              // sort by "user" field, ascending
		From(0).Size(10000).                                                                                   // take documents 0-9
		TrackTotalHits(true).
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		log.Fatalf("error get call_log - %v", err)
	}
	log.Printf("502 %v", search502.TotalHits())

	// 503
	term503 := elastic.NewTermQuery("error_code", 503)
	search503, err := elasticClient.Search().
		Index(index).                                                                                          // search in index "tweets"
		Query(elastic.NewBoolQuery().Must(timeRange).Must(term503).MustNot(elastic.NewExistsQuery("jitter"))). // specify the query
		Sort("@timestamp", true).                                                                              // sort by "user" field, ascending
		From(0).Size(10000).                                                                                   // take documents 0-9
		TrackTotalHits(true).
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		log.Fatalf("error get call_log - %v", err)
	}
	log.Printf("503 %v", search503.TotalHits())

	// 504
	term504 := elastic.NewTermQuery("error_code", 504)
	search504, err := elasticClient.Search().
		Index(index).                                                                                          // search in index "tweets"
		Query(elastic.NewBoolQuery().Must(timeRange).Must(term504).MustNot(elastic.NewExistsQuery("jitter"))). // specify the query
		Sort("@timestamp", true).                                                                              // sort by "user" field, ascending
		From(0).Size(10000).                                                                                   // take documents 0-9
		TrackTotalHits(true).
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		log.Fatalf("error get call_log - %v", err)
	}
	log.Printf("504 %v", search504.TotalHits())

	// 505
	term505 := elastic.NewTermQuery("error_code", 505)
	search505, err := elasticClient.Search().
		Index(index).                                                                                          // search in index "tweets"
		Query(elastic.NewBoolQuery().Must(timeRange).Must(term505).MustNot(elastic.NewExistsQuery("jitter"))). // specify the query
		Sort("@timestamp", true).                                                                              // sort by "user" field, ascending
		From(0).Size(10000).                                                                                   // take documents 0-9
		TrackTotalHits(true).
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		log.Fatalf("error get call_log - %v", err)
	}
	log.Printf("505 %v", search505.TotalHits())

	// 507
	term507 := elastic.NewTermQuery("error_code", 507)
	search507, err := elasticClient.Search().
		Index(index).                                                                                          // search in index "tweets"
		Query(elastic.NewBoolQuery().Must(timeRange).Must(term507).MustNot(elastic.NewExistsQuery("jitter"))). // specify the query
		Sort("@timestamp", true).                                                                              // sort by "user" field, ascending
		From(0).Size(10000).                                                                                   // take documents 0-9
		TrackTotalHits(true).
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		log.Fatalf("error get call_log - %v", err)
	}
	log.Printf("507 %v", search507.TotalHits())

	// 508
	term508 := elastic.NewTermQuery("error_code", 508)
	search508, err := elasticClient.Search().
		Index(index).                                                                                          // search in index "tweets"
		Query(elastic.NewBoolQuery().Must(timeRange).Must(term508).MustNot(elastic.NewExistsQuery("jitter"))). // specify the query
		Sort("@timestamp", true).                                                                              // sort by "user" field, ascending
		From(0).Size(10000).                                                                                   // take documents 0-9
		TrackTotalHits(true).
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		log.Fatalf("error get call_log - %v", err)
	}
	log.Printf("508 %v", search508.TotalHits())

	// 509
	term509 := elastic.NewTermQuery("error_code", 509)
	search509, err := elasticClient.Search().
		Index(index).                                                                                          // search in index "tweets"
		Query(elastic.NewBoolQuery().Must(timeRange).Must(term509).MustNot(elastic.NewExistsQuery("jitter"))). // specify the query
		Sort("@timestamp", true).                                                                              // sort by "user" field, ascending
		From(0).Size(10000).                                                                                   // take documents 0-9
		TrackTotalHits(true).
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		log.Fatalf("error get call_log - %v", err)
	}
	log.Printf("509 %v", search509.TotalHits())

	// 510
	term510 := elastic.NewTermQuery("error_code", 510)
	search510, err := elasticClient.Search().
		Index(index).                                                                                          // search in index "tweets"
		Query(elastic.NewBoolQuery().Must(timeRange).Must(term510).MustNot(elastic.NewExistsQuery("jitter"))). // specify the query
		Sort("@timestamp", true).                                                                              // sort by "user" field, ascending
		From(0).Size(10000).                                                                                   // take documents 0-9
		TrackTotalHits(true).
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		log.Fatalf("error get call_log - %v", err)
	}
	log.Printf("510 %v", search510.TotalHits())

	/*
		MNP
	*/
	// Vinaphone
	queryMnpVnp := elastic.NewQueryStringQuery("Didong*").DefaultField("route")
	searchMnpVnp, err := elasticClient.Search().
		Index(index).                                                                                                            // search in index "tweets"
		Query(elastic.NewBoolQuery().Must(timeRange).Must(term400).Must(queryMnpVnp).MustNot(elastic.NewExistsQuery("jitter"))). // specify the query
		Sort("@timestamp", true).                                                                                                // sort by "user" field, ascending
		From(0).Size(10000).                                                                                                     // take documents 0-9
		TrackTotalHits(true).
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		log.Fatalf("error get call_log - %v", err)
	}
	log.Printf("MNP Vnp %v", searchMnpVnp.TotalHits())

	// Mobifone
	queryMnpMbf := elastic.NewQueryStringQuery("mbf*").DefaultField("route")
	searchMnpMbf, err := elasticClient.Search().
		Index(index).                                                                                                            // search in index "tweets"
		Query(elastic.NewBoolQuery().Must(timeRange).Must(term400).Must(queryMnpMbf).MustNot(elastic.NewExistsQuery("jitter"))). // specify the query
		Sort("@timestamp", true).                                                                                                // sort by "user" field, ascending
		From(0).Size(10000).                                                                                                     // take documents 0-9
		TrackTotalHits(true).
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		log.Fatalf("error get call_log - %v", err)
	}
	log.Printf("MNP Mbf %v", searchMnpMbf.TotalHits())

	// Viettel
	queryMnpVtl := elastic.NewQueryStringQuery("viettel*").DefaultField("route")
	searchMnpVtl, err := elasticClient.Search().
		Index(index).                                                                                                            // search in index "tweets"
		Query(elastic.NewBoolQuery().Must(timeRange).Must(term400).Must(queryMnpVtl).MustNot(elastic.NewExistsQuery("jitter"))). // specify the query
		Sort("@timestamp", true).                                                                                                // sort by "user" field, ascending
		From(0).Size(10000).                                                                                                     // take documents 0-9
		TrackTotalHits(true).
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		log.Fatalf("error get call_log - %v", err)
	}
	log.Printf("MNP Vtl %v", searchMnpVtl.TotalHits())

	// Gtel
	queryMnpGtel := elastic.NewQueryStringQuery("gtel*").DefaultField("route")
	searchMnpGtel, err := elasticClient.Search().
		Index(index).                                                                                                             // search in index "tweets"
		Query(elastic.NewBoolQuery().Must(timeRange).Must(term400).Must(queryMnpGtel).MustNot(elastic.NewExistsQuery("jitter"))). // specify the query
		Sort("@timestamp", true).                                                                                                 // sort by "user" field, ascending
		From(0).Size(10000).                                                                                                      // take documents 0-9
		TrackTotalHits(true).
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		log.Fatalf("error get call_log - %v", err)
	}
	log.Printf("MNP Gtel %v", searchMnpGtel.TotalHits())

	// Itel
	queryMnpItel := elastic.NewQueryStringQuery("itel*").DefaultField("route")
	searchMnpItel, err := elasticClient.Search().
		Index(index).                                                                                                             // search in index "tweets"
		Query(elastic.NewBoolQuery().Must(timeRange).Must(term400).Must(queryMnpItel).MustNot(elastic.NewExistsQuery("jitter"))). // specify the query
		Sort("@timestamp", true).                                                                                                 // sort by "user" field, ascending
		From(0).Size(10000).                                                                                                      // take documents 0-9
		TrackTotalHits(true).
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		log.Fatalf("error get call_log - %v", err)
	}
	log.Printf("MNP Itel %v", searchMnpItel.TotalHits())

	// Mobicast
	queryMnpMbc := elastic.NewQueryStringQuery("mobicast*").DefaultField("route")
	searchMnpMbc, err := elasticClient.Search().
		Index(index).                                                                                                            // search in index "tweets"
		Query(elastic.NewBoolQuery().Must(timeRange).Must(term400).Must(queryMnpMbc).MustNot(elastic.NewExistsQuery("jitter"))). // specify the query
		Sort("@timestamp", true).                                                                                                // sort by "user" field, ascending
		From(0).Size(10000).                                                                                                     // take documents 0-9
		TrackTotalHits(true).
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		log.Fatalf("error get call_log - %v", err)
	}
	log.Printf("MNP Mbc %v", searchMnpMbc.TotalHits())

	records := [][]string{
		{"No", "Description", "Amount", "Percent", "Note", "Action"},
		{"1", "Total Call", fmt.Sprint(searchTotal.TotalHits()), "100%"},
		{"2", "Normal Call", fmt.Sprint(search100.TotalHits()), fmt.Sprintf("%.2f", math.Round(float64(search100.TotalHits())/float64(searchTotal.TotalHits())*float64(100))) + "%", "EC 100", "Bypass"},
		// {"3", "OK", fmt.Sprint(search200.TotalHits()), fmt.Sprintf("%.2f", math.Round(float64(search200.TotalHits())/float64(searchTotal.TotalHits())*float64(100))) + "%", "EC 200", "Bypass"},
		{"3", "VLR is Whitelisted", fmt.Sprint(search201.TotalHits()), fmt.Sprintf("%.2f", math.Round(float64(search201.TotalHits())/float64(searchTotal.TotalHits())*float64(100))) + "%", "EC 201", "Bypass"},
		{"4", "Whitelisted B Number", fmt.Sprint(search202.TotalHits()), fmt.Sprintf("%.2f", math.Round(float64(search202.TotalHits())/float64(searchTotal.TotalHits())*float64(100))) + "%", "EC 202", "Bypass"},
		{"5", "Not Roaming", fmt.Sprint(search400.TotalHits()), fmt.Sprintf("%.2f", math.Round(float64(search400.TotalHits())/float64(searchTotal.TotalHits())*float64(100))) + "%", "EC 400", "Block"},
		{"6", "Blacklisted", fmt.Sprint(search401.TotalHits()), fmt.Sprintf("%.2f", math.Round(float64(search401.TotalHits())/float64(searchTotal.TotalHits())*float64(100))) + "%", "EC 401", "Block"},
		{"7", "Forbidden VLR and IMSI", fmt.Sprint(search403.TotalHits()), fmt.Sprintf("%.2f", math.Round(float64(search403.TotalHits())/float64(searchTotal.TotalHits())*float64(100))) + "%", "EC 403", "Block"},
		{"8", "Unknown Subscriber", fmt.Sprint(search409.TotalHits()), fmt.Sprintf("%.2f", math.Round(float64(search409.TotalHits())/float64(searchTotal.TotalHits())*float64(100))) + "%", "EC 409", "Block"},
		{"9", "Route Not Found", fmt.Sprint(search422.TotalHits()), fmt.Sprintf("%.2f", math.Round(float64(search422.TotalHits())/float64(searchTotal.TotalHits())*float64(100))) + "%", "EC 422", "Bypass"},
		{"10", "Calling is Temporary Blacklisted", fmt.Sprint(search429.TotalHits()), fmt.Sprintf("%.2f", math.Round(float64(search429.TotalHits())/float64(searchTotal.TotalHits())*float64(100))) + "%", "EC 429", "Block"},
		{"11", "Internal Error", fmt.Sprint(search500.TotalHits()), fmt.Sprintf("%.2f", math.Round(float64(search500.TotalHits())/float64(searchTotal.TotalHits())*float64(100))) + "%", "EC 500", "Block"},
		{"12", "Invalid Subscriber Information", fmt.Sprint(search502.TotalHits()), fmt.Sprintf("%.2f", math.Round(float64(search502.TotalHits())/float64(searchTotal.TotalHits())*float64(100))) + "%", "EC 502", "Block"},
		{"13", "Forbidden HLR", fmt.Sprint(search503.TotalHits()), fmt.Sprintf("%.2f", math.Round(float64(search503.TotalHits())/float64(searchTotal.TotalHits())*float64(100))) + "%", "EC 503", "Bypass"},
		{"14", "HLR Late Response", fmt.Sprint(search504.TotalHits()), fmt.Sprintf("%.2f", math.Round(float64(search504.TotalHits())/float64(searchTotal.TotalHits())*float64(100))) + "%", "EC 504", "Bypass"},
		{"15", "Teleservice Not Provisioned", fmt.Sprint(search505.TotalHits()), fmt.Sprintf("%.2f", math.Round(float64(search505.TotalHits())/float64(searchTotal.TotalHits())*float64(100))) + "%", "EC 505", "Block"},
		{"16", "Call Barred", fmt.Sprint(search507.TotalHits()), fmt.Sprintf("%.2f", math.Round(float64(search507.TotalHits())/float64(searchTotal.TotalHits())*float64(100))) + "%", "EC 507", "Block"},
		{"17", "Absent Sri-Sm ", fmt.Sprint(search508.TotalHits()), fmt.Sprintf("%.2f", math.Round(float64(search508.TotalHits())/float64(searchTotal.TotalHits())*float64(100))) + "%", "EC 508", "Block"},
		{"18", "Absent Subscriber ", fmt.Sprint(search509.TotalHits()), fmt.Sprintf("%.2f", math.Round(float64(search509.TotalHits())/float64(searchTotal.TotalHits())*float64(100))) + "%", "EC 509", "Block"},
		{"19", "Not Extended", fmt.Sprint(search510.TotalHits()), fmt.Sprintf("%.2f", math.Round(float64(search510.TotalHits())/float64(searchTotal.TotalHits())*float64(100))) + "%", "EC 510", "Block"},
	}

	f, err := os.Create("vnptspoofing_report_09.csv")
	defer f.Close()

	if err != nil {

		log.Fatalln("failed to open file", err)
	}

	w := csv.NewWriter(f)
	err = w.WriteAll(records) // calls Flush internally

	if err != nil {
		log.Fatal(err)
	}
}
