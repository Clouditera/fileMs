package test

import (
	"fmt"
	"path"
	"strings"
	"testing"
)

func TestUrlJoin(t *testing.T) {
	uuid := "0222e0b9-6c25-4d41-9afa-6f561ed438d3"
	fileName := "1.txt"
	//objectName := strings.TrimPrefix(path.Join("breakpoint", path.Join(uuid[0:1], uuid[1:2], uuid, fileName)), "/")
	objectName := strings.TrimPrefix(path.Join("breakpoint", uuid[0:1], uuid[1:2], uuid, fileName), "/")
	fmt.Println(objectName)
}
