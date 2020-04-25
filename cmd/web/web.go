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

var expansionDLCs = map[Game][]Dlc{ETS: {GoingEast, Scandinavia, LaFrance, Italy, BeyondTheBalticSea, RoadToTheBlackSea}, ATS: {Nevada, Arizona, NewMexico, Oregon, Washington, Utah}}
var cargoDLCs = map[Game][]Dlc{ETS: {PowerCargo, HeavyCargo, SpecialTransport}, ATS: {HeavyCargo, SpecialTransport}}
var trailerDLCs = map[Game][]Dlc{ETS: {Schwarzmuller, Krone}, ATS: {}}

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
		res := make(map[int]map[int]*orderedmap.OrderedMap)
		for _, game := range AllGames {
			res[int(game)] = make(map[int]*orderedmap.OrderedMap)

			res[int(game)][1] = orderedmap.New()
			res[int(game)][2] = orderedmap.New()
			res[int(game)][3] = orderedmap.New()
		}

		for g, gameDLCs := range expansionDLCs {
			for _, d := range gameDLCs {
				res[int(g)][1].Set(d.ToString(), int(d))
			}
		}
		for g, gameDLCs := range cargoDLCs {
			for _, d := range gameDLCs {
				res[int(g)][2].Set(d.ToString(), int(d))
			}
		}
		for g, gameDLCs := range trailerDLCs {
			for _, d := range gameDLCs {
				res[int(g)][3].Set(d.ToString(), int(d))
			}
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

		g := r.Form.Get("game")
		dlcs := r.Form["dlc"]

		val, _ := strconv.Atoi(g)
		game := Game(val)

		offersDlcs := BaseGame
		for _, d := range dlcs {
			val, _ := strconv.Atoi(d)
			offersDlcs |= Dlc(val)
		}

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			return
		}

		newSaveFile, er := savefile.NewSaveFile(bytes.NewReader(buf.Bytes()), game)

		if er != nil {
			writeError(er)
			return // todo
		}

		if len(adminPass) > 0 && r.Form.Get("adminPass") == adminPass {
			FillDbWithJobs(newSaveFile.ExportOffers(), game)
		}

		newSaveFile.ClearOffers()
		PopulateOffers(newSaveFile, game, offersDlcs)

		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", header.Filename))
		w.Header().Set("Content-Type", "application/octet-stream")

		_, _ = newSaveFile.Write(w)
	})

	handler := cors.AllowAll().Handler(mux)

	_ = http.ListenAndServe(fmt.Sprintf(":%s", port), handler)
}
