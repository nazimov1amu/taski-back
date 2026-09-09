package db

import (
	"embed"
	"fmt"
	"strings"
)

//go:embed projects/*.sql tasks/*.sql users/*.sql
var queriesFS embed.FS

type Queries struct {
	GetBulk string
	Get     string
	Create  string
	Update  string
	Delete  string
	GetByEmail string
}

var (
	TaskQueries    = loadQueries("tasks")
	ProjectQueries = loadQueries("projects")
	UserQueries    = loadQueries("users")
)

func loadQueries(entity string) Queries {
	return Queries{
		GetBulk: readQuery(entity, "get_bulk"),
		Get:     readQuery(entity, "get"),
		Create:  readQuery(entity, "create"),
		Update:  readQuery(entity, "update"),
		Delete:  readQuery(entity, "delete"),
		GetByEmail: readQuery(entity, "get_by_email"),
	}
}

func readQuery(entity, name string) string {
	b, err := queriesFS.ReadFile(fmt.Sprintf("%s/%s.sql", entity, name))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}
