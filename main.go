package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"

	cfg "github.com/Hunga1/go-openapi-transform/pkg/Configuration"
	openapi3 "github.com/getkin/kin-openapi/openapi3"
)

const (
	OAS_DOCS_DIR = "./docs"
)

var totalEndpoints int = 0

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
		// doc, err := transformOASDoc(fp)
		_, err := transformOASDoc(fp)
		if err != nil {
			fmt.Println(err)
		} else {
			// printEndpointData(fp, doc)
			fmt.Println("Successfully transformed OAS doc file: ", fp)
			totalTransformed = totalTransformed + 1
		}
	}

	fmt.Printf("Successfully transformed %v OAS doc files!\n", totalTransformed)
	fmt.Printf("Total API Endpoints transformed: %v\n", totalEndpoints)
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

func transformOASDoc(path string) (*openapi3.T, error) {
	data := getOASDocFilesContent(path)
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	doc, err := loader.LoadFromData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OAS doc file: \"%v\". Error: %v", path, err)
	}

	return doc, nil
}

func getOASDocFilesContent(path string) []byte {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read OAS doc file: %v\nError: %v\n", path, err)
	}

	// Remove "Cc" category control codes except tab, CR, LF
	s := string(content)
	r := regexp.MustCompile("([^\\P{Cc}\t\r\n])")
	matches := r.FindAllStringSubmatch(s, -1)
	for _, v := range matches {
		fmt.Printf("Matched Control character: %v in OAS file: %v\n", strconv.QuoteToASCII(v[1]), path)
	}
	s = r.ReplaceAllString(s, "")

	return []byte(s)
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

func printEndpointData(filepath string, doc *openapi3.T) {
	// var totalEndpoints int = 0

	printPath := func(method string, op *openapi3.Operation) {
		// Method
		fmt.Printf("\t\t%v:\n", method)

		// Parameters
		params := op.Parameters
		fmt.Printf("\t\t\tParameters:\n")
		for _, p := range params {
			value := p.Value
			fmt.Printf("\t\t\t\tName: %v, Type: %v\n", value.Name, value.In)
		}
	}

	// Info
	info := doc.Info
	fmt.Printf("Title: %v\n", info.Title)
	fmt.Printf("API Version: %v\n", info.Version)

	// OpenAPI version
	fmt.Printf("OpenAPI version: %v\n", doc.OpenAPI)

	// Servers
	fmt.Printf("Servers:\n")
	servers := doc.Servers
	for _, s := range servers {
		fmt.Printf("\tURL: %v\n", s.URL)
		fmt.Printf("\tDescription: %v\n", s.Description)
	}

	// Paths
	fmt.Printf("Paths:\n")
	paths := doc.Paths
	for uri, p := range paths {
		// Path URI
		fmt.Printf("\tURI: %v\n", uri)

		// Global path parameters
		parameters := p.Parameters
		fmt.Printf("\t\tGlobal Parameters:\n")
		for _, params := range parameters {
			value := params.Value
			fmt.Printf("\t\t\tName: %v, Type: %v\n", value.Name, value.In)
		}

		// Ops
		if p.Get != nil {
			printPath("GET", p.Get)
			totalEndpoints = totalEndpoints + 1
		}
		if p.Put != nil {
			printPath("PUT", p.Put)
			totalEndpoints = totalEndpoints + 1
		}
		if p.Post != nil {
			printPath("POST", p.Post)
			totalEndpoints = totalEndpoints + 1
		}
		if p.Delete != nil {
			printPath("DELETE", p.Delete)
			totalEndpoints = totalEndpoints + 1
		}
		if p.Options != nil {
			printPath("OPTIONS", p.Options)
			totalEndpoints = totalEndpoints + 1
		}
		if p.Head != nil {
			printPath("HEAD", p.Head)
			totalEndpoints = totalEndpoints + 1
		}
		if p.Patch != nil {
			printPath("PATCH", p.Patch)
			totalEndpoints = totalEndpoints + 1
		}
		if p.Trace != nil {
			printPath("TRACE", p.Trace)
			totalEndpoints = totalEndpoints + 1
		}
	}

	// fmt.Printf("Total API Endpoints found in %v: %v\n", filepath, totalEndpoints)
}
