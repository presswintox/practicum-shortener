package migrations

import "embed"

// Files contains all SQL migrations compiled into the application binary.
//
//go:embed *.sql
var Files embed.FS
