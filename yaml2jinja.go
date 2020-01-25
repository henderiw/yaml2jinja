package yaml2jinja

import (
	"fmt"
	"go/format"
	"reflect"
	"strings"

	"gopkg.in/yaml.v1"
)

// New creates Yaml2Jinja object
func New() Yaml2Jinja {
	return Yaml2Jinja{}
}

type line struct {
	structName string
	line       string
}

// Yaml2Jinja to store converted result
type Yaml2Jinja struct {
	Visited   map[line]bool
	StructMap map[string]string
}

// NewStruct creates new entry in StructMap result
func (yj *Yaml2Jinja) NewStruct(structName string, parent string) string {
	// If struct already present with the same name
	// rename struct to ParentStructname
	if _, ok := yj.StructMap[structName]; ok {
		structName = goKeyFormat(parent) + structName
	}
	yj.AppendResult(structName, fmt.Sprintf("// %s\n", structName))
	l := fmt.Sprintf("type %s struct {\n", structName)
	yj.Visited[line{structName, l}] = true
	yj.StructMap[structName] += l
	return structName
}

// AppendResult add lines to the result
func (yj *Yaml2Jinja) AppendResult(structName string, l string) {
	if _, ok := yj.Visited[line{structName, l}]; !ok {
		yj.StructMap[structName] += l
	}
	yj.Visited[line{structName, l}] = true
}

// removeUnderscores and camelize string
func goKeyFormat(key string) string {
	var st string
	strList := strings.Split(key, "_")
	for _, str := range strList {
		st += strings.Title(str)
	}
	if len(st) == 0 {
		st = key
	}
	return st
}

// Convert transforms map[string]interface{} to go struct
func (yj *Yaml2Jinja) Convert(structName string, data []byte) (string, error) {
	yj.Visited = make(map[line]bool)
	yj.StructMap = make(map[string]string)

	// Unmarshal to map[string]interface{}
	var obj map[string]interface{}
	err := yaml.Unmarshal(data, &obj)
	fmt.Println(err)
	if err != nil {
		return "", err
	}
	fmt.Println(obj)

	yj.NewStruct("Yaml2Jinja", "")
	for k, v := range obj {
		fmt.Printf("Key: %v \n", k)
		fmt.Printf("Value: %v \n", v)
		yj.Structify(structName, k, v, false)
	}
	yj.AppendResult("Yaml2Go", "}\n")

	var result string
	for _, value := range yj.StructMap {
		result += fmt.Sprintf("%s\n", value)
	}

	// Convert result into go format
	goFormat, err := format.Source([]byte(result))
	if err != nil {
		return "", err
	}
	return string(goFormat), nil
}

// Structify transforms map key values to struct fields
// structName : parent struct name
// k, v       : fields in the struct
func (yj *Yaml2Jinja) Structify(structName, k string, v interface{}, arrayElem bool) {

	if reflect.TypeOf(v) == nil || len(k) == 0 {
		yj.AppendResult(structName, fmt.Sprintf("%s interface{} `yaml:\"%s\"`\n", goKeyFormat(k), k))
		return
	}

	switch reflect.TypeOf(v).Kind() {

	// If yaml object
	case reflect.Map:
		switch val := v.(type) {
		case map[interface{}]interface{}:
			key := goKeyFormat(k)
			newKey := key
			if !arrayElem {
				// Create new structure
				newKey = yj.NewStruct(key, structName)
				yj.AppendResult(structName, fmt.Sprintf("%s %s `yaml:\"%s\"`\n", key, newKey, k))
			}
			// If array of yaml objects
			for k1, v1 := range val {
				if _, ok := k1.(string); ok {
					yj.Structify(newKey, k1.(string), v1, false)
				}
			}
			if !arrayElem {
				yj.AppendResult(newKey, "}\n")
			}
		}

	// If array
	case reflect.Slice:
		val := v.([]interface{})
		if len(val) == 0 {
			return
		}
		switch val[0].(type) {

		case string, int, bool, float64:
			yj.AppendResult(structName, fmt.Sprintf("%s []%s `yaml:\"%s\"`\n", goKeyFormat(k), reflect.TypeOf(val[0]), k))

		// if nested object
		case map[interface{}]interface{}:
			key := goKeyFormat(k)
			// Create new structure
			newKey := yj.NewStruct(key, structName)
			yj.AppendResult(structName, fmt.Sprintf("%s []%s `yaml:\"%s\"`\n", key, newKey, k))
			for _, v1 := range val {
				yj.Structify(newKey, key, v1, true)
			}
			yj.AppendResult(newKey, "}\n")
		}

	default:
		yj.AppendResult(structName, fmt.Sprintf("%s %s `yaml:\"%s\"`\n", goKeyFormat(k), reflect.TypeOf(v).String(), k))
	}
}
