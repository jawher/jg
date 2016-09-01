# JG

A CLI tool to generate JSON from the command-line.

## Motivation

In my day to day work, I find myself using the excellent [httpie](https://github.com/jkbrzt/httpie) tool to test REST apis
I'm working on.

Usually, those APIs accept and produce JSON payloads.
httpie comes with handy support for generating the JSON body of outgoing requests.

However, this support is pretty basic and suffers from several limitations:

* The generated JSON is only one level deep (jkbrzt/httpie#78)
* Bare bone support for arrays

Out of frustration with these issues, I ended up writing `jg` to solve them.

## At a glance

`jg` is run from the command-line with one or multiple generator expressions describing the desired JSON:
  
```
jg 'customer = { id=:33 name=jawher address.zip=123 links.repos=[:1 :2 :3]}'
```

This produces the following JSON

```json
{
       	"customer": {
       		"id": 33,
       		"name": "jawher",
       		"address": {
       			"zip": "123"
       		},
       		"links": {
       			"repos": [1, 2, 3]
       		}
       	}
}
```

Which can then be fed to other tools (httpie for example).

## Installation

Head to [the releases page](https://github.com/jawher/jg/releases/latest) and grab the latest build for your platform.
For the time being, you'll have to rename the downloaded file from `jg-${platform}` to just `jg`.
Move it somewhere in your `PATH` and ensure that it is executable.

An alternative way if you have go installed is to run `go install github.com/jawher/jg`

## Generators

The generator expressions syntax resembles the JSON syntax with a few caveats:

* It uses the `=` sign instead of `:` to separate a field and its value
* Strings are not quoted (unless they contain spaces)
* Since values are treated as strings by default, other types (numbers, booleans, null) need to be prefixed with a `:`
* No commas required to separate elements of an object or array

Think of it as the lightweight JSON syntax you dreamt of !

### Field

SYNTAX: `field = value`

Adds a field named `field` to the current object.

The value can be one of:
* Literal: will be treated a string by default, unless prefixed with `:` (for numbers, booleans and `null`)
* Variable: prefixed with the `$` sign, e.g. ``$name`. Will be replaced by its value (a variable can be set using the `-s` flag)
* Object generator: see below for details
* Array generator: see below for details

EXAMPLES: 

* `name = foo`: adds a field to a JSON object named `foo` containing the string value `"foo"`
* `age = :42` : prefixing a value with `:` instructs `jg` to handle it as a *raw* value, i.e. as a number, a boolean or as `null` and not treat it like a string
* `id = $x` : prefixing a value with `$` instructs `jg` to handle it as a variable.

### Dotted field

SYNTAX: `parent.child = value`

Same as field, but also creates the intermediary objects.

If you need the field name to contain a dot, enclose it in quotes, e.g. `"field.with.dot"=value`

EXAMPLES:

`customer.address.zip = 123`: adds a field named `customer` containing an object with an `address field`, itself an object with one `zip` field.
 

### Object

SYNTAX `field = { GENERATORS... }`

Creates a JSON object.
To set the object fields, specify the required generator expressions inside the brackets.
The accepted generators are field generators.

EXAMPLES:

`customer = { name=foo, age=:30 }`


### Array

SYNTAX `field = [ GENERATORS... ]`

Creates a JSON array.

To set the array elements, specify the required generator expressions inside the square brackets.
The accepted generators are:

* literal: will be added as is to the array
* variable: a name prefixed with the `$` sign. Its value will be added as is to the array
* field generator: will add a JSON object containing the generated field as an element
* object generator: will add a the generated JSON object as an element
* array generator: will add a the generated JSON array as an element

EXAMPLES:

* `tags = [ foo bar qix $x ]`: creates an array containing the 3 strings `"foo"`, `"bar"` and `"qix"` and the value of the `x` variable.
* `ids = [ :1 :2 :3 ]`: prefix the literals with `:` to handle them as *raw*, i.e. number values
* `results = [ id=:1 id=:2]`: creates an array containing 2 objects, each with one field named `id`
* `results = [ {id=:1 name=foo} {id=:2 name=bar}]`: creates an array containing 2 objects, each with 2 fields `id` and `name`

## CLI usage

### Pretty-print

By default, `jg` outputs unindented JSON in a single line.
The `-p` can be used to instruct `jg` to pretty-print the generated JSON.

### Substitutions

The `-s` flag can be used to declare one or more substitutions:

`jg -s x=World '[ hello $x]'`

Whenever `$x` is encountered in the generator expressions, it will be replaced with the string `"World"`.

A substitution value can be one of:

* Literal: The value will be stored as a string
* Raw literal: Using `:` as prefix (for numbers, booleans and `null`)

### One object mode
Calling `jg` with one or multiple field generators (`field = value`) results in producing a single JSON object

```
jg foo=a bar=c
```

produces a single JSON object with 2 fields:
 
```json
{"bar":"c","foo":"a"}
```

### Multiple objects/arrays
`jg` can also be made to generate one or multiple JSON objects or arrays by using the object and array delimiters (`{}` and `[]`)
at the top level:

Objects:
```
jg {foo=a} {bar=c}
```

```json
{"foo":"a"}
{"bar":"c"}
```

Arrays (The generator expression is singled-quoted to avoid the shell trying to interpret it):
```
jg '[foo=a] [bar=c]'
```

```json
[{"foo":"a"}]
[{"bar":"c"}]
```

## License

This work is published under the MIT license.

Please see the `LICENSE` file for details.