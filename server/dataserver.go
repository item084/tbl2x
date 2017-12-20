package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/gorilla/mux"
	"github.com/nimezhu/data"
)

func main() {
	uri := os.Args[1]
	router := mux.NewRouter()
	l := data.NewLoader("./tbl2x")
	l.Plugins["tsv"] = func(dbname string, data interface{}) (data.DataRouter, error) {
		switch v := data.(type) {
		default:
			fmt.Printf("unexpected type %T", v)
			return nil, errors.New(fmt.Sprintf("bigwig format not support type %T", v))
		case string:
			return nil, errors.New("todo")
		case map[string]interface{}:
			for k0, v0 := range data.(map[string]interface{}) {
				fmt.Println(k0, v0.(string))
			}
			return nil, errors.New("TODO")
		}
	}
	l.Load(uri, router)
}
