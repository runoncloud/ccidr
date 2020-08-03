package ccidr

import (
	"github.com/rakyll/statik/fs"
	"io/ioutil"
	"log"
)

func GetJsonString(filepath string) string {
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	r, err := statikFS.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	return string(contents)
}
