package main

import (
	"fmt"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Log variable
var Log = log.New()

// LogCtxt variable
var LogCtxt = Log.WithFields(log.Fields{
	"tag": "yaml2jinja",
})

var globalForIndent = 0

// New creates Yaml2Jinja object
func New() Yaml2Jinja {
	return Yaml2Jinja{}
}

// Yaml2Jinja to store converted result
type Yaml2Jinja struct {
	FileName  string
	Output    string
	Depth     int
	ForIndent int
	ForDepth  []int
	Index     []string
	NextIndex string
}

//Indent ...
func (yj *Yaml2Jinja) Indent() {
	var ident int
	forIndent := 0
	if globalForIndent < yj.ForIndent { // entered for loop; add - sign/ident
		forIndent = 1
	} else if globalForIndent > yj.ForIndent { // escaped for loop; remove for sign
		globalForIndent--
		forIndent = -1
	}

	ident = yj.Depth + globalForIndent
	if forIndent == 1 {
		globalForIndent++ // increase globalForIndent after for the next processing cycle
	}

	switch {
	case ident == 1:
		yj.Output += "  "
	case ident == 2:
		yj.Output += "    "
	case ident == 3:
		yj.Output += "      "
	case ident == 4:
		yj.Output += "        "
	case ident == 5:
		yj.Output += "          "
	case ident == 6:
		yj.Output += "            "
	case ident == 7:
		yj.Output += "              "
	case ident == 8:
		yj.Output += "                "
	case ident == 9:
		yj.Output += "                  "
	case ident == 10:
		yj.Output += "                    "
	default:
	}

	// append - indent at the end
	if forIndent == 1 {
		yj.Output += "- "
	}

}

//ConcatenateKey sr=tring
func (yj *Yaml2Jinja) ConcatenateKey() string {
	var cKey string
	var start int
	if len(yj.ForDepth) > 0 {
		start = yj.ForDepth[len(yj.ForDepth)-1]
	} else {
		start = 1
	}

	for i := 0; i <= len(yj.Index)-1; i++ {
		if i >= start {
			if i == start {
				cKey += yj.Index[i]
			} else if i > 1 {
				cKey += "." + yj.Index[i]
			}
		}

	}
	if yj.NextIndex != "" {
		if cKey == "" {
			cKey += yj.NextIndex
		} else {
			cKey += "." + yj.NextIndex
		}
	}
	return cKey
}

// Append string to jinja output
func (yj *Yaml2Jinja) Append(action string) {
	LogCtxt.WithFields(log.Fields{
		"function":  "Append",
		"Depth":     yj.Depth,
		"Index":     yj.Index,
		"NextIndex": yj.NextIndex,
		"ForDepth":  yj.ForDepth,
	}).Debug("Append Start")

	switch {
	case action == "ifStart":
		LogCtxt.WithFields(log.Fields{
			"function": "Append",
		}).Debug("Action is ifStart")
		yj.Output += "{% if " + goYamlFormat(yj.ConcatenateKey()) + " is defined %}" + "\n"
	case action == "forStart":
		LogCtxt.WithFields(log.Fields{
			"function": "Append",
		}).Debug("Action is forStart")
		yj.Output += "{% for " + yj.NextIndex + " in " + goYamlFormat(yj.ConcatenateKey()) + " %}" + "\n"
	case action == "ifEnd":
		LogCtxt.WithFields(log.Fields{
			"function": "Append",
		}).Debug("Action is ifEnd")
		yj.Output += "{% endif %}" + "\n"
	case action == "forEnd":
		LogCtxt.WithFields(log.Fields{
			"function": "Append",
		}).Debug("Action is forEnd")
		yj.Output += "{% endfor %}" + "\n"
	case action == "itemEmpty":
		LogCtxt.WithFields(log.Fields{
			"function": "Append",
		}).Debug("Action is itemEmpty")
		yj.Indent()
		yj.Output += yj.NextIndex + ": " + "\n"
	case action == "itemCtxt":
		LogCtxt.WithFields(log.Fields{
			"function": "Append",
		}).Debug("Action is itemCtxt")
		yj.Indent()
		yj.Output += yj.NextIndex + ": " + "{{ " + goYamlFormat(yj.ConcatenateKey()) + " }}" + "\n"
	default:
		fmt.Println("Unknown")
	}
	LogCtxt.WithFields(log.Fields{
		"function": "Append",
	}).Debug("Append End")
}

// goYamlFormat replace - with _
func goYamlFormat(key string) string {
	return strings.ReplaceAll(key, "-", "_")
}

// Convert transforms map[string]interface{} to go struct
func (yj *Yaml2Jinja) Convert(data []byte) (string, error) {
	// Unmarshal to map[string]interface{}
	var obj map[string]interface{}
	err := yaml.Unmarshal(data, &obj)
	if err != nil {
		return "", err
	}

	for k, v := range obj {
		yj.FileName = k + ".j2"
		yj.Jinjify(k, v, false)
	}

	return yj.Output, nil
}

// Jinjify transforms map key values to jinja structs
func (yj *Yaml2Jinja) Jinjify(k string, v interface{}, arrayElem bool) {

	//var key string
	if reflect.TypeOf(v) == nil || len(k) == 0 {
		// nil type
		LogCtxt.WithFields(log.Fields{
			"function": "Jinjify",
		}).Debug("nil type")
		return
	}

	switch reflect.TypeOf(v).Kind() {

	// If yaml object
	case reflect.Map:
		switch val := v.(type) {
		case map[interface{}]interface{}:
			// YAML object MAP
			LogCtxt.WithFields(log.Fields{
				"function": "Jinjify",
			}).Debug("YAML object MAP")
			yj.Depth++
			yj.Index = append(yj.Index, k)

			if !arrayElem {
				// Create new structure
				LogCtxt.WithFields(log.Fields{
					"function":  "Jinjify",
					"Depth":     yj.Depth,
					"Index":     yj.Index,
					"NextIndex": yj.NextIndex,
					"ForDepth":  yj.ForDepth,
				}).Debug("Not an arrayElem")
			}
			// If array of yaml objects
			for k1, v1 := range val {
				// YAML object MAP loop
				yj.NextIndex = k1.(string)

				LogCtxt.WithFields(log.Fields{
					"function":  "Jinjify",
					"Depth":     yj.Depth,
					"Index":     yj.Index,
					"NextIndex": yj.NextIndex,
					"ForDepth":  yj.ForDepth,
				}).Debug("YAML object MAP loop")

				if _, ok := k1.(string); ok {

					yj.Append("ifStart")
					if v1 == nil {
						yj.Append("itemCtxt")
					} else {
						yj.Append("itemEmpty")
					}
					yj.Jinjify(k1.(string), v1, false)
					yj.Append("ifEnd")
				}
			}
			if !arrayElem {
				// Not an arrayElem 2
				yj.Depth--
				yj.Index = yj.Index[:len(yj.Index)-1]

				LogCtxt.WithFields(log.Fields{
					"function":  "Jinjify",
					"Depth":     yj.Depth,
					"Index":     yj.Index,
					"NextIndex": yj.NextIndex,
					"ForDepth":  yj.ForDepth,
				}).Debug("Not an arrayElem 2")

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
			LogCtxt.WithFields(log.Fields{
				"function":  "Jinjify",
				"Depth":     yj.Depth,
				"Index":     yj.Index,
				"NextIndex": yj.NextIndex,
				"ForDepth":  yj.ForDepth,
			}).Debug("string, int, bool, float64")

		// if nested object
		case map[interface{}]interface{}:
			// Slice Object MAP
			LogCtxt.WithFields(log.Fields{
				"function":  "Jinjify",
				"Depth":     yj.Depth,
				"Index":     yj.Index,
				"NextIndex": yj.NextIndex,
				"ForDepth":  yj.ForDepth,
			}).Debug("Slice object MAP")

			key := k

			// start for loop
			yj.Append("forStart")
			yj.ForDepth = append(yj.ForDepth, yj.Depth)
			yj.ForIndent++
			for _, v1 := range val {
				// continue to loop through elements
				yj.Jinjify(key, v1, true)
			}
			// stop for loop
			yj.Append("forEnd")
			if len(yj.ForDepth) > 0 {
				yj.ForDepth = yj.ForDepth[:len(yj.ForDepth)-1]
			}
			yj.Depth--
			yj.ForIndent--
			yj.Index = yj.Index[:len(yj.Index)-1]

		}

	default:
		// Default
		LogCtxt.WithFields(log.Fields{
			"function":  "Jinjify",
			"Depth":     yj.Depth,
			"Index":     yj.Index,
			"NextIndex": yj.NextIndex,
			"ForDepth":  yj.ForDepth,
		}).Debug("Default")
	}
}

func init() {
	// Log as default ASCII formatter.
	Log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: false,
	})

	// Only log the info severity or above.
	//Log.SetLevel(log.InfoLevel)
	//Log.SetLevel(log.DebugLevel)
}
