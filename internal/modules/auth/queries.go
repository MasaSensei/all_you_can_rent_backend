package auth

import "embed"

//go:embed sql/*.sql
var sqlFiles embed.FS

func query(name string) string {
	b, err := sqlFiles.ReadFile("sql/" + name)
	if err != nil {
		panic("auth: missing embedded query: " + name)
	}
	return string(b)
}
