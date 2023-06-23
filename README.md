[![Go](https://github.com/thebirk/json2dot/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/thebirk/json2dot/actions/workflows/go.yml)

---

# json2dot - A simple json to dot conversion tool

## How does it work

Objects are turned into nodes. Where the node label contains all keys and values.

Every field with a nested object turns into a edge to that new object.

All JSON types except for arrays are supported.

## Future changes

Move from `node [shape=record]` to [HTML-like labels](https://graphviz.org/doc/info/shapes.html#html).

Support array type.

## Installation

```bash
go get github.com/thebirk/json2dot
go install github.com/thebirk/json2dot
```

> Make sure `$GOPATH/bin` is in your `$PATH`

## Example

```
$ json2dot
{"a": 123, "b": false, "c": {"d": "another node"}}
digraph G {
node [shape=record];
_N_c [label="{ d : \"another node\" }"]
_N -> _N_c [label="c"]
_N [label="{ a : 123.000000|b : false }"]
}
```

![](/img/example1.png)

### Parameters

We also support using only parameters.

```
$ json2dot -in tree.json -out tree.dot
```

With support for a positional input parameter as well.

```
$ json2dot tree.json | dot -Tpng -otree.png
```

Using only pipes.

```
$ cat tree.json | json2dot | dot -Tpng -otree.png
```

All parameters

```
Usage: json2dot [options] [INPUT]

Parameters:
  INPUT (optional) - Cannot be combined with '-in'
        Input file path

Options:
  -in value
        Input file path
  -name string
        Graph name (default "G")
  -out value
        Output file path
  -preamble string
        Inject text into the output
```