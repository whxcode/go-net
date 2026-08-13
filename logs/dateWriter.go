package logs

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type DateWrite struct {
	dir   string
	file  *os.File
	today string
	mtx   sync.Mutex
	size  int64 // MB
}

var MB_UNIT int64 = 1024 * 1024

func (d *DateWrite) rotate() {
	if d.file != nil {
		d.file.Close()
	}

	today := time.Now().Format("2006-01-02")

	dir := fmt.Sprintf("%s/%s", d.dir, today)

	os.MkdirAll(dir, 0o755)

	date := time.Now().Format("2006-01-02 15:04:05")

	filename := fmt.Sprintf("%s/%s.json", dir, date)
	fmt.Println("filename:", filename)

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		panic(err)
	}

	d.file = f
	d.today = today
}

func (d *DateWrite) Write(p []byte) (n int, err error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	if d.today != time.Now().Format("2006-01-02") {
		d.rotate()
	}

	fileInfo, _ := d.file.Stat()

	if fileInfo.Size()*MB_UNIT > int64(d.size)*MB_UNIT {
		d.rotate()
	}

	reader := bytes.NewReader(p)

	n6, err := io.Copy(d.file, reader)

	return int(n6), nil
}

func MakeDateWrite(dir string, size int64) *DateWrite {
	os.MkdirAll(dir, 0o755)

	return &DateWrite{
		dir:  dir,
		size: size,
	}
}
