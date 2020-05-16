package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func populated() []map[string]interface{} {

	var data []map[string]interface{}

	data = append(data, map[string]interface{}{
		"zero":     0,
		"posInt":   float64(1),
		"negInt":   float64(-1),
		"posFloat": 1.5,
		"negFloat": -1.5,
		"string":   "apple",
		"null":     nil,
	})
	data = append(data, map[string]interface{}{
		"zero":     0,
		"posInt":   float64(1),
		"negInt":   float64(-1),
		"posFloat": 1.5,
		"negFloat": -1.5,
		"string":   "banana",
		"null":     nil,
	})
	data = append(data, map[string]interface{}{
		"zero":     0,
		"posInt":   float64(1),
		"negInt":   float64(-1),
		"posFloat": 1.5,
		"negFloat": -1.5,
		"string":   "cranberry",
		"null":     nil,
	})
	return data
}

var _2x2 Data = Data{
	{"i": float64(-1), "s": "abc"},
	{"n": float64(1), "s": "def"},
}

var _2x3 Data = Data{
	{"n": float64(-1), "s": "abc"},
	{"n": float64(0), "s": "def"},
	{"n": float64(1), "s": "ghi"},
}

var _3x3 Data = Data{
	{"n": float64(-1), "f": float64(-1.5), "s": "abc"},
	{"n": float64(0), "f": float64(0.0), "s": "def"},
	{"n": float64(1), "f": float64(1.5), "s": "ghi"},
}

var data []map[string]interface{} = populated()

func TestData_Empty(t *testing.T) {
	tests := []struct {
		name string
		data Data
		want bool
	}{
		{"populated", data, false},
		{"null", nil, true},
		{"empty", Data(nil), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.data.empty()
			got := tt.data.Empty()
			assert.Equal(t, tt.want, got, "", "")
		})
	}
}

func TestData_getColumn(t *testing.T) {
	type args struct {
		colName string
	}
	tests := []struct {
		name       string
		data       Data
		args       args
		wantColumn []interface{}
	}{
		{"populated", data, args{"string"}, []interface{}{"apple", "banana", "cranberry"}},
		{"empty", Data(nil), args{"x"}, []interface{}(nil)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.data.GetColumn(tt.args.colName)
			gotColumn := tt.data.GetColumn(tt.args.colName)
			assert.Equal(t, tt.wantColumn, gotColumn, "", "")
		})
	}
}

func TestData_sumColumn(t *testing.T) {
	type args struct {
		colName string
	}
	tests := []struct {
		name      string
		data      Data
		args      args
		wantTotal float64
	}{

		{"populated 1", data, args{"posInt"}, 3},
		{"populated -1", data, args{"negInt"}, -3},
		{"populated 0 ", data, args{"zero"}, 0},
		{"populated 1.5 ", data, args{"posFloat"}, 4.5},
		{"populated -1.5 ", data, args{"negFloat"}, -4.5},
		{"populated string ", data, args{"string"}, 0},
		{"empty", Data(nil), args{"string"}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.data.SumColumn(tt.args.colName)
			gotTotal := tt.data.SumColumn(tt.args.colName)
			assert.Equal(t, tt.wantTotal, gotTotal, "", "")
		})
	}
}

func TestData_setColumn(t *testing.T) {
	data2 := Data{
		{"n": float64(1), "s": "abc"},
		{"n": float64(2), "s": "def"},
	}
	f := func(v interface{}) interface{} {
		return v.(float64) * 2
	}
	type args struct {
		colName string
		f       func(interface{}) interface{}
	}
	tests := []struct {
		name     string
		data     Data
		args     args
		wantCopy Data
	}{
		{"mult2", data2, args{"n",
			f}, Data{
			{"n": float64(2), "s": "abc"},
			{"n": float64(4), "s": "def"},
		},
		},
		{"empty", Data(nil), args{"n", f}, Data(nil)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.data.SetColumn(tt.args.colName, tt.args.f)
			gotCopy := tt.data.SetColumn(tt.args.colName, tt.args.f)
			assert.Equal(t, tt.wantCopy, gotCopy, "", "")
		})
	}
}

