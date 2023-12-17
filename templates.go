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
}

func ReadTemplates() map[string]*template.Template {
	templateFiles := []string{}

	err := filepath.WalkDir("templates", func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(s, ".tmpl") {
			templateFiles = append(templateFiles, s)
		}
		return nil
	})
	AssertOk("failed to walk templates directory", err)

	parsedTemplates := make(map[string]*template.Template)

	for _, file := range templateFiles {
		// Open the template file
		tmplFile, err := os.Open(file)
		if err != nil {
			fmt.Println("Error opening file:", err)
			continue
		}
		defer tmplFile.Close()

		name := path.Base(file)
		tmpl, err := template.New(name).Funcs(TemplateFuncs).ParseFiles(file)
		if err != nil {
			fmt.Println("Error parsing template:", err)
			continue
		}

		parsedTemplates[file] = tmpl
	}

	return parsedTemplates
}

func ExecuteTemplates(templates map[string]*template.Template, protocols []Protocol) {
	for name, tmpl := range templates {
		parentDir := filepath.Dir(name)
		parentDir = strings.TrimPrefix(parentDir, "templates/")
		renderDir := "rendered/" + parentDir
		fileName := filepath.Base(name)

		err := os.MkdirAll(renderDir, 0755)
		if err != nil {
			fmt.Println("Error creating output file directory:", err)
			continue
		}

		outputFile, err := os.Create(strings.TrimSuffix(renderDir+"/"+fileName, ".tmpl"))
		if err != nil {
			fmt.Println("Error creating output file:", err)
			continue
		}
		defer outputFile.Close()

		err = tmpl.Execute(outputFile, protocols)
		if err != nil {
			fmt.Println("Error executing template:", err)
			continue
		}
		fmt.Println("Template", name, "executed and saved to", outputFile.Name())
	}
}
