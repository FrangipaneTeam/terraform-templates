package main

import (
	_ "embed"
	"flag"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	fileName := flag.String("filename", "", "filename")
	flag.Parse()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if *fileName == "" {
		log.Fatal().Msg("filename is required")
	}

	if !checkFileExists(*fileName) {
		log.Fatal().Msgf("file %s not found", *fileName)
	}

	log.Info().Msgf("using file %s", *fileName)

	f, err := getFile(*fileName)
	if err != nil {
		log.Fatal().Err(err).Msg("error reading file")
	}

	tfTypes := getTFTypes(*fileName)
	if tfTypes == "" {
		log.Fatal().Msgf("tf types not found. The filename must be like `my_tf_name_datasource.go` or `my_tf_name_resource.go`")
	}
	log.Info().Msgf("tf type : %s", tfTypes)

	packageName := getPackageName(f)
	if packageName == "" {
		log.Fatal().Msg("package name not found")
	}
	log.Info().Msgf("package name : %s", packageName)

	tfName := getTFName(f)
	if tfName == "" {
		log.Fatal().Msg("tfname not found. Please add a comment like `// tfname: my_tf_name")
	}
	log.Info().Msgf("tf name : %s", tfName)

	t := genTemplateConf(tfName, packageName, *fileName)
	errT := t.createTFFile(tfTypes)
	if errT != nil {
		log.Fatal().Err(errT).Msg("error creating file")
	}
}
