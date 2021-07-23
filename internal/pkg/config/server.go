package config

import "github.com/electric-saw/pg-shazam/pkg/util"

type Server string

func (s Server) ConnectionString() string {
	return string(s)
}

func (s Server) ConnectionStringHiddenPass() string {
	return util.HiddePass(s.ConnectionString())
}
