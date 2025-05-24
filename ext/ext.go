package ext

import _ "embed"

//go:embed ai/dist/index.js
var ExtAiScript string

//go:embed pgsql/dist/index.js
var ExtPgsqlScriptDTS string

//go:embed oss/dist/index.js
var ExtOSSScriptDTS string

//go:embed kv/dist/index.js
var ExtKVScriptDTS string
