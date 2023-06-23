package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/iancoleman/orderedmap"
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

func generateGraph(writer io.Writer, prefix string, name string, input orderedmap.OrderedMap) {
	var fields []string
	qualified_name := fmt.Sprintf("%s_%s", prefix, name)

	for _, key := range input.Keys() {
		jsonVal, _ := input.Get(key) // lets assume the key didn't disappaer
		switch val := jsonVal.(type) {
		case bool:
			b := "false"
			if val {
				b = "true"
			}
			fields = append(fields, fmt.Sprintf("%s : %s", key, b))
		case float64:
			fields = append(fields, fmt.Sprintf("%s : %g", key, val))
		case string:
			fields = append(fields, fmt.Sprintf("%s : \\\"%s\\\"", key, val))
		case nil:
			fields = append(fields, fmt.Sprintf("%s : null", key))
		case orderedmap.OrderedMap:
			generateGraph(writer, qualified_name, key, val)
			fmt.Fprintf(writer, "%s -> %s [label=\"%s\"]\n", qualified_name, qualified_name+"_"+key, key)
		case []interface{}:
			fields = append(fields, fmt.Sprintf("%s : array", key))
			fmt.Fprintf(os.Stderr, "Skipping key '%s'. Arrays are currently not supported. \n", key)
		}
	}

	fmt.Fprintf(writer, "%s [label=\"{ %s }\"]\n", qualified_name, strings.Join(fields, "|"))
}

var dotPreamble = flag.String("preamble", "", "Inject text into the output")
var graphName = flag.String("name", "G", "Graph name")
var outPathFlag optionalStringFlag
var inPathFlag optionalStringFlag

func printUsage() {
	fmt.Printf("Usage: %s [options] [INPUT]\n", os.Args[0])
	fmt.Println()
	fmt.Println("Parameters:")
	fmt.Println("  INPUT (optional) - Cannot be combined with '-in'")
	fmt.Println("    \tInput file path")
	fmt.Println()
	fmt.Println("Options:")
	flag.PrintDefaults()
}

func main() {
	reader := os.Stdin
	writer := os.Stdout

	flag.Usage = printUsage
	flag.Var(&outPathFlag, "out", "Output file path")
	flag.Var(&inPathFlag, "in", "Input file path")
	flag.Parse()

	inPath := ""
	if flag.NArg() > 0 {
		if flag.NArg() != 1 {
			fmt.Fprintln(os.Stderr, "Too many positional arguments")
			flag.Usage()
			os.Exit(1)
		}

		if inPathFlag.set {
			fmt.Fprintln(os.Stderr, "Cannot use both positional argument in and '-in'")
			flag.Usage()
			os.Exit(2)
		}

		inPath = flag.Arg(0)
	}

	if inPathFlag.set {
		inPath = inPathFlag.value
	}

	if inPath != "" {
		var err error
		reader, err = os.Open(inPath)
		if err != nil {
			panic(err)
		}
	}

	var input orderedmap.OrderedMap
	err := json.NewDecoder(reader).Decode(&input)
	if err != nil {
		log.Fatalln(err)
	}

	if outPathFlag.set {
		var err error
		writer, err = os.OpenFile(outPathFlag.value, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
	}

	fmt.Fprintf(writer, "digraph %s {\n", *graphName)
	fmt.Fprintf(writer, "node [shape=record];\n")
	if *dotPreamble != "" {
		fmt.Fprintln(writer, *dotPreamble)
	}
	generateGraph(writer, "", "N", input)
	fmt.Fprintf(writer, "}\n")
}
