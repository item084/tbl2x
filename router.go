package tbl2x

import (
	"net/http"

	"github.com/gorilla/mux"
)

type TableRouter struct {
	Id   string
	Data map[string]*Table
}

func (t *TableRouter) Load(d map[string]interface{}) error {
	for k, v := range d {
		switch v.(type) {
		case string:
			tbl, err := Load(v.(string), 0) //TODO 0
			if err == nil {
				t.Data[k] = tbl
			} else {
				return err
			}
		case *Table:
			t.Data[k] = v.(*Table)
		}
	}
	return nil
}
func (t *TableRouter) ServeTo(r *mux.Router) {
	r.HandleFunc("/"+t.Id, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("todo test"))
	})
}
