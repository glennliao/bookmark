package inits

import _ "embed"

//go:embed _access.json
var AccessConfigJson string

//go:embed _reuqest.json
var RequestConfigJson string
