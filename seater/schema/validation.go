package schema

import (
	"github.com/xeipuuv/gojsonschema"
)

// Validate validates schema with document
func Validate(schema, document string) (result *gojsonschema.Result, err error) {
	schemaLoader := gojsonschema.NewStringLoader(schema)
	documentLoader := gojsonschema.NewStringLoader(document)
	return gojsonschema.Validate(schemaLoader, documentLoader)
}