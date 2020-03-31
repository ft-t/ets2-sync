package main

import (
	"bytes"
	"encoding/json"
	"ets2-sync/db"
	"ets2-sync/dlc"
	"ets2-sync/global"
	"ets2-sync/savefile"
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	global.IsDebug = true

	if err := db.InitializeDb(); err != nil {
		panic(err)
	}

	if err := initOfferManager(); err != nil {
		panic(err)
	}

	Start()
}

func Start() {
	port := os.Getenv("httpPort")

	if len(port) == 0 {
		port = "8080"
	}

	server := &http.Server{
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		Addr:         fmt.Sprintf(":%s", port),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		_ = tmpl.Execute(w, []dlc.Dlc{dlc.GoingEast, dlc.Scandinavia, dlc.LaFrance, dlc.Italy, dlc.BeyondTheBalticSea, dlc.RoadToTheBlackSea, dlc.PowerCargo, dlc.HeavyCargo, dlc.SpecialTransport, dlc.Krone, dlc.Schwarzmuller})
	})

	http.HandleFunc("/dlc", func(w http.ResponseWriter, r *http.Request) {
		res := make(map[string]int)

		for _, d := range dlc.AllDLCs {
			res[d.ToString()] = int(d)
		}

		w.Header().Set("Content-Type", "application/jsons")
		b, _ := json.Marshal(res)

		_, _ = w.Write(b)
		return
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

		dlcs := r.Form["dlc"]
		offersDlcs := dlc.BaseGame
		for _, d := range dlcs {
			val, _ := strconv.Atoi(d)
			offersDlcs |= dlc.Dlc(val)
		}

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			return
		}

		newSaveFile, er := savefile.NewSaveFile(bytes.NewReader(buf.Bytes()))

		if er != nil {
			return // todo
		}

		FillDbWithJobs(newSaveFile.ExportOffers())
		newSaveFile.ClearOffers()
		PopulateOffers(newSaveFile, offersDlcs)

		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", header.Filename))
		w.Header().Set("Content-Type", "application/octet-stream")

		_, _ = newSaveFile.Write(w)
	})

	_ = server.ListenAndServe()
}
