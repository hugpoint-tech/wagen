# WAGEN

(WA)yland (GEN)erator

Language independent, humble and unpretentious wayland code generation.

- Are you tired of reading macro gobbledygook?
- Do your eyes bleed when you look at ugly generated code?
- Is generated code too complicated, opinionated or API does not fit your project?
- Have you tried "jump to definition" and couldn't because macros only expand at compile time?
- Do you want native wayland code for your favorite language?

Then this project is for you:

- SupExecuteTemplatesply your go templates and generate code for your language that looks the way you like and fits your needs.
- Put the generated code into your project and use as if you wrote it by hand.
- Use code completion and navigation tools that work for the rest of your code.
- Be happy, get stuff done, move on.


## Usage

The tool takes 2 directories:

- `in` with contains go tempates or subdirectories with go templates.
- `out` where rendered results will be saved. The file names and directories will mirror the structure of `in-dir`

XML Protocol definitions are embedded into the codebase. Update `protocols/` directory and rebuild wagen if
you need to update 

```
Usage: wagen -in <TEMPLATE_DIR> -out <OUTPUT_DIR>
       wagen -list
       wagen -help
       wagen -version
       
Options:
  -in <TEMPLATE_DIR>
  -out <OUTPUT_DIR>
  -p <PROTOCOL>                 Name of the protocol to pass into the template.
                                Can be used multiple times. 
                                Special names are recognized:
                                - core - core wayland protocol (default)
                                - stable - all stable protocols
                                - staging - all staging protocols
                                - unstable - all unstable protocols
  -list                         List available protocols
  -help                         Print this help
  -version                      Print version
```

## Build / Install

Wagen is written in golang using only standard library. Building it is as simple as `go build`


## Templates

Each template receives a list of protocols as input. See [types.go](types.go) for struct definitions/


See [railway templates](./templates/railway) for examples for now.


On top of standard go template functionality the following functions are also available for use in templates:

- `ToLines`  - takes a string and returns a list of strings separated by newline. (effectively `strings.Split(s, "\n")`)
- `ToPascal` - takes any amount of strings, and turns them into a single PascalCase string. 
               It can also convert from snake_case.
- `Trim` - `strings.TrimSpace` trims whitespaces on both ends of the string.