package main

import (
	"bytes"
	"ets2-sync/savefile"
	"fmt"
	"io/ioutil"
)

func main() {
	d, e := ioutil.ReadFile("/home/skydev//go/src/ets2-sync/game.sii")
	r, _ := savefile.NewSaveFile(bytes.NewReader(d))

	fmt.Println(r,e)
}
