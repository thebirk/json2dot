package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type optionalStringFlag struct {
	set   bool
	value string
}

func (flag *optionalStringFlag) Set(str string) error {
	flag.set = true
	flag.value = str

	return nil
}

func (flag *optionalStringFlag) String() string {
	return flag.value
}

func walk(writer io.Writer, prefix string, name string, input map[string]interface{}) {
	var fields []string
	qualified_name := fmt.Sprintf("%s_%s", prefix, name)

	for key, jsonVal := range input {
		switch val := jsonVal.(type) {
		case bool:
			b := "false"
			if val {
				b = "true"
			}
			fields = append(fields, fmt.Sprintf("%s : %s", key, b))
		case float64:
			fields = append(fields, fmt.Sprintf("%s : %f", key, val))
		case string:
			fields = append(fields, fmt.Sprintf("%s : \\\"%s\\\"", key, val))
		case nil:
			fields = append(fields, fmt.Sprintf("%s : null", key))
		case map[string]interface{}:
			walk(writer, qualified_name, key, val)
			fmt.Fprintf(writer, "%s -> %s [label=\"%s\"]\n", qualified_name, qualified_name+"_"+key, key)
		case []interface{}:
			log.Println("arrays not supported")
		}
	}

	fmt.Fprintf(writer, "%s [label=\"{ %s }\"]\n", qualified_name, strings.Join(fields, "|"))
}

var dotPreamble = flag.String("preamble", "", "Inject text into the output")
var graphName = flag.String("name", "G", "Graph name")
var outPath optionalStringFlag

func main() {
	reader := os.Stdin
	writer := os.Stdout

	flag.Var(&outPath, "out", "Output file path")
	flag.Parse()

	inPath := ""
	if flag.NArg() > 0 {
		if flag.NArg() != 1 {
			fmt.Fprintln(os.Stderr, "Too many positional arguments")
			flag.Usage()
			os.Exit(1)
		}

		inPath = flag.Arg(0)
	}

	if inPath != "" {
		var err error
		reader, err = os.Open(inPath)
		if err != nil {
			panic(err)
		}
	}

	var input map[string]interface{}
	err := json.NewDecoder(reader).Decode(&input)
	if err != nil {
		log.Fatalln(err)
	}

	if outPath.set {
		var err error
		writer, err = os.OpenFile(outPath.value, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
	}

	fmt.Fprintf(writer, "digraph %s {\n", *graphName)
	fmt.Fprintf(writer, "node [shape=record];\n")
	if *dotPreamble != "" {
		fmt.Fprintln(writer, *dotPreamble)
	}
	walk(writer, "", "N", input)
	fmt.Fprintf(writer, "}\n")
}
