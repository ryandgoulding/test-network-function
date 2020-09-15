package main

import (
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"os"
	"path"
)

func main() {
	curDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	schema := path.Join(curDirectory, "ideal-test-schema.json")
	testFile := path.Join(curDirectory, "crazy.json")
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + schema)
	documentLoader := gojsonschema.NewReferenceLoader("file://" + testFile)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}
}
