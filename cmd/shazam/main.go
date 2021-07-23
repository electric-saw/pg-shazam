package main

import (
	"os"
	"path/filepath"

	"github.com/electric-saw/pg-shazam/pkg/shazam"
	"github.com/electric-saw/pg-shazam/pkg/util"
)

func main() {
	baseName := filepath.Base(os.Args[0])

	err := shazam.NewShazamCommand(baseName).Execute()
	util.CheckErr(err)
}
