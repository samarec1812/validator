Validator
=================
Validator implements value validations for structs and individual fields based on tags.

Installation
------------

Use go get.

	go get github.com/samarec1812/validator

Then import the validator package into your own code.

	import "github.com/samarec1812/validator"

Struct Validations
------

### Tags:
| Tag | Description   | 
|-----|---------------|
| len | Length        |
| max | Maximum       |
| min | Minimum       | 
| in  | Included With | 
