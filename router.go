package tbl2x

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type TableRouter struct {
	Id   string
	Data map[string]*Table
}

func smartParseIdxs(s string, key map[string]int) []int {
	//TODO
	return []int{0}
}
func (t *TableRouter) Load(d map[string]interface{}) error {
	for k, v := range d {
		switch v.(type) {
		case []string:
			n, _ := strconv.Atoi(v.([]string)[1])
			tbl, err := Load(v.([]string)[0], n) //TODO 0
			if err == nil {
				t.Data[k] = tbl
			} else {
				return err
			}
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
	uriMap := make(map[string]string)
	rowMap := make(map[string]map[string]int)
	colMap := make(map[string]map[string]int)
	for k, v := range t.Data {
		uriMap[k] = v.FileName
		rowMap[k] = make(map[string]int)
		colMap[k] = make(map[string]int)
		for i0, v0 := range v.ColNames {
			colMap[k][v0] = i0
		}
		for i0, v0 := range v.RowNames {
			rowMap[k][v0] = i0
		}
	}
	r.HandleFunc("/"+t.Id, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("todo test"))
	})
	r.HandleFunc("/"+t.Id+"/list", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		keys := []string{}
		for key, _ := range t.Data {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		j, _ := json.Marshal(keys)
		w.Write(j)
	})
	r.HandleFunc("/"+t.Id+"/ls", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		j, _ := json.Marshal(uriMap)
		w.Write(j)
	})
	r.HandleFunc("/"+t.Id+"/get/{id}/size", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		params := mux.Vars(r)
		id := params["id"]
		a, ok := t.Data[id]
		if ok {
			r, c := a.Dims()
			o := map[string]int{
				"row": r,
				"col": c,
			}
			j, _ := json.Marshal(o)
			w.Write(j)
		} else {
			w.Write([]byte("{error:'not found'}"))
		}
	})
	r.HandleFunc("/"+t.Id+"/get/{id}/print/{res}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		params := mux.Vars(r)
		res := params["res"]
		id := params["id"]
		n, err := strconv.Atoi(res)
		if err != nil {
			n = 0
		}
		a, ok := t.Data[id]
		if ok {
			w.Write([]byte(a.PrettyString(n)))
		} else {
			w.Write([]byte("{error:'not found'}"))
		}
	})
	r.HandleFunc("/"+t.Id+"/get/{id}/colnames", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		params := mux.Vars(r)
		id := params["id"]
		a, ok := t.Data[id]
		if ok {
			j, _ := json.Marshal(a.ColNames)
			w.Write(j)
		} else {
			w.Write([]byte("{error:'not found'}"))
		}
	})
	r.HandleFunc("/"+t.Id+"/get/{id}/rownames", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		params := mux.Vars(r)
		id := params["id"]

		a, ok := t.Data[id]
		if ok {
			j, _ := json.Marshal(a.RowNames)
			w.Write(j)
		} else {
			w.Write([]byte("{error:'not found'}"))
		}
	})
	r.HandleFunc("/"+t.Id+"/get/{id}/rows/{rows}/print/{res}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		params := mux.Vars(r)
		id := params["id"]
		res := params["res"]
		n, err := strconv.Atoi(res)
		if err != nil {
			n = 0
		}
		a, ok := t.Data[id]
		_rows := params["rows"]
		rows := strings.Split(_rows, ",")
		rowidxs := make([]int, len(rows))
		if ok {
			for i, rowid := range rows {
				if v, ok := rowMap[id][rowid]; ok {
					rowidxs[i] = v
				} else {
					rowidxs[i], _ = strconv.Atoi(rowid)
				}
			}
			w.Write([]byte(a.PrettyStringChosenRows(rowidxs, n)))
		} else {
			w.Write([]byte("{error:'not found'}"))
		}
	})
	r.HandleFunc("/"+t.Id+"/get/{id}/cols/{cols}/print/{res}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		params := mux.Vars(r)
		id := params["id"]
		res := params["res"]
		n, err := strconv.Atoi(res)
		if err != nil {
			n = 0
		}
		a, ok := t.Data[id]
		_cols := params["cols"]
		cols := strings.Split(_cols, ",")
		colidxs := make([]int, len(cols))
		if ok {
			for i, colid := range cols {
				if v, ok := colMap[id][colid]; ok {
					colidxs[i] = v
				} else {
					colidxs[i], _ = strconv.Atoi(colid)
				}
			}
			w.Write([]byte(a.PrettyStringChosenCols(colidxs, n)))
		} else {
			w.Write([]byte("{error:'not found'}"))
		}
	})

	r.HandleFunc("/"+t.Id+"/get/{id}/submat/{rows}/{cols}/print/{res}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		params := mux.Vars(r)
		id := params["id"]
		res := params["res"]
		n, err := strconv.Atoi(res)
		if err != nil {
			n = 0
		}
		a, ok := t.Data[id]
		_cols := params["cols"]
		cols := strings.Split(_cols, ",")
		colidxs := make([]int, len(cols))
		_rows := params["rows"]
		rows := strings.Split(_rows, ",")
		rowidxs := make([]int, len(rows))
		if ok {
			for i, colid := range cols {
				if v, ok := colMap[id][colid]; ok {
					colidxs[i] = v
				} else {
					colidxs[i], _ = strconv.Atoi(colid)
				}
			}
			for i, rowid := range rows {
				if v, ok := rowMap[id][rowid]; ok {
					rowidxs[i] = v
				} else {
					rowidxs[i], _ = strconv.Atoi(rowid)
				}
			}
			w.Write([]byte(a.PrettyStringSubmat(rowidxs, colidxs, n)))
		} else {
			w.Write([]byte("{error:'not found'}"))
		}
	})
}
