package jsoncompare

import "fmt"

type ErrorPart string

const (
	JSON1 ErrorPart = "JSON1"
	JSON2 ErrorPart = "JSON2"
)

type JsonComparisonError struct {
	Part    ErrorPart
	Message string
}

func (j *JsonComparisonError) String() string {
	return fmt.Sprintf("%s: %s", j.Part, j.Message)
}

func AddUnsupportedType(errors []JsonComparisonError, part ErrorPart) []JsonComparisonError {
	errors = append(errors, JsonComparisonError{
		Part:    part,
		Message: "Invalid type. Only strings or structs are supported.",
	})
	return errors
}

func AddEmptyJsonStringError(errors []JsonComparisonError, part ErrorPart) []JsonComparisonError {
	errors = append(errors, JsonComparisonError{
		Part:    part,
		Message: "JSON string cannot be empty.",
	})
	return errors
}

func AddMissingField(errors []JsonComparisonError, part ErrorPart, fieldName string) []JsonComparisonError {
	oposite := map[ErrorPart]ErrorPart{
		JSON1: JSON2,
		JSON2: JSON1,
	}
	errors = append(errors, JsonComparisonError{
		Part:    part,
		Message: fmt.Sprintf("\"%s\" not found in %s.", fieldName, oposite[part]),
	})
	return errors
}
