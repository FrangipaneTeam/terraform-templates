package main

import (
	"bytes"
	_ "embed"
	"html/template"
	"os"

	"github.com/iancoleman/strcase"
)

type templateDef struct {
	Name           string
	PackageName    string
	LowerCamelName string
	CamelName      string
	Filename       string
}

//go:embed templates/datasource.go.tmpl
var templateDatasource string

//go:embed templates/resource.go.tmpl
var templateResource string

func genTemplateConf(tfName, packageName, fileName string) templateDef {
	t := templateDef{
		Name:           tfName,
		PackageName:    packageName,
		LowerCamelName: strcase.ToLowerCamel(tfName),
		CamelName:      strcase.ToCamel(tfName),
		Filename:       fileName,
	}
	return t
}

func (t templateDef) createTFFile(tfTypes string) error {
	templateS := templateDatasource
	if tfTypes == "resource" {
		templateS = templateResource
	}

	tmpl, err := template.New("template").Parse(templateS)
	if err != nil {
		return err
	}

	var tpl bytes.Buffer

	errExec := tmpl.Execute(&tpl, t)

	if errExec != nil {
		return errExec
	}

	errWrite := os.WriteFile(t.Filename, tpl.Bytes(), 0o644)
	if errWrite != nil {
		return errWrite
	}

	return nil
}
