package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type data struct {
	Date    int
	Value   float32
	Quality string
}
type station struct {
	Key  string
	Name string
}
type values struct {
	Value   []data
	Station station
}

func (v values) String() string {
	var ret string
	for key, _ := range v.Value {
		ret += fmt.Sprintf("Date: %v, Value %v, Quality: %v\n", v.Value[key].Date, v.Value[key].Value, v.Value[key].Quality)
	}
	ret += fmt.Sprintf("Key: %v, Name: %v", v.Station.Key, v.Station.Name)
	return ret
}

type safevalues struct {
	v   values
	mux sync.Mutex
}

var s safevalues

func (sv *safevalues) GetWaterFlow() {
	for {
		url := "https://opendata-download-hydroobs.smhi.se/api/version/latest/parameter/1/station/855/period/latest-hour/data.json"
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		var b []byte
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		sv.mux.Lock()
		err = json.Unmarshal(b, &sv.v)
		if err != nil {
			log.Fatal(err)
		}
		sv.mux.Unlock()
		time.Sleep(time.Hour * 1)
	}
}
func main() {
	go s.GetWaterFlow()
	http.HandleFunc("/", roothandler)
	http.HandleFunc("/data", datahandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func roothandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, nil)
}

func datahandler(w http.ResponseWriter, r *http.Request) {
	s.mux.Lock()
	json.NewEncoder(w).Encode(s.v)
	s.mux.Unlock()
}
