package main

import (
	"bytes"
	"encoding/json"
	"errors"
	savefile2 "ets2-sync/pkg/savefile"
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/iancoleman/orderedmap"
	"github.com/rs/cors"

	"ets2-sync/db"
	"ets2-sync/dlc"
	"ets2-sync/global"
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
	adminPass := os.Getenv("adminPass")

	if len(port) == 0 {
		port = "8080"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		_ = tmpl.Execute(w, []dlc.Dlc{dlc.GoingEast, dlc.Scandinavia, dlc.LaFrance, dlc.Italy, dlc.BeyondTheBalticSea, dlc.RoadToTheBlackSea, dlc.PowerCargo, dlc.HeavyCargo, dlc.SpecialTransport, dlc.Krone, dlc.Schwarzmuller})
	})

	mux.HandleFunc("/stat", func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")

		b, _ := json.Marshal(map[string]interface{}{
			"last_sync" : lastUpdatedSync.Format(time.Stamp),
			"total_offers" : totalOffersForSync,
		})

		_, _ = w.Write(b)
		return
	})

	mux.HandleFunc("/dlc", func(w http.ResponseWriter, r *http.Request) {
		res := make(map[int]*orderedmap.OrderedMap)

		res[1] = orderedmap.New()
		res[2] = orderedmap.New()
		res[3] = orderedmap.New()

		for _, d := range dlc.ExpansionDLCs {
			res[1].Set(d.ToString(), int(d))
		}
		for _, d := range dlc.CargoDLCs {
			res[2].Set(d.ToString(), int(d))
		}
		for _, d := range dlc.TrailerDLCs {
			res[3].Set(d.ToString(), int(d))
		}

		w.Header().Set("Content-Type", "application/json")

		b, _ := json.Marshal(res)

		_, _ = w.Write(b)
		return
	})

	mux.HandleFunc("/save_upload", func(w http.ResponseWriter, r *http.Request) {
		writeError := func (er interface{}){
			b, _ := json.Marshal(map[string]interface{}{
				"error" : r,
			})

			w.WriteHeader(500)
			_, _ = w.Write(b)
		}
		defer func() {
			if r := recover(); r != nil {
				writeError(r)
				return
			}
		}()

		w.Header().Add("Access-Control-Expose-Headers", "Content-Disposition")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		_ = r.ParseMultipartForm(32 << 20)

		file, header, err := r.FormFile("savefile")
		if err != nil {
			writeError(err)
			return
		}
		defer file.Close()

		if !strings.HasSuffix(header.Filename, ".sii") {
			writeError(errors.New("not a .sii file"))

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

		newSaveFile, er := savefile2.NewSaveFile(bytes.NewReader(buf.Bytes()))

		if er != nil {
			writeError(er)
			return // todo
		}

		if len(adminPass) > 0 && r.Form.Get("adminPass") == adminPass {
			FillDbWithJobs(newSaveFile.ExportOffers())
		}

		newSaveFile.ClearOffers()
		PopulateOffers(newSaveFile, offersDlcs)

		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", header.Filename))
		w.Header().Set("Content-Type", "application/octet-stream")

		_, _ = newSaveFile.Write(w)
	})

	handler := cors.AllowAll().Handler(mux)

	_ = http.ListenAndServe(fmt.Sprintf(":%s", port),handler)
}
