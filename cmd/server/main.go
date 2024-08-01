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

func updatePage(storage *MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/update/")
		parts := strings.Split(path, "/")

		if len(parts) != 3 {
			http.Error(w, "invalid request", http.StatusNotFound)
			return
		}

		metricType, metricName, metricValue := parts[0], parts[1], parts[2]

		switch metricType {
		case "gauge":
			val, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				http.Error(w, "invalid gauge", http.StatusBadRequest)
				return
			}
			storage.SetGauge(metricName, val)

		case "counter":
			val, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				http.Error(w, "invalid counter", http.StatusBadRequest)
			}
			storage.AddCounter(metricName, val)
			return
		default:
			http.Error(w, "invalid metric type", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	}
}

func main() {

	ms := &MemStorage{gauge: make(map[string]float64), counter: make(map[string]int64)}

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", updatePage(ms))

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
