package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const filename = "config/data.json"

var T = flag.String("t", "clean up my room", "put your todo-list")
var D = flag.String("d", "3000-00-00", "set up a deadline for your todo (format-3000-00-00)")

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "<usage> %s <flag option> -t -d\n<usage2> %s update\n<usage2> %s history\n", os.Args[0], os.Args[0], os.Args[0])
		os.Exit(1)
	}
	switch os.Args[1] {
	case "history":
		_ = History()
	case "update":
		Update()
	default:
		HandleDef()
	}
}

func HandleDef() {
	flag.Parse()
	id := 0
	check, _ := getAllFromFile()
	if len(check) > 0 {
		tem := getTheLastId()
		id = tem.Isbn + 1
	}
	timeC := time.Now()
	data := &Todo{
		Isbn:        id,
		Title:       *T,
		TimeCreated: timeC,
		DeadLine:    *D,
	}

	bs := ToJson(data)
	SaveToFile(bs)
}

func Update() {
	Datas, err := getAllFromFile()
	if err != nil {
		log.Fatalln(err)
	}
	var ret Todos
	for i := 0; i < len(Datas); i++ {

		day := time.Now()
		const layout = "2006-01-02"
		time := day.Format(layout)
		intDate := strings.Split(Datas[i].DeadLine, "-")
		timestr := strings.Split(time, "-")
		fmt.Println(intDate, timestr)
		if !Compare(intDate[0], timestr[0]) {
			if !Compare(intDate[1], timestr[1]) {
				if !Compare(intDate[2], timestr[2]) {
					ret = append(ret, Datas[i])
					fmt.Println(ret)
					for _, data := range ret {
						bs := ToJson(data)
						_ = ioutil.WriteFile(filename, bs, 0644)
					}
				}
			}
		}

	}
}

func updateFile(b []byte) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	f.Write(b)
}

func Compare(s string, x string) bool {
	if intfy(s) > intfy(x) {
		return false
	}
	return true
}

func intfy(s string) int {
	a, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println(err)
	}
	return a
}

func ToJson(s interface{}) []byte {
	data, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func SaveToFile(bs []byte) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.Write(bs)
}

func getTheLastId() Todo {
	datas, _ := getAllFromFile()
	return datas[len(datas)-1]
}

func History() error {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	dec := json.NewDecoder(bytes.NewReader(contents))

	for {
		var data Todo
		if err = dec.Decode(&data); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		// datas = append(datas, data)
		fmt.Println(data.Isbn, data.Title, "|", data.DeadLine)
	}
	return err
}

func getAllFromFile() (Todos, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return []Todo{}, err
	}
	var datas Todos
	dec := json.NewDecoder(bytes.NewReader(contents))

	for {
		var data Todo
		if err = dec.Decode(&data); err == io.EOF {
			break
		} else if err != nil {
			return []Todo{}, err
		}
		datas = append(datas, data)
	}
	return datas, nil
}
