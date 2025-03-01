package util

import (
	"strings"

	"github.com/google/uuid"
)

func GenerateId() string {
	uid := uuid.New()
	return strings.ReplaceAll(uid.String(), "-", "")
}
