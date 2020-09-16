package jsoncompare

import (
	"encoding/json"
	"fmt"
	"reflect"
)

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
	errors = jsonStringMatches(jsonMap1, jsonMap2, JSON1, errors)
	errors = jsonStringMatches(jsonMap2, jsonMap1, JSON2, errors)
	return errors
}

func toJsonMap(json interface{}, part ErrorPart, errors []JsonComparisonError) (map[string]interface{}, []JsonComparisonError) {
	jsonType := getType(json)
	switch jsonType {
	case "string":
		return stringToJsonMap(json.(string), part, errors)
	case "struct":
		return interfaceToJsonMap(json, part, errors)
	}
	return nil, errors
}

func stringToJsonMap(j string, part ErrorPart, errors []JsonComparisonError) (map[string]interface{}, []JsonComparisonError) {
	if j == "" {
		errors = AddEmptyJsonStringError(errors, part)
		return nil, errors
	}
	var result map[string]interface{}
	json.Unmarshal([]byte(j), &result)
	return result, errors
}

func interfaceToJsonMap(j interface{}, part ErrorPart, errors []JsonComparisonError) (map[string]interface{}, []JsonComparisonError) {
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

func jsonStringMatches(jsonMap, otherJsonMap map[string]interface{}, part ErrorPart, errors []JsonComparisonError) []JsonComparisonError {
	if len(jsonMap) != len(otherJsonMap) {
		// return false
	}
	for key, value := range jsonMap {
		otherValue, found := otherJsonMap[key]
		if !found {
			errors = AddMissingField(errors, part, key)
			return errors
		}
		valueType := fmt.Sprintf("%T", value)
		otherValueType := fmt.Sprintf("%T", otherValue)
		if valueType != otherValueType {
			return errors
		}
		switch valueType {
		case "map[string]interface {}": // object
			lenBefore := len(errors)
			errors = jsonStringMatches(jsonMap[key].(map[string]interface{}), otherJsonMap[key].(map[string]interface{}), part, errors)
			lenAfter := len(errors)
			if lenBefore != lenAfter {
				return errors
			}
			continue
		case "[]interface {}": // repeated object
			// naive implementation of comparison of repeated messages.
			// TODO investigate a more performant way to compare
			items := jsonMap[key].([]interface{})
			otherItems := otherJsonMap[key].([]interface{})
			if len(items) != len(otherItems) {
				return errors
			}
			for _, item := range items {
				var found = false
				for _, otherItem := range otherItems {
					itemType := fmt.Sprintf("%T", item)
					otherItemType := fmt.Sprintf("%T", otherItem)
					if itemType != otherItemType {
						// Not sure if they can be different
						continue
					}
					switch itemType {
					case "map[string]interface {}":
						lenBefore := len(errors)
						errors = jsonStringMatches(item.(map[string]interface{}), otherItem.(map[string]interface{}), part, errors)
						lenAfter := len(errors)
						if lenBefore != lenAfter {
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
					return errors
				}
			}
			continue
		}
		if value != otherValue {
			return errors
		}
	}
	return errors
}
