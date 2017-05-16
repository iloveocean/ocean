// construct MongoDB query condition as key-value pair
package mongo

import (
	"encoding/json"
	"errors"

	"reflect"
)

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
	return m.M(key, value, operators...)
}

//condition version of M function
func MCond(key string, value interface{}, cond bool, operators ...Operator) Map {
	m := Map{}
	return m.MCond(key, value, cond, operators...)
}

//M method, construct query condition as the following styles:
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
func (m Map) M(key string, value interface{}, operators ...Operator) Map {
	if len(operators) > 0 {
		m.handleWithOperator(key, value, operators[0])
	} else {
		m.handleWithoutOperator(key, value)
	}
	return m
}

//condition version of M method
func (m Map) MCond(key string, value interface{}, cond bool, operators ...Operator) Map {
	if !cond {
		return m
	}
	return m.M(key, value, operators...)
}

//force construct the value in Map of the style of []interface{}
//to meet the syntax requiremnt of MongoDB of allowing multiple Or conditions
//in one pipeline document.
//MSlice can wrap multiple Or conditions into one AND conditon
func (m Map) MSlice(key string, value interface{}) {
	if old, ok := m[key]; ok {
		if oldIsSlice, oldValue := isSlice(old); oldIsSlice {
			if newIsSlice, newValue := isSlice(value); newIsSlice {
				m[key] = append(oldValue, newValue...)
			} else {
				m[key] = append(oldValue, value)
			}
		} else {
			if newIsSlice, newValue := isSlice(value); newIsSlice {
				tmp := []interface{}{old}
				m[key] = append(tmp, newValue...)
			} else {
				//force construct []Map
				m[key] = []interface{}{old, value}
			}
		}
	} else {
		if newIsSlice, _ := isSlice(value, true); newIsSlice {
			m[key] = value
		} else {
			m[key] = []interface{}{value}
		}
	}
}

func (m Map) Document() string {
	//docRaw, err := json.Marshal(m)
	docRaw, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		return errors.New("json marshal error: " + err.Error()).Error()
	}
	return string(docRaw)
}

//func (m Map) MSlice(key string, value interface{}, operators ...Operator) Map {

//}

//Merge two Maps into one
//Please note that both from Map and to One should not contain the
//same key(s), otherwise, the value in to Map will be replaced by the one
//in from One.
func (to Map) Merge(from Map) Map {
	for k, v := range from {
		to[k] = v
	}
	return to
}

func (m Map) handleWithoutOperator(key string, value interface{}) {
	if old, ok := m[key]; ok {
		if oldIsSlice, oldValue := isSlice(old); oldIsSlice {
			if newIsSlice, newValue := isSlice(value); newIsSlice {
				m[key] = append(oldValue, newValue...)
			}
		} else {
			m[key] = value
		}
	} else {
		m[key] = value
	}
}

func (m Map) handleWithOperator(key string, newOne interface{}, op Operator) {
	if old, ok := m[key]; ok {
		if mapValue, isMap := old.(Map); isMap {
			mapValue[string(op)] = newOne
			m[key] = mapValue
		} else if oldIsSlice, oldValue := isSlice(old); oldIsSlice {
			if newIsSlice, newValue := isSlice(newOne); newIsSlice {
				m[key] = append(oldValue, newValue...)
			} else {
				m[key] = append(oldValue, newOne)
			}
		}
	} else {
		m[key] = Map{string(op): newOne}
	}
}

//helper function to judge whether input value is slice type
//if yes, extract each elemnt
func isSlice(input interface{}, onlyJudge ...bool) (bool, []interface{}) {
	if input == nil {
		return false, nil
	}

	//check type
	value := reflect.ValueOf(input)
	if value.Kind() != reflect.Slice {
		return false, nil
	}
	if len(onlyJudge) > 0 && onlyJudge[0] {
		return true, nil
	}
	//extract value
	ret := make([]interface{}, value.Len())
	for i := 0; i < value.Len(); i++ {
		ret[i] = value.Index(i).Interface()
	}
	return true, ret
}
