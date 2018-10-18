package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

type data struct {
	//Date    int
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
		//ret += fmt.Sprintf("Date: %v, Value %v, Quality: %v\n", v.Value[key].Date, v.Value[key].Value, v.Value[key].Quality)
		ret += fmt.Sprintf("Value %v, Quality: %v\n", v.Value[key].Value, v.Value[key].Quality)
	}
	ret += fmt.Sprintf("Key: %v, Name: %v", v.Station.Key, v.Station.Name)
	return ret
}

type safevalues struct {
	v values
}

const dataFolderPath = "C:\\tmp\\smhi"
const loopTime = time.Second * 15

func (sv *safevalues) GetWaterFlow() {
	for {
		url := "https://opendata-download-hydroobs.smhi.se/api/version/latest/parameter/1/station/855/period/latest-hour/data.json"
		resp, err := http.Get(url)
		if err != nil {
			log.Print(err)
			sv.errorFileHandler(1)
			time.Sleep(loopTime)
			continue
		}
		defer resp.Body.Close()

		var b []byte
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
			sv.errorFileHandler(1)
			time.Sleep(loopTime)
			continue
		}

		err = json.Unmarshal(b, &sv.v)
		if err != nil {
			log.Print(err)
			sv.errorFileHandler(1)
			time.Sleep(loopTime)
			continue
		}
		sv.fileHandler()
		sv.errorFileHandler(0)
		time.Sleep(loopTime)
	}
}
func main() {
	os.Mkdir(dataFolderPath, 0777)
	f, err := os.OpenFile(path.Join(dataFolderPath, "log"), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Error opening log file", err)
	}

	defer f.Close()

	log.SetOutput(f)
	var s safevalues
	s.GetWaterFlow()

}

func (sv *safevalues) fileHandler() {
	var filenameFlow string = "flode1.txt"

	var b []byte = []byte(fmt.Sprintf("%.3f", sv.v.Value[0].Value))

	err := ioutil.WriteFile(path.Join(dataFolderPath, filenameFlow), b, 0777)
	if err != nil {
		log.Print("Error writing to file", err)
		sv.errorFileHandler(1)
	}

}

func (sv *safevalues) errorFileHandler(e int) {
	var filenameError string = "error.txt"
	var s string = ""

	if e == 1 {
		s = "true"
	}
	if e == 0 {
		s = "false"
	}
	var b []byte = []byte(string(s))

	err := ioutil.WriteFile(path.Join(dataFolderPath, filenameError), b, 0777)
	if err != nil {
		log.Print("Error writing to errorfile", err)
	}

}
