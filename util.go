package gorequests

import (
	"errors"
	"fmt"
	"log"
	"os"
)

func NewFilesMap(files map[string]string) (filemap map[string]*os.File, err error) {
	filemap = make(map[string]*os.File, len(files))
	for key, filename := range files {
		fh, err := os.Open(filename)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error opening file(%s): %v", key, filename))
			log.Println(err)
			return filemap, err
		}
		filemap[key] = fh
	}
	return
}

func NewQueryData(v interface{}) (data map[string][]string, err error) {
	data = make(map[string][]string)
	switch d := v.(type) {
	case map[string]string:
		for k, v := range d {
			data[k] = []string{v}
		}
	case map[string][]string:
		for k, v := range d {
			data[k] = v
		}
	default:
		err = errors.New(fmt.Sprintf("NewQueryData only accepts map[string]string or map[string][]string."))
	}
	return
}
