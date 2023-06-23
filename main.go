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

var dotPreamble = flag.String("preamble", "", "")
var graphName = flag.String("name", "G", "")

func main() {
	reader := os.Stdin
	writer := os.Stdout
	flag.Parse()

	var input map[string]interface{}
	err := json.NewDecoder(reader).Decode(&input)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Fprintf(writer, "digraph %s {\n", *graphName)
	fmt.Fprintf(writer, "node [shape=record];\n")
	fmt.Fprintln(writer, *dotPreamble)
	walk(writer, "", "N", input)
	fmt.Fprintf(writer, "}\n")
}
