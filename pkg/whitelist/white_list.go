package whitelist

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"os"
)

func init1() {
	f, err := os.OpenFile("whitelist.txt", os.O_RDONLY, 0644)
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
	reader.ReadLine()

}
