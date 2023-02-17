package whitelist

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var Whitelist = make(map[string]bool)

func init() {
	f, err := os.OpenFile("./pkg/whitelist/whitelist.txt", os.O_RDONLY, 0644)
	if nil != err {
		logrus.Error("failed to open whitelist")
		return
	}
	defer f.Close()
	fInfo, err := f.Stat()
	if nil != err {
		logrus.Error("failed to state file")
		return
	}
	reader := bufio.NewReaderSize(f, int(fInfo.Size()))
	for {
		l, e := reader.ReadBytes('\n')
		Whitelist[string(l)] = true
		if io.EOF == e {
			break
		}
	}
	return
}

func CheckList(file string) bool {
	if Whitelist[file] {
		return true
	}
	return false
}
