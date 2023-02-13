package test

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"testing"
)

var WhitelistTest = make(map[string]bool)

func TestReadOneLine(t *testing.T) {
	f, err := os.OpenFile("../pkg/whitelist/whitelist.txt", os.O_RDONLY, 0644)
	if nil != err {
		logrus.Error("failed to open whitelist")
	}
	defer f.Close()
	fInfo, err := f.Stat()
	if nil != err {
		logrus.Error("failed to state file")
	}

	reader := bufio.NewReaderSize(f, int(fInfo.Size()))
	for {
		l, e := reader.ReadBytes('\n')
		WhitelistTest[string(l)] = true
		if io.EOF == e {
			break
		}
	}
	fmt.Println(WhitelistTest["部署脚本说明.pdf"])
	fmt.Println(WhitelistTest["部署脚本说明1.pdf"])

}
