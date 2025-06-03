package ext

import _ "embed"

//go:embed ai/dist/index.js
var ExtAiScript string

//go:embed pgsql/dist/index.js
var ExtPgsqlScript string

//go:embed oss/dist/index.js
var ExtOSSScript string

//go:embed kv/dist/index.js
var ExtKVScript string

//go:embed assets/dist/index.js
var ExtAssetsScript string
