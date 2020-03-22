package dlc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

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

type companyFile struct {
	Cities     []string `json:"cities"`
	CargoesIn  []string `json:"cargoes_in"`
	CargoesOut []string `json:"cargoes_out"`
}

func GetRequiredDlc(source string, target string, cargo string, trailerDef string, trailerVariant string) Dlc {
	getCompanyAndCity := func(str string) (city string, company string) {
		if len(str) > 0 {
			str = strings.Replace(str, "\"", "", 2)
			companyData := strings.Split(str, ".")
			return companyData[1], companyData[0]
		}

		return "", ""
	}

	targetCity, targetCompany := getCompanyAndCity(target)
	sourceCity, sourceCompany := getCompanyAndCity(source)


	dlc1, _ := mapCargoToDlc(cargo)
	dlc2, _ := mapCompanyToDlc(targetCompany, targetCity)
	dlc3, _ := mapCompanyToDlc(sourceCompany, sourceCity)
	dlc4, _ := mapTrailerDefToDlc(trailerDef)
	dlc5, _ := mapTrailerVariantToDlc(trailerVariant)

	totalDlc := dlc1 | dlc2 | dlc3 | dlc4 | dlc5

	return totalDlc
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

func readCompanyFile(d Dlc) map[string]*companyFile {
	data, er := ioutil.ReadFile(fmt.Sprintf("./data/companies_%s.json", d.ToString()))

	if er != nil {
		return nil // todo log
	}

	r := make(map[string]*companyFile)
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

func mapCompanyToDlc(companyName string, cityName string) (Dlc, error) {
	for _, d := range allDLCs {
		if res := readCompanyFile(d); res != nil {
			if company, ok := res[companyName]; ok && utils.Contains(company.Cities, cityName) {
				return d, nil
			}
		}
	}

	return None, errors.New("company not found")
}

func mapCargoToDlc(cargoName string) (Dlc, error) {
	for _, d := range allDLCs {
		if res := readSimpleJsonArr("cargoes", d); res != nil && utils.Contains(res, cargoName) {
			return d, nil
		}
	}

	return None, errors.New("trailer not found")
}

func mapTrailerVariantToDlc(trailerVariant string) (Dlc, error) {
	for _, d := range allDLCs {
		if res := readTrailerFile(d); res != nil && utils.Contains(res.Variants, trailerVariant) {
			return d, nil
		}
	}

	return None, errors.New("trailer not found")
}

func mapTrailerDefToDlc(trailerDef string) (Dlc, error) {
	for _, d := range allDLCs {
		if res := readTrailerFile(d); res != nil && utils.Contains(res.Definition, trailerDef) {
			return d, nil
		}
	}

	return None, errors.New("trailer not found")
}
