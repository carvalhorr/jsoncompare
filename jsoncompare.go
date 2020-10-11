package jsoncompare

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func IsEqual(json1, json2 interface{}) bool {
	errors := Compare(json1, json2)
	return len(errors) == 0
}

// Compare two JSONs and return an array of JsonComparisonError with the differences.
// The two JSONs can be either string or a struct.
func Compare(json1, json2 interface{}) []JsonComparisonError {
	var errors []JsonComparisonError
	isType1Supported := isTypeSupported(json1)
	isType2Supported := isTypeSupported(json2)
	if !isType1Supported {
		errors = AddUnsupportedType(errors, JSON1)
	}
	if !isType2Supported {
		errors = AddUnsupportedType(errors, JSON2)
	}
	if !isType1Supported || !isType2Supported {
		return errors
	}

	return compare(json1, json2, errors)
}
func compare(json1, json2 interface{}, errors []JsonComparisonError) []JsonComparisonError {
	jsonMap1, errors := toJsonMap(json1, JSON1, errors)
	jsonMap2, errors := toJsonMap(json2, JSON2, errors)
	if len(errors) > 0 {
		return errors
	}
	errors = jsonStringMatches(jsonMap1, jsonMap2, JSON1, "", errors)
	errors = jsonStringMatches(jsonMap2, jsonMap1, JSON2, "", errors)
	return errors
}

func toJsonMap(json interface{}, part ErrorPart, errors []JsonComparisonError) (interface{}, []JsonComparisonError) {
	jsonType := getType(json)
	switch jsonType {
	case "string":
		return stringToJsonMap(json.(string), part, errors)
	case "struct":
		return interfaceToJsonMap(json, part, errors)
	}
	return nil, errors
}

func stringToJsonMap(j string, part ErrorPart, errors []JsonComparisonError) (interface{}, []JsonComparisonError) {
	if j == "" {
		errors = AddEmptyJsonStringError(errors, part)
		return nil, errors
	}
	var result interface{}
	json.Unmarshal([]byte(j), &result)
	return result, errors
}

func jsonMapToString(m map[string]interface{}) string {
	result, _ := json.Marshal(m)
	// TODO handle error
	return string(result)
}

func interfaceToJsonMap(j interface{}, part ErrorPart, errors []JsonComparisonError) (interface{}, []JsonComparisonError) {
	bytes, _ := json.Marshal(j)
	// TODO handle error
	return stringToJsonMap(string(bytes), part, errors)
}

func isTypeSupported(json interface{}) bool {
	jsonType := getType(json)
	switch jsonType {
	case "string":
		return true
	case "struct":
		return true
	}
	return false
}

func getType(json interface{}) string {
	return reflect.ValueOf(json).Kind().String()
}

func jsonStringMatches(json1, json2 interface{}, part ErrorPart, path string, errors []JsonComparisonError) []JsonComparisonError {
	json1Type := fmt.Sprintf("%T", json1)
	switch json1Type {
	case "map[string]interface {}":
		return matchesObject(json1.(map[string]interface{}), json2.(map[string]interface{}), part, path, errors)
	case "[]interface {}":
		return matchesArray(json1.([]interface{}), json2.([]interface{}), part, path, errors)
	default:
		return errors
	}
}

func matchesObject(jsonMap, otherJsonMap map[string]interface{}, part ErrorPart, path string, errors []JsonComparisonError) []JsonComparisonError {
	for key, value := range jsonMap {
		fullPath := getFullPath(path, key)
		otherValue, found := otherJsonMap[key]
		if !found {
			errors = AddMissingField(errors, part, fullPath)
			continue
		}
		valueType := fmt.Sprintf("%T", value)
		switch valueType {
		case "map[string]interface {}": // object
			errors = jsonStringMatches(jsonMap[key].(map[string]interface{}), otherJsonMap[key].(map[string]interface{}), part, fullPath, errors)
		case "[]interface {}": // repeated object
			errors = matchesArray(jsonMap[key].([]interface{}), otherJsonMap[key].([]interface{}), part, fullPath, errors)
		default:
			errors = matchesTypeAndValue(part, fullPath, errors, value, otherValue, valueType)
		}
	}
	return errors
}

func matchesArray(items []interface{}, otherItems []interface{}, part ErrorPart, fullPath string, errors []JsonComparisonError) []JsonComparisonError {
	for i, item := range items {
		var found = false
		item1Type := fmt.Sprintf("%T", item)
		for _, otherItem := range otherItems {
			switch item1Type {
			case "map[string]interface {}":
				if IsEqual(jsonMapToString(item.(map[string]interface{})), jsonMapToString(otherItem.(map[string]interface{}))) {
					found = true
					break
				}
			default:
				if item == otherItem {
					found = true
					break
				}
			}
		}
		if !found {
			errors = AddItemNotFoundArray(errors, part, fullPath, i)
		}
	}
	return errors
}

func matchesTypeAndValue(part ErrorPart, fullPath string, errors []JsonComparisonError, value, otherValue interface{}, valueType string) []JsonComparisonError {
	otherValueType := fmt.Sprintf("%T", otherValue)
	if valueType != otherValueType {
		errors = AddTypeMismatchField(errors, part, fullPath, valueType, otherValueType)
		return errors
	}
	if value != otherValue {
		errors = AddDifferentValue(errors, part, fullPath, value, otherValue)
	}
	return errors
}

func getFullPath(path, field string) string {
	if path == "" {
		return field
	}
	return path + "." + field
}
