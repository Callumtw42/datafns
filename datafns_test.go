package datafns

import (
	"reflect"
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
			if got := tt.data.Empty(); got != tt.want {
				t.Errorf("Data.Empty() = %v, want %v", got, tt.want)
			}
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
			if gotColumn := tt.data.getColumn(tt.args.colName); !reflect.DeepEqual(gotColumn, tt.wantColumn) {
				t.Errorf("Data.getColumn() = %v, want %v", gotColumn, tt.wantColumn)
			}
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
			gotTotal := tt.data.sumColumn(tt.args.colName)
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
			gotCopy := tt.data.setColumn(tt.args.colName, tt.args.f)
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
		{"empty", Data(nil), args{"n", "new", f},
			Data(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.data.addColumn(tt.args.colName, tt.args.newColName, tt.args.f)
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
			got := tt.data.removeColumns(tt.args.colNames...)
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
			got := tt.data.columns(tt.args.colNames...)
			assert.Equal(t, tt.want, got, "", "")
		})
	}
}
