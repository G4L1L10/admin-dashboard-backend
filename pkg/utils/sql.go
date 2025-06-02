// utils/sql.go
package utils

import (
	"database/sql"
	"log"
)

// SafeCloseRows ensures rows.Close() is called and logs any error.
func SafeCloseRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		log.Printf("error closing rows: %v", err)
	}
}
