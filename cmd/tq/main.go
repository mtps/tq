package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mtps/tq/toml"
	"github.com/mtps/tq/version"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Arguments struct {
	Ver        bool
	ToJson     bool
	ToToml     bool
	ScriptFile string
}

var (
	args Arguments
)

func init() {
	flag.BoolVar(&args.Ver, "v", false, "Print version information")
	flag.BoolVar(&args.ToJson, "json", false, "Convert input stream from toml to json")
	flag.BoolVar(&args.ToToml, "toml", false, "Convert input stream from json to toml")
	flag.StringVar(&args.ScriptFile, "f", "", "Command file to append to current editing commands")
}

func main() {
	flag.Parse()

	if args.Ver {
		fmt.Printf("%s\n", strings.Join(version.BuildInfo(), "\n"))
		return
	}

	var fn func(io.Reader) (string, error)
	if args.ToToml {
		fn = processJsonToToml
	} else if args.ToJson {
		fn = processTomlToJson
	} else {
		scripts, err := readScripts()
		if err != nil {
			panic(err)
		}
		fn = runScript(scripts)
	}

	if fn != nil {
		output, err := fn(os.Stdin)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s\n", output)
	}
}

func readScripts() ([]string, error) {
	var scripts []string
	for _, arg := range flag.Args() {
		scripts = append(scripts, strings.TrimSpace(arg))
	}

	// Append anything in the argScriptFile if provided
	if args.ScriptFile == "" {
		return scripts, nil
	}

	f, err := os.Open(args.ScriptFile)
	if err != nil {
		return []string{}, fmt.Errorf("failed to open script file %s: %w", args.ScriptFile, err)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		script := strings.TrimSpace(scanner.Text())
		if script != "" {
			scripts = append(scripts, script)
		}
	}
	return scripts, nil
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

func runScript(scripts []string) func (r io.Reader) (string, error) {
	return func(r io.Reader) (string, error) {
		tree, err := toml.LoadReader(r)
		if err != nil {
			return "", fmt.Errorf("faled to load toml file: %w", err)
		}

		var gets []string
		var sets []string
		for _, script := range scripts {
			parts := strings.Split(script, "=")
			if len(parts) == 1 {
				// Get operation
				gets = append(gets, parts[0])
			} else if len(parts) == 2 {
				// Set operation
				sets = append(sets, script)
			}
		}

		if len(gets) != 0 && len(sets) != 0 {
			return "", fmt.Errorf("set and get not allowed in same call")
		}

		if len(gets) > 0 {
			sb := strings.Builder{}
			for _, get := range gets {
				value := tree.Get(get)
				sb.WriteString(fmt.Sprintf("%s", value))
				sb.WriteString("\n")
			}
			return sb.String(), nil
		}

		for _, set := range sets {
			parts := strings.Split(set, "=")
			value := tree.Get(parts[0])

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
					return "", fmt.Errorf("expected int: %w", err)
				}
			case int64:
				i, err := strconv.Atoi(p)
				if err != nil {
					return "", fmt.Errorf("expected int64: %w", err)
				}
				v = int64(i)

			case float64:
				v, err = strconv.ParseFloat(p, 64)
				if err != nil {
					return "", fmt.Errorf("expected float64: %w", err)
				}

			case bool:
				v, err = strconv.ParseBool(p)
				if err != nil {
					return "", fmt.Errorf("expected bool: %w", err)
				}

			case string:
				v = parts[1]
			}

			tree.Set(parts[0], v)
		}

		return tree.ToTomlString()
	}
}

