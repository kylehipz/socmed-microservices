package utils

import (
	"github.com/google/uuid"
)

func ToUUID(id string) uuid.UUID {
	parsedUuid, _ := uuid.Parse(id)

	return parsedUuid
}
