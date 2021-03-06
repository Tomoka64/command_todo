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
	"time"

	"github.com/fatih/color"
)

const (
	filename = "config/data.json"
	layout   = "2006-01-02"
)

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
		t, _ := time.Parse(layout, Datas[i].DeadLine)
		if !day.After(t) {
			ret = append(ret, Datas[i])
			fmt.Println(len(ret))
			switch len(ret) {
			case 1:
				bs := ToJson(Datas[i])
				first(bs)
				// _ = ioutil.WriteFile(filename, bs, 0644)
			default:
				for _, data := range ret {
					b := ToJson(data)
					SaveToFile(b)
				}
			}
		} else if day.After(t) {
			color.Red("deleted: %v | %v", Datas[i].Title, Datas[i].DeadLine)
		}
	}
}

func ToJson(s interface{}) []byte {
	data, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func first(bs []byte) {
	f, err := os.OpenFile(filename, os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.Write(bs)
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
