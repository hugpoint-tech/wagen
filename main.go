package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type arrayFlag []string

func (af *arrayFlag) String() string {
	return strings.Join(*af, ", ")
}

func (af *arrayFlag) Set(value string) error {
	*af = append(*af, value)
	return nil
}

func AssertOk(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	}
}

func main() {

	var in string
	var out string
	var help bool
	var version bool
	var list bool
	var show bool
	var names arrayFlag
	var protocols []Protocol

	flag.StringVar(&in, "in", "", "Input directory containing Go templates")
	flag.StringVar(&out, "out", "", "Output directory for rendered results")
	flag.BoolVar(&list, "list", false, "List available protocols")
	flag.BoolVar(&show, "show", false, "Show protocol definitions")
	flag.BoolVar(&help, "help", false, "Print this help")
	flag.BoolVar(&version, "version", false, "Print version")
	flag.Var(&names, "p", "protocol name, can be used multiple times")
	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	if version {
		fmt.Println("wagen version 0.1.0")
		os.Exit(0)
	}

	if list {
		fmt.Println("Available protocols:")
		for k := range PROTOCOLS {

			fmt.Printf(" - %s\n", k)
		}
		os.Exit(0)
	}

	if len(names) == 0 {
		protocols = append(protocols, PROTOCOLS["wayland"])
	}

	for _, name := range names {
		if name == "core" {
			protocols = append(protocols, PROTOCOLS["wayland"])
			continue
		}

		if name == "stable" {
			for _, p := range PROTOCOLS {
				if p.Maturity == Stable {
					protocols = append(protocols, p)
				}
			}
			continue
		}

		if name == "staging" {
			for _, p := range PROTOCOLS {
				if p.Maturity == Staging {
					protocols = append(protocols, p)
				}
			}
			continue
		}

		if name == "unstable" {
			for _, p := range PROTOCOLS {
				if p.Maturity == Unstable {
					protocols = append(protocols, p)
				}
			}
			continue
		}

		p, ok := PROTOCOLS[name]
		if !ok {
			fmt.Println("Unknown protocol name:", name)
			fmt.Println("Run with -list to see available protocols")
			os.Exit(1)
		}
		protocols = append(protocols, p)
	}

	if in == "" || out == "" {
		fmt.Println("Error: Both --in-dir and --out-dir are required")
		flag.Usage()
		os.Exit(1)
	}

	templates := ReadTemplates(in)
	if len(templates) == 0 {
		fmt.Println("No templates were found in ", in)
		os.Exit(1)
	}

	errors := RenderTemplates(templates, protocols, out)
	if len(errors) != 0 {
		fmt.Println("There were errors during template rendering")
		for _, e := range errors {
			fmt.Println(e)
		}
		os.Exit(1)
	}
}
