# jsoncompare

A small library to compare if two JSON's are the same. It takes either JSON strings or structs and return a list of differences.

## Functions
```

Compare
- Takes two jSON's and return list of differences.

```

```
IsEqual
- A wrapper around Compare that returns true if no differences are found.

```

## List of differences

* Unsupported Type

Returned if either of the JSON's is not a string or a struct.

* Empty JSON String

Returned if an empty JSON string is passed as a parameter.

* Missing Field

Returned if a field is only exist in one of the JSONs.

* Type Mismatch

Returned when fields exist in both JSONs, but they are of different types. 

* Value Mismatch

Returned when fields exist in both JSONs, they are of the same type, but they don have the same value.

* Missing Array Item

Returned when one item in an array only exist in one of the JSONs.

## Example

```go
type TestStruct struct {
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
	SubStruct SubStruct `json:"sub_struct"`
}

type SubStruct struct {
	Field2 float64 `json:"field2"`
	Array1 []string `json:"array1"`
}

func main() {
	testJson1 := TestStruct{
		Field1:    "value1",
		Field2:    "value2",
		SubStruct: SubStruct{
			Field2: 0,
			Array1: []string{"str1", "str2"},
		},
	}
	testJson2 := "{\"field1\": \"value1\", \"field2\": 2, \"sub_struct\":{\"field2\": 1, \"array1\":[\"str1\", \"str2\", \"str3\"]}, \"field3\": 3}"
	errs := Compare(testJson1, testJson2)
	for _, e := range errs {
		fmt.Println(e.String())
	}
}
```

Expected results:
```
JSON1: "field2" type mismatch in JSON2. Expected string. Found float64.
JSON1: "sub_struct.field2" value mismatch in JSON2. Expected 0. Found 1.
JSON2: "field3" not found in JSON1.
JSON2: "sub_struct.array1[2]" no corresponding item found in JSON1.
```