package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	cfg "github.com/Hunga1/go-openapi-transform/pkg/Configuration"
	"github.com/neotoolkit/openapi"
)

const (
	OAS_DOCS_DIR = "./docs"
)

func main() {
	var config cfg.Configuration
	var paths []string
	var totalTransformed int = 0

	// Get Configuration
	config = *cfg.NewConfig()

	// Get OAS Doc files
	filenames := getOASDocFilenames()
	if len(filenames) == 0 {
		fmt.Println("No OAS doc files found.")
		os.Exit(0)
	} else {
		fmt.Println("OAS docs to transform:")
		for _, f := range filenames {
			relPath := OAS_DOCS_DIR + "/" + f

			if len(config.Whitelist) != 0 {
				// Include whitelisted OAS files
				if isWhitelistedFile(f, config.Whitelist) {
					fmt.Println(relPath)
					paths = append(paths, relPath)
				}
			} else {
				// Ignore blacklisted OAS files
				if isIgnoredFile(f, config.IgnoreFiles) {
					fmt.Printf("%v (Ignored)\n", relPath)
				} else {
					fmt.Println(relPath)
					paths = append(paths, relPath)
				}
			}
		}
	}

	fmt.Printf("Total OAS doc files to transform: %v\n", len(paths))

	for _, fp := range paths {
		_ = transformOASDoc(fp)
		fmt.Println("Successfully transformed OAS doc file: ", fp)
		totalTransformed = totalTransformed + 1
	}

	fmt.Printf("Successfully transformed %v OAS doc files!\n", totalTransformed)
}

func getOASDocFilenames() []string {
	var filenames []string

	files, err := ioutil.ReadDir(OAS_DOCS_DIR)
	if err != nil {
		log.Fatal("Failed to read OAS docs directory. Error: ", err)
	}

	for _, f := range files {
		filenames = append(filenames, f.Name())
	}

	return filenames
}

func transformOASDoc(path string) openapi.OpenAPI {
	fmt.Printf("Parsing OAS doc file '%v'...\n", path)

	oapi, err := openapi.Parse(getOASDocFilesContent(path))
	if err != nil {
		log.Fatalf("Failed to parse OAS doc file: %v\nError: %v\n", path, err)
	}

	return oapi
}

func getOASDocFilesContent(path string) []byte {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read OAS doc file: %v\nError: %v\n", path, err)
	}
	return content
}

func isIgnoredFile(filename string, ignoredFiles []string) bool {
	for _, ignoredFile := range ignoredFiles {
		if filename == ignoredFile {
			return true
		}
	}
	return false
}

func isWhitelistedFile(filename string, whitelistedFiles []string) bool {
	for _, whitelistedFile := range whitelistedFiles {
		if filename == whitelistedFile {
			return true
		}
	}
	return false
}
