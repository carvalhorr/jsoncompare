package jsoncompare

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompare_TwoEmptyJson(t *testing.T) {
	jsonStr1 := ""
	jsonStr2 := ""
	errors := Compare(jsonStr1, jsonStr2)
	assert.Equal(t, 2, len(errors), "Empty JSON should return error")
	assert.Equal(t, "JSON1: JSON string cannot be empty.", errors[0].String(), "First JSON empty should return error")
	assert.Equal(t, "JSON2: JSON string cannot be empty.", errors[1].String(), "Second JSON empty should return error")
}

func TestCompare_FirstJsonIsEmpty(t *testing.T) {
	jsonStr1 := ""
	jsonStr2 := "{}"
	errors := Compare(jsonStr1, jsonStr2)
	assert.Equal(t, 1, len(errors), "Empty JSON should return error")
	assert.Equal(t, "JSON1: JSON string cannot be empty.", errors[0].String(), "First JSON empty should return error")
}

func TestCompare_SecondJsonIsEmpty(t *testing.T) {
	jsonStr1 := "{}"
	jsonStr2 := ""
	errors := Compare(jsonStr1, jsonStr2)
	assert.Equal(t, 1, len(errors), "Empty JSON should return error")
	assert.Equal(t, "JSON2: JSON string cannot be empty.", errors[0].String(), "Second JSON empty should return error")
}

func TestCompare_StructsEqual(t *testing.T) {
	jsonStruct1 := TestStruct{Field1: "blah1"}
	jsonStruct2 := TestStruct{Field1: "blah1"}
	errors := Compare(jsonStruct1, jsonStruct2)
	assert.Equal(t, 0, len(errors))
}

func TestCompare_UnsupportedJson1(t *testing.T) {
	json1 := float64(4)
	json2 := ""
	errors := Compare(json1, json2)
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, JSON1, errors[0].Part)
	assert.Equal(t, "Invalid type. Only strings or structs are supported.", errors[0].Message)
}

type TestStruct struct {
	Field1 string `json:"field1"`
}
