package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	. "ets2-sync/dlc_mapper"
	"ets2-sync/savefile"
	"github.com/iancoleman/orderedmap"
	"github.com/pkg/errors"
	"github.com/rs/cors"
)

var expansionDLCs = []Dlc{GoingEast, Scandinavia, LaFrance, Italy, BeyondTheBalticSea, RoadToTheBlackSea}
var cargoDLCs = []Dlc{PowerCargo, HeavyCargo, SpecialTransport}
var trailerDLCs = []Dlc{Schwarzmuller, Krone}

func Run() {
	rand.Seed(time.Now().UTC().UnixNano())

	if err := InitializeDb(); err != nil {
		panic(err)
	}

	if err := initOfferManager(); err != nil {
		panic(err)
	}

	start()
}

func start() {
	port := os.Getenv("httpPort")
	adminPass := os.Getenv("adminPass")

	if len(port) == 0 {
		port = "8080"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/stat", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		b, _ := json.Marshal(map[string]interface{}{
			"last_sync":    lastUpdatedSync.Format(time.Stamp),
			"total_offers": totalOffersForSync,
		})

		_, _ = w.Write(b)
		return
	})

	mux.HandleFunc("/dlc", func(w http.ResponseWriter, r *http.Request) {
		res := make(map[int]*orderedmap.OrderedMap)

		res[1] = orderedmap.New()
		res[2] = orderedmap.New()
		res[3] = orderedmap.New()

		for _, d := range expansionDLCs {
			res[1].Set(d.ToString(), int(d))
		}
		for _, d := range cargoDLCs {
			res[2].Set(d.ToString(), int(d))
		}
		for _, d := range trailerDLCs {
			res[3].Set(d.ToString(), int(d))
		}

		w.Header().Set("Content-Type", "application/json")

		b, _ := json.Marshal(res)

		_, _ = w.Write(b)
		return
	})

	mux.HandleFunc("/save_upload", func(w http.ResponseWriter, r *http.Request) {
		writeError := func(er interface{}) {
			b, _ := json.Marshal(map[string]interface{}{
				"error": r,
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
		offersDlcs := BaseGame
		for _, d := range dlcs {
			val, _ := strconv.Atoi(d)
			offersDlcs |= Dlc(val)
		}

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			return
		}

		newSaveFile, er := savefile.NewSaveFile(bytes.NewReader(buf.Bytes()))

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

	_ = http.ListenAndServe(fmt.Sprintf(":%s", port), handler)
}