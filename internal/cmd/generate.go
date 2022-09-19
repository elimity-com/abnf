package cmd

import (
	"bytes"
	"github.com/elimity-com/abnf"
	"github.com/elimity-com/abnf/internal/config"
	"io"
	"log"
	"os"
	"path/filepath"
)

func Generate(dir, filename string, stderr io.Writer) error {
	configPath, configuration, err := config.ReadConfig(stderr, dir, filename)
	if err != nil {
		return err
	}

	if configuration.VerboseFlag {
		log.Printf("Config Path: %v, Spec: %v, Package: %v, Output: %v", configPath, configuration.SpecFile, configuration.PackageName, configuration.OutputPath)
	}

	rawABNF, err := os.ReadFile(configuration.SpecFile)
	if err != nil {
		log.Fatalf("Error %v", err)
	}

	g := abnf.CodeGenerator{
		PackageName: configuration.PackageName,
		RawABNF:     rawABNF,
	}
	buf := new(bytes.Buffer)

	if configuration.Generation == "alternatives" {
		g.GenerateABNFAsAlternatives(buf)
	} else {
		g.GenerateABNFAsOperators(buf)
	}

	assembledOutputPath := filepath.Join(configuration.OutputPath, configuration.PackageName)

	if configuration.VerboseFlag {
		log.Printf("Assembled output path: %v", assembledOutputPath)
	}

	err = os.MkdirAll(assembledOutputPath, 0775)
	if err != nil {
		log.Fatalf("Unable to make assembled output path directory structure: %v", err)
	}

	assembledFilePath := filepath.Join(assembledOutputPath, configuration.GoFileName)
	if configuration.VerboseFlag {
		log.Printf("Assembled Go file path: %v", assembledFilePath)
	}

	err = os.WriteFile(assembledFilePath, buf.Bytes(), 0644)
	if err != nil {
		log.Fatalf("Error writing Go source code to file: %v", err)
	}

	if configuration.VerboseFlag {
		log.Printf("%v ABNF specification successfully processed", configuration.SpecFile)
		log.Printf("Output Directory: %v", configuration.OutputPath)
		log.Printf("Go Code %v/%v", configuration.PackageName, configuration.GoFileName)
	}
	return nil
}
