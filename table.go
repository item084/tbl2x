package tbl2x

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"strconv"

	"gonum.org/v1/gonum/mat"

	"github.com/nimezhu/netio"
)

type Table struct {
	ColNames      []string
	ColLabelNames []string
	ColLabelData  []string
	RowNames      []string
	ColSize       int
	RowSize       int
	Mat           []float64
	FileName      string
	Name          string
}

func NewTable(r int, c int, data []float64) *Table {
	colNames := make([]string, r)
	rowNames := make([]string, c)
	colLabelName := make([]string, 0)
	colLabels := make([]string, 0)
	for i := range colNames {
		colNames[i] = "C" + strconv.Itoa(i)
	}
	for i := range rowNames {
		rowNames[i] = "R" + strconv.Itoa(i)
	}
	if data == nil {
		data = make([]float64, r*c)
		for i := range data {
			data[i] = 0.0
		}
	}
	fileName := "noname.tsv"
	name := "table"
	return &Table{colNames, colLabelName, colLabels, rowNames, r, c, data, fileName, name}

}
func (t *Table) Dims() (int, int) {
	return t.RowSize, t.ColSize
}
func (t *Table) Cols() []string {
	return t.ColNames
}
func (t *Table) Rows() []string {
	return t.RowNames
}
func (t *Table) Dense() *mat.Dense {
	return mat.NewDense(t.RowSize, t.ColSize, t.Mat)
}
func (t *Table) ColLabelNum() int {
	return len(t.ColLabelData) / t.RowSize
}

/*
func (t *Table) ColLabel(i int) []string { //TODO

}
*/
func (t *Table) String() string {
	return t.PrettyString(-1)
}
func (t *Table) TxtEncode() string {
	return t.PrettyString(2)
}

/* TODO handle colLabels?? */
func (t *Table) T() error {
	m := t.Dense().RawMatrix().Data
	t.ColNames, t.RowNames = t.RowNames, t.ColNames
	t.ColSize, t.RowSize = t.RowSize, t.ColSize
	t.FileName = t.FileName + "_transpose.tsv"
	t.Name = t.Name + "_transpose"
	data := make([]float64, t.RowSize*t.ColSize)
	for i := 0; i < t.ColSize; i++ {
		for j := 0; j < t.RowSize; j++ {
			data[j*t.ColSize+i] = m[i*t.RowSize+j]
		}
	}
	t.Mat = data
	return nil
}
func (t *Table) Info() string {
	r, c := t.Dims()
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Dims: %d * %d\n\nRowNames:\n\t", r, c))
	k := r
	s := ""
	if k > 20 {
		k = 20
		s = "..."
	}
	for i := 0; i < k; i++ {
		buffer.WriteString(t.RowNames[i] + ",")
	}
	buffer.WriteString(s)
	buffer.WriteString("\n\nColNames:\n\t")
	s = ""
	k = c
	if k > 20 {
		k = 20
		s = "..."
	}
	for j := 0; j < k; j++ {
		buffer.WriteString(t.ColNames[j] + ",")
	}
	buffer.WriteString(s)
	buffer.WriteString("\n\n")
	matv := t.Dense()
	max := mat.Max(matv)
	min := mat.Min(matv)
	buffer.WriteString(fmt.Sprintf("Domain: [%f , %f]", min, max))
	buffer.WriteString("\n\n")
	return buffer.String()
}

func (t *Table) Log(e float64, pseudo float64) error {
	a := make([]float64, len(t.Mat))
	root := math.Log(e)
	for i, v := range t.Mat {
		a[i] = math.Log(v+pseudo) / root
	}
	t.Mat = a
	t.Name += "|log"
	return nil
}

func (table *Table) loadReader(f io.Reader, fn string, n int) error {
	r := csv.NewReader(f)
	r.Comma = '\t'
	table.FileName = fn
	iter, err := r.ReadAll()
	if err != nil {
		return err
	}
	table.Name = iter[0][0]
	table.ColNames = iter[0][(n + 1):]
	table.ColLabelNames = iter[0][1:(n + 1)]
	table.ColSize = len(table.ColNames)
	table.RowSize = len(iter) - 1
	table.ColLabelData = make([]string, len(table.ColLabelNames)*table.RowSize)
	table.RowNames = make([]string, table.RowSize)
	table.Mat = make([]float64, table.ColSize*table.RowSize)
	for i := 1; i < len(iter); i++ {
		name, values := iter[i][0], iter[i][(n+1):]
		for j := 0; j < len(values); j++ {
			if values[j] == "NA" {
				table.Mat[(i-1)*table.ColSize+j] = math.NaN()
			} else {
				table.Mat[(i-1)*table.ColSize+j], err = strconv.ParseFloat(values[j], 64)
				if err != nil {
					return err
				}
			}
		}
		colLabelValues := iter[i][1:(n + 1)]
		for j := 0; j < len(colLabelValues); j++ {
			table.ColLabelData[(i-1)*table.ColLabelNum()+j] = colLabelValues[j]
		}
		table.RowNames[i-1] = name
	}
	return err
}
func (table *Table) LoadTsv(file string, n int) error {
	f, err := netio.Open(file)
	defer f.Close()
	if err != nil {
		return err
	}
	err = table.loadReader(f, file, n)
	//table := new(Table)
	return err
}
func Load(fn string, n int) (*Table, error) {
	t := &Table{}
	err := t.LoadTsv(fn, n)
	return t, err
}
func (t *Table) PrettyStringSubmat(rows []int, cols []int, f int) string {
	var buffer bytes.Buffer
	buffer.WriteString(t.Name + "_sub")
	for i0 := 0; i0 < len(cols); i0++ {
		s := fmt.Sprintf("\t%s", t.ColNames[cols[i0]])
		buffer.WriteString(s)
	}
	buffer.WriteString("\n")
	format := "\t" + "%." + strconv.Itoa(f) + "f"
	if f == -1 {
		format = "\t%f"
	}
	m := t.Dense()
	for i := 0; i < len(rows); i++ {
		buffer.WriteString(fmt.Sprintf("%s", t.RowNames[rows[i]]))
		for j := 0; j < len(cols); j++ {
			buffer.WriteString(fmt.Sprintf(format, m.At(rows[i], cols[j])))
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}
func (t *Table) PrettyStringChosenRows(rows []int, f int) string {
	_, c := t.Dims()
	cols := make([]int, c)
	for i := 0; i < c; i++ {
		cols[i] = i
	}
	return t.PrettyStringSubmat(rows, cols, f)
}
func (t *Table) PrettyStringChosenCols(cols []int, f int) string {
	r, _ := t.Dims()
	rows := make([]int, r)
	for i := 0; i < r; i++ {
		rows[i] = i
	}
	return t.PrettyStringSubmat(rows, cols, f)
}

func (t *Table) PrettyString(f int) string {
	r, c := t.Dims()
	var buffer bytes.Buffer
	buffer.WriteString(t.Name)

	for i0 := 0; i0 < c; i0++ {
		s := fmt.Sprintf("\t%s", t.ColNames[i0])
		buffer.WriteString(s)
	}

	buffer.WriteString("\n")
	format := "\t" + "%." + strconv.Itoa(f) + "f"
	if f == -1 {
		format = "\t%f"
	}
	m := t.Dense()
	for i := 0; i < r; i++ {
		buffer.WriteString(fmt.Sprintf("%s", t.RowNames[i]))
		for j := 0; j < c; j++ {
			buffer.WriteString(fmt.Sprintf(format, m.At(i, j)))
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}
