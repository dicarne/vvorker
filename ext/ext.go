package ext

import _ "embed"

//go:embed ai/dist/index.js
var ExtAiScript string

//go:embed pgsql/dist/index.js
var ExtPgsqlScriptDTS string
