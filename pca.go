package tbl2x

import (
	"github.com/gonum/matrix/mat64"
	"github.com/gonum/stat"
)

type tablePCA struct {
	vars []float64
	vecs *mat64.Dense
}

func (t *Table) PCA() *tablePCA {
	a := t.Dense()
	weight := make([]float64, t.RowSize)
	for i := 0; i < t.RowSize; i++ {
		weight[i] = 1.0
	}
	var pc stat.PC
	ok := pc.PrincipalComponents(a, weight)
	var vecs *mat64.Dense
	var vars []float64
	if ok {
		vecs = pc.Vectors(vecs)
		vars = pc.Vars(vars)
		return &tablePCA{vars, vecs}
	} else {
		return nil
	}
	return nil
}