func TestData_addColumn(t *testing.T) {
	f := func(v interface{}) interface{} {
		return v.(float64) * 2
	}
	data := Data{
		{"n": float64(1), "s": "abc"},
		{"n": float64(2), "s": "def"},
	}
	type args struct {
		colName    string
		newColName string
		f          func(interface{}) interface{}
	}
	tests := []struct {
		name string
		data Data
		args args
		want Data
	}{
		{"mult2", data, args{"n", "new", f},
			Data{
				{"n": float64(1), "s": "abc", "new": float64(2)},
				{"n": float64(2), "s": "def", "new": float64(4)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.data.AddColumn(tt.args.colName, tt.args.newColName, tt.args.f)
			got := tt.data.AddColumn(tt.args.colName, tt.args.newColName, tt.args.f)
			assert.Equal(t, tt.want, got, "", "")
		})
	}
}

func TestData_removeColumns(t *testing.T) {
	type args struct {
		colNames []string
	}
	data := Data{
		{"n": float64(1), "s": "abc"},
		{"n": float64(2), "s": "def"},
	}
	tests := []struct {
		name string
		data Data
		args args
		want Data
	}{
		{"remove string", data, args{[]string{"s"}},
			Data{
				{"n": float64(1)},
				{"n": float64(2)},
			},
		},
		{"remove all", data, args{[]string{"s", "n"}},
			Data{{}, {}},
		},
		{"empty", Data(nil), args{[]string{"s"}},
			Data(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.data.RemoveColumns(tt.args.colNames...)
			got := tt.data.RemoveColumns(tt.args.colNames...)
			assert.Equal(t, tt.want, got, "", "")
		})
	}
}

func TestData_columns(t *testing.T) {
	type args struct {
		colNames []string
	}
	data := Data{
		{"n": float64(1), "s": "abc", "f": 1.5},
		{"n": float64(2), "s": "def", "f": 2.5},
	}
	tests := []struct {
		name string
		data Data
		args args
		want Data
	}{
		{"n", data, args{[]string{"n", "f"}}, Data{
			{"n": float64(1), "f": 1.5},
			{"n": float64(2), "f": 2.5},
		}},
		{"wrong", data, args{[]string{"x"}}, Data{{"x": nil}, {"x": nil}}},
		{"empty", Data(nil), args{[]string{"x"}}, Data(nil)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.data.Columns(tt.args.colNames...)
			got := tt.data.Columns(tt.args.colNames...)
			assert.Equal(t, tt.want, got, "", "")
		})
	}
}

func TestData_getUniqueValues(t *testing.T) {
	type f = float64
	data = Data{
		{"n": f(1), "s": "x"},
		{"n": f(1), "s": "x"},
		{"n": f(2), "s": "y"},
		{"n": f(2), "s": "y"},
	}
	type args struct {
		colName string
	}
	tests := []struct {
		name     string
		data     Data
		args     args
		wantVals []interface{}
	}{
		{"string", data, args{"s"}, []interface{}{
			"x", "y",
		}},

		{"float", data, args{"n"}, []interface{}{
			f(1), f(2),
		}},

		{"wrong", data, args{"banana"}, []interface{}{nil}},
		{"empty", Data(nil), args{"n"}, []interface{}(nil)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.data.GetUniqueValues(tt.args.colName)
			gotVals := tt.data.GetUniqueValues(tt.args.colName)
			assert.Equal(t, tt.wantVals, gotVals, "", "")
		})
	}
}

