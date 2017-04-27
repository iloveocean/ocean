// construct MongoDB query condition as key-value pair
package mongo

import "reflect"

//using map to storage and present key-value pair
type Map map[string]interface{}

//M function, construct query condition as the following styles:
///////////////////////////////
//    {key: value}
//    {key: {operator: value}}
//////////////////////////////
//please not that, value and can be nested
//that is value can be another Map or Map slice
////////////////////////////////////////////////////
//please also note the following two principles:
///////////////////////////////////////////////////
//1. for multiple key-value pairs whose key are the same one, M will always keep the latest one
//2. If value is a slice style, the old value must also be the slice type! And each elemnt in the new slice
//   will be appended to the old one, if either of old value and new one isn't slice, the new value will not
//   be stored in Map
func M(key string, value interface{}, operators ...Operator) Map {
	m := Map{}
	return m.M(key, value, operators)
}

func (this *Map) M(key string, value interface{}, operators ...Operator) Map {
	if len(operators) > 0 {
		this.handleWithOperator(key, value, operators[0])
	} else {
		this.handleWithoutOperator(key, value)
	}
	return this
}

func (this *Map) handleWithoutOperator(key string, value interface{}) {
	//update directly
	this[key] = value
}

func (this *Map) handleWithOperator(string key, newOne interface{}, op Operator) {
	if old, ok := this[key]; ok {
		if mapValue, isMap := Map(old); isMap {
			mapValue[string(op)] = newOne
		} else if isSliceOld, oldValue := isSlice(old); isSliceOld {
			if isSliceNew, newValue := isSlice(newOne); isSliceNew {
				old[key] = append(oldValue, newValue...)
			}
		}
	} else {
		tmp := Map{
			string(op): newOne,
		}
		this[key] = tmp
	}
}

//helper function to judge whether input value is slice type
//if yes, extract each elemnt
func isSlice(input interface{}) (bool, []interface{}) {
	if input == nil {
		return false, nil
	}

	//check type
	value := reflect.ValueOf(input)
	if value.Kind() != reflect.Slice {
		return false, nil
	}

	//extract value
	ret := make([]interface{}, value.Len())
	for i := 0; i < value.Len(); i++ {
		ret[i] = value.Index(i).Interface()
	}
	return true, ret
}
