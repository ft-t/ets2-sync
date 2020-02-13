package dlc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"ets2-sync/utils"
)

type Dlc int

const (
	None               Dlc = 0
	BaseGame           Dlc = 1
	GoingEast          Dlc = 1 << 1
	LaFrance           Dlc = 1 << 2
	Italy              Dlc = 1 << 3
	PowerCargo         Dlc = 1 << 4
	HeavyCargo         Dlc = 1 << 5
	BeyondTheBalticSea Dlc = 1 << 6
	Krone              Dlc = 1 << 7
	Schwarzmuller      Dlc = 1 << 8
	Scandinavia        Dlc = 1 << 9
	RoadToTheBlackSea  Dlc = 1 << 10
	SpecialTransport   Dlc = 1 << 11
)

var allDLCs = []Dlc{BaseGame, Krone, Schwarzmuller, Scandinavia, GoingEast, LaFrance, Italy, PowerCargo, HeavyCargo, BeyondTheBalticSea}

func (t Dlc) ToString() string {
	switch t {
	case None:
		return "none"
	case BaseGame:
		return "base_game"
	case Scandinavia:
		return "scandinavia"
	case GoingEast:
		return "going_east"
	case LaFrance:
		return "la_france"
	case Italy:
		return "italy"
	case PowerCargo:
		return "power_cargo"
	case HeavyCargo:
		return "heavy_cargo"
	case BeyondTheBalticSea:
		return "beyond_the_baltic_sea"
	case Krone:
		return "krone"
	case Schwarzmuller:
		return "schwarzmuller"
	case RoadToTheBlackSea:
		return "road_to_the_black_sea"
	case SpecialTransport:
		return "special_transport"
	}

	return "unk"
}

type trailerFile struct {
	Variants   []string `json:"variants"`
	Definition []string `json:"definitions"`
}

func readTrailerFile(d Dlc) *trailerFile {
	data, er := ioutil.ReadFile(fmt.Sprintf("./data/trailers_%s.json", d.ToString()))

	if er != nil {
		return nil // todo log
	}

	r := &trailerFile{}
	_ = json.Unmarshal(data, r)

	return r
}
func readSimpleJsonArr(prefix string, d Dlc) []string {
	data, er := ioutil.ReadFile(fmt.Sprintf("./data/%s_%s.json", prefix, d.ToString()))

	if er != nil {
		return nil // todo log
	}

	var r []string
	_ = json.Unmarshal(data, r)

	return r
}

func MapCompanyToDlc(companyName string) (Dlc, error) {
	for _, d := range allDLCs {
		if res := readSimpleJsonArr("companies", d); res != nil && utils.Contains(res, companyName) {
			return d, nil
		}
	}

	return None, errors.New("trailer not found")
}

func MapCargoToDlc(cargoName string) (Dlc, error) {
	for _, d := range allDLCs {
		if res := readSimpleJsonArr("cargoes", d); res != nil && utils.Contains(res, cargoName) {
			return d, nil
		}
	}

	return None, errors.New("trailer not found")
}

func MapTrailerVariantToDlc(trailerVariant string) (Dlc, error) {
	for _, d := range allDLCs {
		if res := readTrailerFile(d); res != nil && utils.Contains(res.Variants, trailerVariant) {
			return d, nil
		}
	}

	return None, errors.New("trailer not found")
}

func MapTrailerDefToDlc(trailerDef string) (Dlc, error) {
	for _, d := range allDLCs {
		if res := readTrailerFile(d); res != nil && utils.Contains(res.Definition, trailerDef) {
			return d, nil
		}
	}

	return None, errors.New("trailer not found")
}
