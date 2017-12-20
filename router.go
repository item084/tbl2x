package tbl2x

import "github.com/gorilla/mux"

type TableRouter struct {
	id   string
	data map[string]*Table
}

func (t *TableRouter) ServeTo(r *mux.Router) {

}
