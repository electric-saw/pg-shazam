package auth

import (
	"context"

	"github.com/electric-saw/pg-shazam/internal/pkg/backend"
)

func ValidateUser(conn *backend.Conn, user, pass string) (bool, string) {
	r, err := conn.Query(context.Background(), "select rolpassword from pg_authid where rolname=$1", user)
	if err != nil {
		return false, err.Error()
	}

	var passDb string

	if !r.Next() {
		return false, "Wrong user/password"
	}
	err = r.Scan(&passDb)
	if err != nil {
		return false, err.Error()
	}

	if passDb != pass {
		return false, "Wrong user/password"
	}

	return true, "Success"

}
