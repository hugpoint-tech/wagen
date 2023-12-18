package main

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

//func compose(parts ...string) *template.Template {
//	files := make([]string, 0, 5)
//	for _, p := range parts {
//		files = append(files, p+".gohtml")
//	}
//	page, err := template.ParseFS(fs, files...)
//	AssertOk("template composition failed", err)
//	return page
//}

func ToLines(s string) []string {
	return strings.Split(s, "\n")
}

func ToPascal(chunks ...string) string {
	concatenated := strings.Join(chunks, "_")
	concatenated = strings.ReplaceAll(concatenated, ".", "_")

	words := strings.Split(concatenated, "_")

	for i, word := range words {
		words[i] = strings.Title(word)
	}
	return strings.Join(words, "")
}

var TemplateFuncs = template.FuncMap{
	"ToLines":  ToLines,
	"ToPascal": ToPascal,
	"Trim":     strings.TrimSpace,
	"Contains": strings.Contains,
	"Upper":    strings.ToUpper,
}

func ReadTemplates(dir string) map[string]*template.Template {
	files := []string{}
	templates := make(map[string]*template.Template)

	dir = filepath.Clean(dir)
	if !strings.HasSuffix(dir, string(filepath.Separator)) {
		dir += string(filepath.Separator)
	}

	err := filepath.WalkDir(dir, func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(s, ".tmpl") {

			files = append(files, s)
		}
		return nil
	})
	AssertOk("failed to walk templates directory", err)

	for _, fname := range files {

		file, err := os.Open(fname)
		if err != nil {
			fmt.Println("error opening file:", err)
			continue
		}
		defer file.Close()

		name := path.Base(fname)
		tmpl, err := template.New(name).Funcs(TemplateFuncs).ParseFiles(fname)
		if err != nil {
			fmt.Println("error parsing template", fname, err)
			continue
		}

		templates[strings.TrimPrefix(fname, dir)] = tmpl
	}

	return templates
}

func RenderTemplates(templates map[string]*template.Template, protocols []Protocol, out string) []error {
	var errors []error

	for tname, tmpl := range templates {
		indir := filepath.Dir(tname)
		outdir := out + "/" + indir
		inname := filepath.Base(tname)

		err := os.MkdirAll(outdir, 0755)
		if err != nil {
			errors = append(errors, fmt.Errorf("error creating output file directory: %s", err))
			continue
		}

		outname := strings.TrimSuffix(outdir+"/"+inname, ".tmpl")
		outfile, err := os.Create(outname)
		if err != nil {
			errors = append(errors, fmt.Errorf("error creating output file %s: %s", outname, err))
			continue
		}
		defer outfile.Close()

		err = tmpl.Execute(outfile, protocols)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to execute template: %s", err))
			continue
		}
		fmt.Println("OK:", tname, "->", outfile.Name())
	}
	return errors
}
