package utils

import (
	"github.com/satori/go.uuid"
)

func CreateUUID() string {
	// 创建
	u1, _ := uuid.NewV4()
	var uuidString = u1.String()
	return uuidString
}

