package utils

import (
	"strings"

	"github.com/satori/go.uuid"
)

// UUID4 generates an uuid string for version 4
func UUID4() string {
	return uuid.NewV4().String()
}

// UUID generates an uuid string no '-'
func UUID() string {
	return strings.Replace(UUID4(), "-", "", -1)
}
