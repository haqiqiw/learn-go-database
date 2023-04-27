package learn_go_database

import (
	"testing"
)

func TestOpenConnection(t *testing.T) {
	db := GetConnection()

	defer db.Close()
}
