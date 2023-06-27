package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/FrangipaneTeam/terraform-templates/internal/terraform"
	"github.com/FrangipaneTeam/terraform-templates/pkg/file"

	_ "embed"
)

func main() {
	fileName := flag.String("filename", "", "filename")
	flag.Parse()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if *fileName == "" {
		log.Fatal().Msg("filename is required")
	}

	if !file.IsFileExists(*fileName) {
		log.Fatal().Msgf("file %s not found", *fileName)
	}

	// Get the absolute path of the file
	absPath, err := filepath.Abs(*fileName)
	if err != nil {
		log.Fatal().Err(err).Msg("error getting absolute path")
	}

	// Determine the test Path with the absolute path
	d := filepath.Dir(absPath)
	testDir := filepath.Join(absPath, "../../../tests/", filepath.Base(d))

	// Determine the schema Path with the absolute path
	schemaDir := filepath.Join(absPath, "../")

	log.Info().Msgf("using schema dir %s", schemaDir)

	// test if filedir exists and is a directory
	dir, err := os.Stat(testDir)
	if err != nil || !dir.IsDir() {
		log.Fatal().Err(err).Msgf("testdir %s not found or not a directory", testDir)
	}

	log.Info().Msgf("using file %s", *fileName)

	f, err := file.ToString(*fileName)
	if err != nil {
		log.Fatal().Err(err).Msg("error reading file")
	}

	tfTypes := terraform.GetTFTypes(*fileName)
	if tfTypes == "" {
		log.Fatal().Msgf("tf types not found. The filename must be like `my_tf_name_datasource.go` or `my_tf_name_resource.go`")
	}
	log.Info().Msgf("tf type : %s", tfTypes)

	packageName := terraform.GetPackageName(f)
	if packageName == "" {
		log.Fatal().Msg("package name not found")
	}
	log.Info().Msgf("package name : %s", packageName)

	categoryName, resourceName := terraform.GetTFName(f)
	if categoryName == "" {
		log.Fatal().Msg("tfname not found. Please add a comment like `// tfname: category_resource_name")
	}
	log.Info().Msgf("categoryName : %s -- resourceName : %s", categoryName, resourceName)

	t := genTemplateConf(categoryName, resourceName, packageName, testDir, *fileName, schemaDir)
	errT := t.createTFFile(tfTypes)
	if errT != nil {
		log.Fatal().Err(errT).Msg("error creating file")
	}
}
