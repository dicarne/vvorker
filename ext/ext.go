package ext

import _ "embed"

//go:embed ai/dist/index.js
var ExtAiScript string

//go:embed pgsql/dist/index.js
var ExtPgsqlScript string

//go:embed mysql/dist/index.js
var ExtMysqlScript string

//go:embed oss/dist/index.js
var ExtOSSScript string

//go:embed kv/dist/index.js
var ExtKVScript string

//go:embed assets/dist/index.js
var ExtAssetsScript string

//go:embed task/dist/index.js
var ExtTaskScript string

////////////////////////////

//go:embed ai/src/binding.ts
var TypeBindingAI string

//go:embed pgsql/src/binding.ts
var TypeBindingPgsql string

//go:embed mysql/src/binding.ts
var TypeBindingMysql string

//go:embed oss/src/binding.ts
var TypeBindingOSS string

//go:embed kv/src/binding.ts
var TypeBindingKV string

//go:embed assets/src/binding.ts
var TypeBindingAssets string

//go:embed task/src/binding.ts
var TypeBindingTask string

///////////////////////////

//go:embed control/dist/index.js
var ControlScript string
