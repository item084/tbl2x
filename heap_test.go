package tbl2x

import (
	//"container/heap"
	"testing"
	//"github.com/gonum/matrix/mat64"
)

func TestHeap(t *testing.T) {
	//mat := mat64.NewDense(3, 3, []float64{0.0, 1.0, 2.0, 1.0, 0.0, 3.0, 2.0, 3.0, 0.0})
	tbl := Table{[]string{"1", "2", "3"}, []string{}, []string{}, []string{"a", "b", "c"}, 3, 3, []float64{0.0, 1.0, 2.0, 1.0, 10.0, 3.0, 2.0, 3.0, 10.0}, "noname", "test"}
	a, _ := TblTopK(&tbl, 1, 1, 3)

	t.Log(tbl.PrettyStringChosenRows(a, 2))

}
