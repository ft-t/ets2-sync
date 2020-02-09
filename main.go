package main

import (
	"bufio"
	"bytes"
	"ets2-sync/savefile"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	d, err := ioutil.ReadFile("/home/skydev//go/src/ets2-sync/game.sii")
	r, _ := savefile.NewSaveFile(bytes.NewReader(d))

	save, _ := savefile.NewSaveManager(r)
	save.ClearJobs()

	targetPath := "/home/skydev/go/src/ets2-sync/game.sii_2"
	os.Remove(targetPath)
	f, _ := os.Create(targetPath)

	wr := bufio.NewWriter(f)

	_, err = r.Write(wr)

	if err != nil {
		fmt.Println(err)
	}

	_ = wr.Flush()


	//fmt.Println(r, e)
}
