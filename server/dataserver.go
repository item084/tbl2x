package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/nimezhu/data"
	"github.com/nimezhu/tbl2x"
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
			r := &tbl2x.TableRouter{dbname, make(map[string]*tbl2x.Table)}
			err := r.Load(data.(map[string]interface{}))
			return r, err
		}
	}
	l.Load(uri, router)
	log.Println("Listening...")
	log.Println("Please open http://127.0.0.1:8080")
	http.ListenAndServe(":8080", router)
}
