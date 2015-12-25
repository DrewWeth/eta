package models

import (
	"github.com/gocql/gocql"
)

type Querier struct {
	Session *gocql.Session
}
