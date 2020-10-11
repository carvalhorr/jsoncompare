package jsoncompare

import (
	"fmt"
)

type ErrorPart string

const (
	JSON1 ErrorPart = "JSON1"
	JSON2 ErrorPart = "JSON2"
)

type ErrorType string

const (
	UnsupportedType  = "UnsupportedType"
	EmptyJsonString  = "EmptyJsonString"
	MissingField     = "MissingField"
	TypeMismatch     = "TypeMismatch"
	ValueMismatch    = "ValueMismatch"
	MissingArrayItem = "MissingArrayItem"
)

type JsonComparisonError struct {
	Part    ErrorPart
	Type    ErrorType
	Field   string
	Message string
}

func (j *JsonComparisonError) String() string {
	return fmt.Sprintf("%s: %s", j.Part, j.Message)
}

func AddUnsupportedType(errors []JsonComparisonError, part ErrorPart) []JsonComparisonError {
	errors = append(errors, JsonComparisonError{
		Part:    part,
		Type:    UnsupportedType,
		Message: "Invalid type. Only strings or structs are supported.",
	})
	return errors
}

func AddEmptyJsonStringError(errors []JsonComparisonError, part ErrorPart) []JsonComparisonError {
	errors = append(errors, JsonComparisonError{
		Part:    part,
		Type:    EmptyJsonString,
		Message: "JSON string cannot be empty.",
	})
	return errors
}

var oppositeJson = map[ErrorPart]ErrorPart{
	JSON1: JSON2,
	JSON2: JSON1,
}

func AddMissingField(errors []JsonComparisonError, part ErrorPart, fieldName string) []JsonComparisonError {
	errors = append(errors, JsonComparisonError{
		Part:    part,
		Type:    MissingField,
		Field:   fieldName,
		Message: fmt.Sprintf("\"%s\" not found in %s.", fieldName, oppositeJson[part]),
	})
	return errors
}

func AddTypeMismatchField(errors []JsonComparisonError, part ErrorPart, fieldName string, type1 string, type2 string) []JsonComparisonError {
	var errType ErrorType = TypeMismatch
	if alreadyAddedInOtherSide(errors, part, errType, fieldName) {
		return errors
	}
	errors = append(errors, JsonComparisonError{
		Part:    part,
		Type:    errType,
		Field:   fieldName,
		Message: fmt.Sprintf("\"%s\" type mismatch in %s. Expected %s. Found %s.", fieldName, oppositeJson[part], type1, type2),
	})
	return errors
}

func AddDifferentValue(errors []JsonComparisonError, part ErrorPart, fieldName string, value1 interface{}, value2 interface{}) []JsonComparisonError {
	var errType ErrorType = ValueMismatch
	if alreadyAddedInOtherSide(errors, part, errType, fieldName) {
		return errors
	}
	errors = append(errors, JsonComparisonError{
		Part:    part,
		Type:    errType,
		Field:   fieldName,
		Message: fmt.Sprintf("\"%s\" value mismatch in %s. Expected %s. Found %s.", fieldName, oppositeJson[part], fmt.Sprint(value1), fmt.Sprint(value2)),
	})
	return errors
}

func AddItemNotFoundArray(errors []JsonComparisonError, part ErrorPart, fieldName string, position int) []JsonComparisonError {
	var errType ErrorType = MissingArrayItem
	if alreadyAddedInOtherSide(errors, part, errType, fieldName) {
		return errors
	}
	errors = append(errors, JsonComparisonError{
		Part:    part,
		Type:    errType,
		Field:   fieldName,
		Message: fmt.Sprintf("\"%s[%d]\" no corresponding item found in %s.", fieldName, position, oppositeJson[part]),
	})
	return errors
}

func alreadyAddedInOtherSide(errors []JsonComparisonError, part ErrorPart, errorType ErrorType, fieldName string) bool {
	for _, err := range errors {
		if err.Field == fieldName && err.Type == errorType {
			return true
		}
	}
	return false
}
