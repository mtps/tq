package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mtps/tq/toml"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
)


var (
	argToJson bool
	argToToml bool
)

func init() {
	flag.BoolVar(&argToJson, "json", false, "Convert input stream from toml to json")
	flag.BoolVar(&argToToml, "toml", false, "Convert input stream from json to toml")
}

func main() {
	flag.Parse()

	var fn func(io.Reader) (string, error)
	if argToToml {
		fn = processJsonToToml
	} else if argToJson {
		fn = processTomlToJson
	} else {
		fn = runScript
	}

	if fn != nil {
		output, err := fn(os.Stdin)
		if err != nil {
			panic(err)
		}

		fmt.Printf("~~~~~~~~~~~~~~~~~~~~~~\n%s\n~~~~~~~~~~~~~~~~~~~~~~\n", output)
	}
}

func processJsonToToml(r io.Reader) (string, error) {
	bz, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	var m map[string]interface{}
	err = json.Unmarshal(bz, &m)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal json: %w", err)
	}
	tree, err := toml.TreeFromMap(m)
	if err != nil {
		return "", fmt.Errorf("faled to convert map to toml: %w", err)
	}

	return tree.String(), nil
}

func processTomlToJson(r io.Reader) (string, error) {
	tree, err := toml.LoadReader(r)
	if err != nil {
		return "", fmt.Errorf("faled to load toml file: %w", err)
	}
	m := tree.ToMap()
	bz, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("failed to marshal toml to json: %w", err)
	}
	return string(bz), nil
}

func runScript(r io.Reader) (string, error) {
	var scripts []string
	for _, arg := range flag.Args() {
		scripts = append(scripts, strings.TrimSpace(arg))
	}

	tree, err := toml.LoadReader(r)
	if err != nil {
		return "", fmt.Errorf("faled to load toml file: %w", err)
	}

	for _, script := range scripts {
		fmt.Printf("script: %v\n", script)
		parts := strings.Split(script, "=")

		// Get operation
		if len(parts) == 1 {
			part := tree.Get(parts[0])
			return fmt.Sprintf("%s\n", part), nil
		} else if len(parts) == 2 {
			value := tree.Get(parts[0])

			fmt.Printf("value:%v\n", reflect.TypeOf(value))
			var v interface{}
			// Only fields of the following types are supported:
			//   * string
			//   * bool
			//   * int
			//   * int64
			//   * float64
			p := strings.TrimSpace(parts[1])
			switch value.(type) {
			case int:
				v, err = strconv.Atoi(p)
				if err != nil {
					return "", err
				}

			case int64:
				i, err := strconv.Atoi(p)
				if err != nil {
					return "", err
				}
				v = int64(i)

			case float64:
				v, err = strconv.ParseFloat(p, 64)
				if err != nil {
					return "", err
				}

			case bool:
				v, err = strconv.ParseBool(p)
				if err != nil {
					return "", err
				}

			case string:
				v = parts[1]
			}
			fmt.Printf("v:%v\n", v)
			tree.Set(parts[0], v)
			return tree.ToTomlString()
		}
	}

	return "", nil
}

