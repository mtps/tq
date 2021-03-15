package main

import (
	toml "github.com/mtps/tq/toml"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
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

	if argToJson == false && argToToml == false {
		flag.Usage()
		os.Exit(1)
	}

	var fn func(io.Reader) (string, error)
	if argToToml {
		fn = processJsonToToml
	} else if argToJson {
		fn = processTomlToJson
	}

	if fn != nil {
		output, err := fn(os.Stdin)
		if err != nil {
			panic(err)
		}

		fmt.Printf(output)
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



