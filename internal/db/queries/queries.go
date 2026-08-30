package db

import (
	"embed"
	"fmt"
	"strings"
)

//go:embed notes/*.sql projects/*.sql tasks/*.sql
var queriesFS embed.FS

type Queries struct {
	GetBulk string
	Get     string
	Create  string
	Update  string
	Delete  string
	Search  string
}

var (
	TaskQueries    = loadQueries("tasks")
	ProjectQueries = loadQueries("projects")
	NoteQueries    = loadQueries("notes")
)

func loadQueries(entity string) Queries {
	return Queries{
		GetBulk: mustQuery(entity, "get_bulk"),
		Get:     mustQuery(entity, "get"),
		Create:  mustQuery(entity, "create"),
		Update:  mustQuery(entity, "update"),
		Delete:  mustQuery(entity, "delete"),
		Search:  optionalQuery(entity, "search"),
	}
}

func mustQuery(entity, name string) string {
	b, err := queriesFS.ReadFile(fmt.Sprintf("%s/%s.sql", entity, name))
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(b))
}

func optionalQuery(entity, name string) string {
	b, err := queriesFS.ReadFile(fmt.Sprintf("%s/%s.sql", entity, name))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}
