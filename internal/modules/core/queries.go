package core

import "embed"

// sqlFiles embeds every named .sql file in sql/. It lives at the module
// root (not inside repository/postgres) because sql/ is a module-level
// concern shared by whichever data-access implementation the module
// uses — go:embed cannot reach outside its own package directory tree,
// so postgres repositories receive query strings as constructor
// arguments instead of embedding the folder themselves.
//
//go:embed sql/*.sql
var sqlFiles embed.FS

// query reads an embedded .sql file by name. Called once per repository
// construction in module.go, never on the hot path.
func query(name string) string {
	b, err := sqlFiles.ReadFile("sql/" + name)
	if err != nil {
		// A missing embedded query is a build-time packaging error, not a
		// runtime condition callers can recover from.
		panic("core: missing embedded query: " + name)
	}
	return string(b)
}
