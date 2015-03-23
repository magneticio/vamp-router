package tools

import (
	"code.google.com/p/go-uuid/uuid"
	"strings"
)

func GetUUID() string {
	myUUID := uuid.NewUUID()
	return strings.Split(myUUID.String(), "-")[0]
}
