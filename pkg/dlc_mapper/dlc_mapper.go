package dlc_mapper

import (
	"encoding/json"
	"errors"
	"ets2-sync/internal"
	"fmt"
	"io/ioutil"
	"strings"
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

var AllDLCs = []Dlc{BaseGame, Krone, Schwarzmuller, Scandinavia, GoingEast, LaFrance, Italy, PowerCargo, HeavyCargo, BeyondTheBalticSea, SpecialTransport, RoadToTheBlackSea}

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
	Variants    []string                     `json:"variants"`
	Definitions map[string]trailerDefinition `json:"definitions"`
}

type trailerDefinition struct {
	Countries []string `json:"countries"`
}

type companyFile struct {
	Cities     []string `json:"cities"`
	CargoesIn  []string `json:"cargoes_in"`
	CargoesOut []string `json:"cargoes_out"`
}

type cityFile struct {
	Country string `json:"country"`
}

func GetRequiredDlc(source string, target string, cargo string, trailerDef string, trailerVariant string) Dlc {
	getCompanyAndCity := func(str string) (city string, company string) {
		if len(str) > 0 {
			str = strings.Replace(str, "\"", "", 2)
			companyData := strings.Split(str, ".")
			return companyData[len(companyData)-1], companyData[len(companyData)-2]
		}

		return "", ""
	}

	targetCity, targetCompany := getCompanyAndCity(target)
	sourceCity, sourceCompany := getCompanyAndCity(source)
	targetCountry := getCountryByCity(targetCity)
	sourceCountry := getCountryByCity(sourceCity)


	validators := []func() (Dlc, error){
		func() (Dlc, error) { return mapCargoToDlc(cargo) },
		func() (Dlc, error) { return mapCompanyToDlc(targetCompany, targetCity) },
		func() (Dlc, error) { return mapCompanyToDlc(sourceCompany, sourceCity) },
		func() (Dlc, error) { return mapTrailerDefToDlc(trailerDef, targetCountry, sourceCountry) },
		func() (Dlc, error) { return mapTrailerVariantToDlc(trailerVariant) },
	}

	totalDlc := None

	for _, v := range validators {
		parsed, er := v()

		if er != nil {
			return None
		}

		if parsed == None {
			return None
		}

		totalDlc |= parsed
	}

	return totalDlc
}

var cityToCountry map[string]string

func getCountryByCity(city string) string {
	if cityToCountry != nil {
		s, _ := cityToCountry[city]
		return s
	}

	dir := "data"
	files, er := ioutil.ReadDir(dir)

	if er != nil {
		return ""
	}

	cityToCountry = make(map[string]string)
	for _, file := range files {
		input, _ := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, file.Name()))
		cityItem := map[string]cityFile{}
		_ = json.Unmarshal(input, &cityItem)

		for k,v := range cityItem {
			cityName := k
			if strings.HasPrefix(cityName,"city."){
				cityName = cityName[5:]
			}

			cityToCountry[cityName] = v.Country
		}
	}

	s, _ := cityToCountry[city]
	return s
}

func mapCompanyToDlc(companyName string, cityName string) (Dlc, error) {
	for _, d := range AllDLCs {
		if res := readCompanyFile(d); res != nil {
			if company, ok := res[companyName]; ok && internal.Contains(company.Cities, cityName) {
				return d, nil
			}
		}
	}

	return None, errors.New("company not found")
}


func mapCargoToDlc(cargoName string) (Dlc, error) {
	for _, d := range AllDLCs {
		if res := readSimpleJsonArr("cargoes", d); res != nil && internal.Contains(res, cargoName) {
			return d, nil
		}
	}

	return None, errors.New("trailer not found")
}

func mapTrailerVariantToDlc(trailerVariant string) (Dlc, error) {
	for _, d := range AllDLCs {
		if res := readTrailerFile(d); res != nil && internal.Contains(res.Variants, trailerVariant) {
			return d, nil
		}
	}

	return None, errors.New("trailer not found")
}

func mapTrailerDefToDlc(trailerDef string, targetCountry string, sourceCountry string) (Dlc, error) {
	for _, d := range AllDLCs {
		if res := readTrailerFile(d); res != nil {
			if def, ok := res.Definitions[trailerDef]; ok  {
				if len(def.Countries) == 0 { // that trailer is allowed for all counties
					return d, nil
				}

				if internal.Contains(def.Countries, targetCountry) && internal.Contains(def.Countries, sourceCountry) {
					return d, nil
				}

				return None, errors.New("invalid target or source country")
			}
		}
	}

	return None, errors.New("trailer not found")
}

var parsedTrailerFile = make(map[string]*trailerFile)
var parsedCompanyFile = make(map[string]map[string]*companyFile)

var dataCache = make(map[string][]string)

func readTrailerFile(d Dlc) *trailerFile {
	name := fmt.Sprintf("./data/trailers_%s.json", d.ToString())

	if v, ok := parsedTrailerFile[name]; ok {
		return v
	}

	data, er := ioutil.ReadFile(name)

	if er != nil {
		return nil // todo log
	}

	parsed := &trailerFile{}
	_ = json.Unmarshal(data, parsed)

	parsedTrailerFile[name] = parsed

	return parsed
}

func readCompanyFile(d Dlc) map[string]*companyFile {
	name := fmt.Sprintf("./data/companies_%s.json", d.ToString())

	if v, ok := parsedCompanyFile[name]; ok {
		return v
	}

	data, er := ioutil.ReadFile(name)

	if er != nil {
		return nil // todo log
	}

	parsed := make(map[string]*companyFile)
	_ = json.Unmarshal(data, &parsed)

	parsedCompanyFile[name] = parsed

	return parsed
}

func readSimpleJsonArr(prefix string, d Dlc) []string {
	path := fmt.Sprintf("./data/%s_%s.json", prefix, d.ToString())

	if v, ok := dataCache[path]; ok {
		return v
	}

	data, er := ioutil.ReadFile(path)

	if er != nil {
		return nil // todo log
	}

	var r []string
	_ = json.Unmarshal(data, &r)

	dataCache[path] = r

	return r
}
