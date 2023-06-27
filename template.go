package main

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"

	_ "embed"
)

type templateDef struct {
	CategoryName   string
	ResourceName   string
	Name           string
	PackageName    string
	LowerCamelName string
	CamelName      string
	Filename       string
	TestDir        string
	SchemaDir      string
}

//go:embed templates/datasource.go.tmpl
var templateDatasource string

//go:embed templates/resource.go.tmpl
var templateResource string

//go:embed templates/acc_test_resource.go.tmpl
var templateAccTestResource string

//go:embed templates/acc_test_datasource.go.tmpl
var templateAccTestDataSource string

//go:embed templates/schema.go.tmpl
var templateSchema string

func genTemplateConf(categoryName, resourceName, packageName, testDir, fileName, schemaDir string) templateDef {
	t := templateDef{
		CategoryName:   categoryName,
		ResourceName:   resourceName,
		PackageName:    packageName,
		LowerCamelName: strcase.ToLowerCamel(resourceName),
		CamelName:      strcase.ToCamel(resourceName),
		Filename:       fileName,
		TestDir:        testDir,
		SchemaDir:      schemaDir,
	}

	if resourceName == "" {
		t.LowerCamelName = strcase.ToLowerCamel(categoryName)
		t.CamelName = strcase.ToCamel(categoryName)
	}

	return t
}

func (t templateDef) createTFFile(tfTypes string) error {
	templateS := templateDatasource
	templateAccTest := templateAccTestDataSource
	if tfTypes == "resource" {
		templateS = templateResource
		templateAccTest = templateAccTestResource
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

	errWrite := os.WriteFile(t.Filename, tpl.Bytes(), 0o600)
	if errWrite != nil {
		return errWrite
	}

	// for acc test
	tmplAccTest, errAccTest := template.New("template").Parse(templateAccTest)
	if errAccTest != nil {
		return errAccTest
	}

	var tplAccTest bytes.Buffer

	errAccTestExec := tmplAccTest.Execute(&tplAccTest, t)
	if errAccTestExec != nil {
		return errAccTestExec
	}

	errWriteAccTest := os.WriteFile(t.TestDir+"/"+fileNameWithoutExtAndPath(t.Filename)+"_test.go", tplAccTest.Bytes(), 0o600)
	if errWriteAccTest != nil {
		return errWriteAccTest
	}

	var tplSchema bytes.Buffer

	// if file not already exists create schema file
	if _, err := os.Stat(t.SchemaDir + "/" + (t.Filename) + "_schema.go"); os.IsNotExist(err) {
		tmplSchema, errSchemaTmpl := template.New("template").Parse(templateSchema)
		if errSchemaTmpl != nil {
			return errSchemaTmpl
		}

		errSchema := tmplSchema.Execute(&tplSchema, t)
		if errSchema != nil {
			return errSchema
		}

		errWriteSchema := os.WriteFile(t.SchemaDir+"/"+fileNameWithoutExtAndPath(t.Filename)+"_schema.go", tplSchema.Bytes(), 0o600)
		if errWriteSchema != nil {
			return errWriteSchema
		}
	}

	return nil
}

func fileNameWithoutExtAndPath(fileName string) string {
	f := filepath.Base(fileName)
	return strings.TrimSuffix(f, filepath.Ext(f))
}
