package schema

// schemas
const (
	id = `{
	    "type": "integer"
	}`

	integer = `{
		"type": "integer"
	}`

	str = `{
		"type": "string",
		"minLength": 1
	}`

	name = `{
		"type": "string",
		"minLength": 1,
		"maxLength": 128
	}`

	email = `{
		"type": "string",
		"format": "email"
  }`

	password = `{
		"type": "string",
		"minLength": 5,
		"maxLength": 128,
		"pattern": "^[^\\s]+$"
  }`

	url = `{
		"type": "string"
	}`
)
