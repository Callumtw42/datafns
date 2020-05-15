package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sort"
)

//Data an array of untyped maps
type Data []map[string]interface{}

//Empty returns true if data is empty
func (data Data) Empty() bool {
	if data == nil {
		fmt.Println("Data = null")
	}
	return data == nil || len(data) == 0
}

//GetColumn returns a column as a slice
func (data Data) GetColumn(colName string) (column []interface{}) {
	for _, e := range data {
		column = append(column, e[colName])
	}
	return column
}

//SumColumn returns the sum of a column
func (data Data) SumColumn(colName string) (total float64) {
	if !data.Empty() {
		var typ string = reflect.TypeOf(data[0][colName]).Name()
		if typ == "int64" || typ == "float64" {
			for _, e := range data {
				total += e[colName].(float64)
			}
		}
	}
	return total
}

//SetColumn returns a copy modifying an existing column by running function f on each value in an existing column
func (data Data) SetColumn(colName string, f func(e interface{}) interface{}) (cpy Data) {
	for _, e := range data {
		cpyE := deepCopy(e)
		cpyE[colName] = f(cpyE[colName])
		cpy = append(cpy, cpyE)
	}
	return cpy
}

//AddColumn returns a copy adding a new column by running the function f on each value in an existing column
func (data Data) AddColumn(colName string, newColName string, f func(interface{}) interface{}) Data {
	for _, e := range data {
		e[newColName] = f(e[colName])
	}
	return data
}

//RemoveColumns returns a copy with 1 or more columns removed
func (data Data) RemoveColumns(colNames ...string) Data {
	for _, e := range data {
		for _, colName := range colNames {
			delete(e, colName)
		}
	}
	return data
}

//Columns returns a copy with only the selected columns
func (data Data) Columns(colNames ...string) (subset Data) {
	for i, e := range data {
		subset = append(subset, map[string]interface{}{})
		for _, colName := range colNames {
			subset[i][colName] = e[colName]
		}
	}
	return subset
}

//GetUniqueValues returns an slice holding only unique values of a given column
func (data Data) GetUniqueValues(colName string) (vals []interface{}) {
	set := map[interface{}]bool{}
	for _, e := range data {
		if !set[e[colName]] {
			set[e[colName]] = true
			vals = append(vals, e[colName])
		}
	}
	return vals
}

//SplitByValues splits the data into seperate slices according to distinct values in the given column
func (data Data) SplitByValues(colName string) (split []Data) {
	valSet := map[interface{}]interface{}{}
	var valCount int
	for _, e := range data {
		var key interface{} = e[colName]
		if valSet[key] == nil {
			valSet[key] = valCount
			split = append(split, Data{e})
			valCount++
		} else {
			var splitItem *Data = &split[valSet[key].(int)]
			*splitItem = append(*splitItem, e)
		}
	}
	return split
}

// //SumAndGroup
// func (data Data) SumAndGroup(colName string) (subset Data) {

// 	valSet := map[interface{}]interface{}{}
// 	var valCount int

// 	var numeric []string
// 	var sampleRow map[string]interface{} = data[0]
// 	for k, v := range sampleRow {
// 		typ := reflect.TypeOf(v).Name()
// 		if typ == "float64" && k != colName {
// 			numeric = append(numeric, k)
// 		}
// 	}

// 	for _, e := range data {
// 		var key interface{} = e[colName]
// 		if valSet[key] == nil {
// 			valSet[key] = valCount
// 			subset = append(subset, e)
// 			valCount++
// 		} else {
// 			accRow := &subset[valSet[key].(int)]
// 			for _, k := range numeric {
// 				(*accRow)[k] = (*accRow)[k].(float64) + e[k].(float64)
// 			}
// 		}
// 	}
// 	return subset
// }

//SumAndGroup
func (data Data) SumAndGroup(colName string) (subset Data) {

	valSet := map[interface{}]interface{}{}
	var valCount int

	var numeric []string
	var sampleRow map[string]interface{} = data[0]
	for k, v := range sampleRow {
		typ := reflect.TypeOf(v).Name()
		if typ == "float64" && k != colName {
			numeric = append(numeric, k)
		}
	}

	for _, e := range data {
		var key interface{} = e[colName]
		if valSet[key] == nil {
			valSet[key] = valCount
			subset = append(subset, deepCopy(e))
			valCount++
		} else {
			accRow := &subset[valSet[key].(int)]
			for _, k := range numeric {
				(*accRow)[k] = (*accRow)[k].(float64) + e[k].(float64)
			}
		}
	}
	return subset
}

func deepCopy(m map[string]interface{}) map[string]interface{} {
	j, err := json.Marshal(m)
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	var c *map[string]interface{}
	json.Unmarshal(j, &c)
	return *c
}

func (data Data) getElementsWithValue(colName string, value interface{}) (rows Data) {
	for _, e := range data {
		if e[colName] == value {
			rows = append(rows, e)
		}
	}
	return rows
}

func (data Data) SortBy(colName string, order string) Data {
	conditions := [2]string{reflect.TypeOf(data[0][colName]).Name(), order}
	var less func(i, j int) bool
	switch conditions {
	case [2]string{"float64", "asc"}:
		less = func(i, j int) bool { return data[i][colName].(float64) < data[j][colName].(float64) }
	case [2]string{"float64", "desc"}:
		less = func(i, j int) bool { return data[i][colName].(float64) > data[j][colName].(float64) }
	case [2]string{"string", "asc"}:
		less = func(i, j int) bool { return data[i][colName].(string)[0] < data[j][colName].(string)[0] }
	case [2]string{"string", "desc"}:
		less = func(i, j int) bool { return data[i][colName].(string)[0] < data[j][colName].(string)[0] }
	}
	sort.Slice(data, less)
	return data
}

// function sort(data, property, order) {
//     return data.sort(sortByProperty(property, order))
// }

// function sortByProperty(property, order) {
//     order = (order === 'asc') ? -1 : 1;
//     return function (a, b) {
//         if (a[property] > b[property])
//             return order;
//         else if (a[property] < b[property])
//             return -order;
//         return 0;
//     }
// }
