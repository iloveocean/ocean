// construct MongoDB query condition as key-value pair
package mongo

import "reflect"

//using map to storage and present key-value pair
type Map map[string]interface{}

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
