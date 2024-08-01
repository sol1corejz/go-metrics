package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func (ms *MemStorage) SetGauge(name string, value float64) {
	ms.gauge[name] = value
}

func (ms *MemStorage) AddCounter(name string, value int64) {
	ms.counter[name] += value
}

func mainPage(storage *MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Сервис сбора метрик и алертинга\n")

		fmt.Fprint(w, "\ngauge:\n", storage.gauge)
		fmt.Fprint(w, "\n\ncounter:\n", storage.counter)
	}
}

func updatePage(storage *MemStorage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			res.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
			return
		}

		path := strings.TrimPrefix(req.URL.Path, "/update/")
		parts := strings.Split(path, "/")

		metricType := parts[0]
		metricName := parts[1]
		metricValue := parts[2]

		switch metricType {
		case "gauge":
			if val, err := strconv.ParseFloat(metricValue, 64); err == nil {
				storage.SetGauge(metricName, val)
			} else {
				res.Write([]byte(http.StatusText(http.StatusBadRequest)))
			}
			break
		case "counter":
			if val, err := strconv.ParseInt(metricValue, 0, 64); err == nil {
				storage.AddCounter(metricName, val)
			} else {
				res.Write([]byte(http.StatusText(http.StatusBadRequest)))
			}
			break
		default:
			res.Write([]byte(http.StatusText(http.StatusBadRequest)))
			return
		}

		res.Write([]byte(http.StatusText(http.StatusOK)))
	}
}

func main() {

	ms := &MemStorage{gauge: make(map[string]float64), counter: make(map[string]int64)}

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", updatePage(ms))
	mux.HandleFunc("/", mainPage(ms))

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
