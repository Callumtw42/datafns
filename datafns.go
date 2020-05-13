package datafns

import (
	"fmt"
	"reflect"
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

func (data Data) getColumn(colName string) (column []interface{}) {
	for _, e := range data {
		column = append(column, e[colName])
	}
	return column
}

func (data Data) sumColumn(colName string) (total float64) {
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

func (data Data) setColumn(colName string, f func(e interface{}) interface{}) Data {
	for _, e := range data {
		e[colName] = f(e[colName])
	}
	return data
}

func (data Data) addColumn(colName string, newColName string, f func(interface{}) interface{}) Data {
	for _, e := range data {
		e[newColName] = f(e[colName])
	}
	return data
}

func (data Data) removeColumns(colNames ...string) Data {
	for _, e := range data {
		for _, colName := range colNames {
			delete(e, colName)
		}
	}
	return data
}

func (data Data) columns(colNames ...string) (subset Data) {
	for i, e := range data {
		subset = append(subset, map[string]interface{}{})
		for _, colName := range colNames {
			subset[i][colName] = e[colName]
		}
	}
	return subset
}

// function sumAndGroup(data, col, ...dontSum) {
//     let groups = getUniqueValues(data, col);
//     let split = groups.map(e => { return getElementsWithValue(data, col, e) });
//     const sumObjectsByKey = (obj1, obj2) => {
//         Object.keys(obj1).forEach(k => { obj1[k] = (typeof obj1[k] === 'number' && k !== col && !dontSum.includes(k)) ? obj1[k] + obj2[k] : obj1[k] });
//         return obj1;
//     }
//     split = JSON.parse(JSON.stringify(split));
//     let grouped = split.map(a => { return a.reduce(sumObjectsByKey) });
//     return grouped;
// }

// function split(data, col) {
//     let groups = getUniqueValues(data, col);
//     return groups.map(e => { return getElementsWithValue(data, col, e) });
// }

// function getUniqueValues(data, col) {
//     return [...new Set(data.map(i => {
//         return i[col];
//     }))];
// }

// function getElementsWithValue(data, key, value) {
//     return data.filter(e =>
//         e[key] === value)
// }

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
