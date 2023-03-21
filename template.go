package main

import (
	"bytes"
	_ "embed"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
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
}

//go:embed templates/datasource.go.tmpl
var templateDatasource string

//go:embed templates/resource.go.tmpl
var templateResource string

//go:embed templates/acc_test_resource.go.tmpl
var templateAccTestResource string

//go:embed templates/acc_test_datasource.go.tmpl
var templateAccTestDataSource string

func genTemplateConf(categoryName, resourceName, packageName, testDir, fileName string) templateDef {
	t := templateDef{
		CategoryName:   categoryName,
		ResourceName:   resourceName,
		PackageName:    packageName,
		LowerCamelName: strcase.ToLowerCamel(resourceName),
		CamelName:      strcase.ToCamel(resourceName),
		Filename:       fileName,
		TestDir:        testDir,
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

	errWrite := os.WriteFile(t.Filename, tpl.Bytes(), 0o644)
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

	errWriteAccTest := os.WriteFile(t.TestDir+"/"+fileNameWithoutExtAndPath(t.Filename)+"_test.go", tplAccTest.Bytes(), 0o644)
	if errWriteAccTest != nil {
		return errWriteAccTest
	}

	return nil
}

func fileNameWithoutExtAndPath(fileName string) string {
	f := filepath.Base(fileName)
	return strings.TrimSuffix(f, filepath.Ext(f))
}
