package tbl2x

import (
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
)

type tablePCA struct {
	vars []float64
	vecs *mat.Dense
}

func (t *Table) PCA() *tablePCA {
	a := t.Dense()
	weight := make([]float64, t.RowSize)
	for i := 0; i < t.RowSize; i++ {
		weight[i] = 1.0
	}
	var pc stat.PC
	ok := pc.PrincipalComponents(a, weight)
	var vecs *mat.Dense
	var vars []float64
	if ok {
		vecs = pc.VectorsTo(vecs)
		vars = pc.VarsTo(vars)
		return &tablePCA{vars, vecs}
	} else {
		return nil
	}
	return nil
}
