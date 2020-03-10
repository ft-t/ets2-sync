package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"ets2-sync/db"
	"ets2-sync/global"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"ets2-sync/savefile"
)

func tryMerge() {
	b, _ := ioutil.ReadFile("/home/skydev/go/src/ets2-sync/test.json")
	var jobs []*savefile.JobOffer

	json.Unmarshal(b, &jobs)

	d, _ := ioutil.ReadFile("/home/skydev/go/src/ets2-sync/game.sii")
	r, _ := savefile.NewSaveFile(bytes.NewReader(d))

	save, _ := savefile.NewSaveManager(r)
	save.ClearOffers()

	for _, j := range jobs {
		save.TryAddOffer(j)
	}

	targetPath := "/home/skydev/go/src/ets2-sync/game.sii_2"
	os.Remove(targetPath)
	f, _ := os.Create(targetPath)

	wr := bufio.NewWriter(f)

	_, _ = r.Write(wr)

	//if err != nil {
	//	fmt.Println(err)
	//}

	_ = wr.Flush()
}

func main() {
	global.IsDebug = true
	err := db.InitializeDb()

	if err != nil {
		panic(err)
	}

	tryMerge()
	//d, _ := ioutil.ReadFile("/home/skydev/go/src/ets2-sync/game_damanox.sii")
	//r, _ := savefile.NewSaveFile(bytes.NewReader(d))
	//
	//_, _ = savefile.NewSaveManager(r)
	//
	//b, _ := json.Marshal(r.ExportJobs())
	//x := string(b)
	//fmt.Print(x)
	//save.ClearOffers()
	//
	//targetPath := "/home/skydev/go/src/ets2-sync/game.sii_2"
	//os.Remove(targetPath)
	//f, _ := os.Create(targetPath)
	//
	//wr := bufio.NewWriter(f)
	//
	//_, err = r.Write(wr)
	//
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//_ = wr.Flush()

	//Start()

	//fmt.Println(r, e)
}

func Start() {
	server := &http.Server{
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		Addr:         ":8080",
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		_ = tmpl.Execute(w, nil)
	})
	http.HandleFunc("/save_upload", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseMultipartForm(32 << 20)

		file, header, err := r.FormFile("savefile")
		if err != nil {
			return
		}
		defer file.Close()

		if !strings.HasSuffix(header.Filename, ".sii") {
			return // not a save file
		}

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			return
		}

		newSaveFile, _ := savefile.NewSaveFile(bytes.NewReader(buf.Bytes()))

		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", header.Filename))
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		wr := bufio.NewWriter(w)
		_, _ = newSaveFile.Write(wr)
	})

	_ = server.ListenAndServe()
}
