package main

import (
	"fmt"

	"strings"

	"os"

	"bytes"

	"encoding/json"

	"github.com/jawher/mow.cli"
)

var (
	version   string
	gitCommit string
	buildDate string
)

func main() {
	app := cli.App("jg", "a CLI to generate JSON")
	app.LongDesc = HELP
	app.Version("v version", fmt.Sprintf("%s [sha: %s, time: %s]", version, gitCommit, buildDate))

	app.Spec = "[-p] [-s...] GENERATORS..."

	prettyPrint := app.Bool(cli.BoolOpt{
		Name: "p pretty-print",
		Desc: "Pretty print the generated JSON",
	})

	substitutions := Substitutions{}
	app.Var(cli.VarOpt{
		Name:  "s substitution",
		Desc:  "Substition in the format name=value where value can be a literal (string) or a raw literal (prefixed by ':', e.g. :true, :42, :null)",
		Value: &substitutions,
	})

	generators := app.Strings(cli.StringsArg{
		Name: "GENERATORS",
		Desc: "Generator expressions",
	})

	app.Action = func() {
		spec := strings.Join(*generators, " ")
		parsedGenerators, err := ParseGenerators(spec)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			cli.Exit(1)
		}

		for _, g := range parsedGenerators {
			output, _ := json.Marshal(g.Gen(substitutions))
			if *prettyPrint {
				var indentedOutput bytes.Buffer
				json.Indent(&indentedOutput, output, "", "\t")
				output = indentedOutput.Bytes()
			}
			fmt.Printf("%s\n", output)
		}
	}

	app.Run(os.Args)
}

type Substitutions map[string]Any

func (s *Substitutions) Set(v string) error {
	parts := strings.SplitN(v, "=", 2)

	if len(parts) != 2 {
		return fmt.Errorf("Invalid substitution %s: should be in the form name=value", v)
	}

	varName := parts[0]
	varValue, err := ParseSubstValue(parts[1])
	if err != nil {
		return err
	}
	(*s)[varName] = varValue
	return nil
}

func (s *Substitutions) String() string {
	res := ""
	sep := ""
	for name, value := range *s {
		res += fmt.Sprintf("%s%s=%q", sep, name, value)
		sep = ", "
	}
	return res
}

func (s *Substitutions) IsMulti() bool {
	return true
}

func (s *Substitutions) Clear() {
	*s = map[string]Any{}
}

const HELP = `jg - a CLI tool to generate JSON

EXAMPLES:

	$ jg name=foo
	{"name":"foo"}

	$ jg id=:32
	{"id":32}

	$ jg name=foo id=:32
	{"id":32,"name":"foo"}

	$ jg foo.bar=hello
	{"foo":{"bar":"hello"}}

	$ jg parent.child1=v1 parent.child2=v2
	{"parent":{"child1":"v1","child2":"v2"}}

	$ jg 'parent={child1=v1 child2=v2}'
	{"parent":{"child1":"v1","child2":"v2"}}

	$ jg '[foo bar]'
	["foo","bar"]

	$ jg '[:1 :2 :3]'
	[1,2,3]

	$ jg '[ foo=:true bar=:false]'
	[{"foo":true},{"bar":false}]

	$ jg '[ foo.yes=:true bar.yes=:false]'
	[{"foo":{"yes":true}},{"bar":{"yes":false}}]

	$ jg '[ {foo=a yes=:true} {bar=b yes=:false} ]'
	[{"foo":"a","yes":true},{"bar":"b","yes":false}]


GENERATOR SYNTAX:

	* Field

		field = value

	If a field names contains dots, jg will treat it as a path and will
	create the intermediary objects.
	Enclose it in double quotes to disable that.

	Values can be a literal, an object generator or an array generator.
	Literals are treated as strings by default.
	Prefix the literal with ':' for it to be treated as a number or a boolean.
	Multi word literals must be enclosed in double quotes.

	* Object

		field = { GENERATOR... }

	Creates a JSON object.
	To set the object fields, specify the required generator expressions inside the brackets.
	The accepted generators are field generators.

	* Array

		field = [ GENERATOR... ]

	Creates a JSON array.
	To set the array elements, specify the required generator expressions inside the square brackets.
	The accepted generators are:

		- literal: will be added as is to the array
		- field generator: will add a JSON object containing the generated field as an element
		- object generator: will add a the generated JSON object as an element
		- array generator: will add a the generated JSON array as an element

`
