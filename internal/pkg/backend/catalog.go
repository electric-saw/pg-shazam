package backend

import (
	"fmt"
)

func InitShazamCatalog(shazam *Shazam) {
	shazam.RunAllPrimaryHosts(string(createTableShazamShardDefinition()))
}

func InsertCatalogShardDefinition(shazam *Shazam, tableName string, columns []string) []error {
	for _, column := range columns {
		err := shazam.RunAllPrimaryHosts(fmt.Sprintf("INSERT INTO _shazam_shard_definition (table_name, column_name) VALUES('%v', '%v');", tableName, column))

		if err != nil {
			return err
		}
	}

	return nil
}

func createTableShazamShardDefinition() string {
	return `
		CREATE TABLE IF NOT EXISTS _shazam_shard_definition (
			table_name text,
			column_name text,
			PRIMARY KEY (table_name, column_name)
		);
	`
}