func TestData_splitByValues(t *testing.T) {
	type args struct {
		colName string
	}
	data = Data{
		{"n": float64(3), "s": "a"},
		{"n": float64(2), "s": "a"},
		{"n": float64(1), "s": "a"},
		{"n": float64(1), "s": "b"},
		{"n": float64(2), "s": "b"},
		{"n": float64(3), "s": "b"},
	}
	expTestInt := []Data{
		Data{
			{"n": float64(3), "s": "a"},
			{"n": float64(3), "s": "b"},
		},
		Data{
			{"n": float64(2), "s": "a"},
			{"n": float64(2), "s": "b"},
		},
		Data{
			{"n": float64(1), "s": "a"},
			{"n": float64(1), "s": "b"},
		},
	}
	expTestString := []Data{
		Data{
			{"n": float64(3), "s": "a"},
			{"n": float64(2), "s": "a"},
			{"n": float64(1), "s": "a"},
		},
		Data{
			{"n": float64(1), "s": "b"},
			{"n": float64(2), "s": "b"},
			{"n": float64(3), "s": "b"},
		},
	}
	tests := []struct {
		name      string
		data      Data
		args      args
		wantSplit []Data
	}{
		{"test int", data, args{"n"}, expTestInt},
		{"test string", data, args{"s"}, expTestString},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.data.SplitByValues(tt.args.colName)
			gotSplit := tt.data.SplitByValues(tt.args.colName)
			assert.Equal(t, tt.wantSplit, gotSplit, "", "")
		})
	}
}

func TestData_sumAndGroup(t *testing.T) {
	data = Data{
		{"x": float64(2), "y": float64(4), "s": "a"},
		{"x": float64(2), "y": float64(3), "s": "b"},
		{"x": float64(3), "y": float64(2), "s": "a"},
		{"x": float64(3), "y": float64(1), "s": "b"},
	}
	expS := Data{
		{"x": float64(5), "y": float64(6), "s": "a"},
		{"x": float64(5), "y": float64(4), "s": "b"},
	}
	expX := Data{
		{"x": float64(2), "y": float64(7), "s": "a"},
		{"x": float64(3), "y": float64(3), "s": "a"},
	}
	type args struct {
		colName string
	}
	tests := []struct {
		name       string
		data       Data
		args       args
		wantSubset Data
	}{
		{"group by x", data, args{"x"}, expX},
		{"group by s", data, args{"s"}, expS},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.data.SumAndGroup(tt.args.colName)
			gotSubset := tt.data.SumAndGroup(tt.args.colName)
			assert.Equal(t, tt.wantSubset, gotSubset, "", "")
		})
	}
}

func TestData_getElementsWithValue(t *testing.T) {
	data = Data{
		{"x": float64(1), "s": "a"},
		{"x": float64(0), "s": "a"},
		{"x": float64(1), "s": "b"},
		{"x": float64(0), "s": "b"},
	}
	exp := Data{
		{"x": float64(1), "s": "a"},
		{"x": float64(1), "s": "b"},
	}
	type args struct {
		colName string
		value   interface{}
	}
	tests := []struct {
		name     string
		data     Data
		args     args
		wantRows Data
	}{
		{"get 1", data, args{"x", float64(1)}, exp},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.data.getElementsWithValue(tt.args.colName, tt.args.value)
			gotRows := tt.data.getElementsWithValue(tt.args.colName, tt.args.value)
			assert.Equal(t, tt.wantRows, gotRows, "", "")
		})
	}
}

func TestData_sortBy(t *testing.T) {
	data = Data{
		{"x": float64(4), "y": "b"},
		{"x": float64(1), "y": "a"},
		{"x": float64(3), "y": "da"},
		{"x": float64(2), "y": "ca"},
	}

	exp1 := Data{
		{"x": float64(1), "y": "a"},
		{"x": float64(2), "y": "ca"},
		{"x": float64(3), "y": "da"},
		{"x": float64(4), "y": "b"},
	}

	exp2 := Data{
		{"x": float64(4), "y": "b"},
		{"x": float64(3), "y": "da"},
		{"x": float64(2), "y": "ca"},
		{"x": float64(1), "y": "a"},
	}

	exp3 := Data{
		{"x": float64(1), "y": "a"},
		{"x": float64(4), "y": "b"},
		{"x": float64(2), "y": "ca"},
		{"x": float64(3), "y": "da"},
	}

	type args struct {
		colName string
		order   string
	}
	tests := []struct {
		name string
		data Data
		args args
		want Data
	}{
		{"num asc", data, args{"x", "asc"}, exp1},
		{"num desc", data, args{"x", "desc"}, exp2},
		{"string", data, args{"y", "asc"}, exp3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.data.SortBy(tt.args.colName, tt.args.order)
			got := tt.data.SortBy(tt.args.colName, tt.args.order)
			assert.Equal(t, tt.want, got, "", "")
		})
	}
}
