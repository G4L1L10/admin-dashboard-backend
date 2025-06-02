// utils/closer.go
package utils

import (
	"io"
	"log"
)

// SafeClose handles closing objects that return an error
func SafeClose(closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Printf("error during Close(): %v", err)
	}
}
