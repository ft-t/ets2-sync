package dlc_mapper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"ets2-sync/internal"
)

type Game int

const (
	ETS         = 1
	ATS         = 2
	ETS_PROMODS = 3
)

var AllGames = []Game{ETS, ATS, ETS_PROMODS}

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
	Arizona            Dlc = 1 << 12
	Nevada             Dlc = 1 << 13
	NewMexico          Dlc = 1 << 14
	Oregon             Dlc = 1 << 15
	Utah               Dlc = 1 << 17
	Washington         Dlc = 1 << 18
	ForestMachinery    Dlc = 1 << 19
	Idaho 			   Dlc = 1 << 20
)

var AllDLCs = map[Game][]Dlc{
	ETS:         {BaseGame, Krone, Schwarzmuller, Scandinavia, GoingEast, LaFrance, Italy, PowerCargo, HeavyCargo, BeyondTheBalticSea, SpecialTransport, RoadToTheBlackSea},
	ETS_PROMODS: {BaseGame, Krone, Schwarzmuller, Scandinavia, GoingEast, LaFrance, Italy, PowerCargo, HeavyCargo, BeyondTheBalticSea, SpecialTransport, RoadToTheBlackSea},
	ATS:         {BaseGame, Arizona, Nevada, NewMexico, Oregon, Utah, Washington, Idaho, ForestMachinery, SpecialTransport, HeavyCargo},
}

func (t Game) ToString() string {
	switch t {
	case ETS:
		return "ets"
	case ATS:
		return "ats"
	case ETS_PROMODS:
		return "ets-promods"
	}

	return "unk"
}

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
	case Arizona:
		return "arizona"
	case Nevada:
		return "nevada"
	case NewMexico:
		return "new_mexico"
	case Oregon:
		return "oregon"
	case Utah:
		return "utah"
	case Washington:
		return "washington"
	case ForestMachinery:
		return "forest_harvesting"
	case Idaho:
		return "idaho"
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

func GetRequiredDlc(source string, target string, cargo string, trailerDef string, trailerVariant string, game Game) Dlc {
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
	targetCountry := getCountryByCity(targetCity, game)
	sourceCountry := getCountryByCity(sourceCity, game)

	validators := []func() (Dlc, error){
		func() (Dlc, error) { return mapCargoToDlc(cargo, game) },
		func() (Dlc, error) { return mapCompanyToDlc(targetCompany, targetCity, game) },
		func() (Dlc, error) { return mapCompanyToDlc(sourceCompany, sourceCity, game) },
		func() (Dlc, error) { return mapTrailerDefToDlc(trailerDef, targetCountry, sourceCountry, game) },
		func() (Dlc, error) { return mapTrailerVariantToDlc(trailerVariant, game) },
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

func getCountryByCity(city string, game Game) string {
	if cityToCountry != nil {
		s, _ := cityToCountry[city]
		return s
	}

	dir := fmt.Sprintf("data/%s", game.ToString())
	files, er := ioutil.ReadDir(dir)

	if er != nil {
		return ""
	}

	cityToCountry = make(map[string]string)
	for _, file := range files {
		input, _ := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, file.Name()))
		cityItem := map[string]cityFile{}
		_ = json.Unmarshal(input, &cityItem)

		for k, v := range cityItem {
			cityName := k
			if strings.HasPrefix(cityName, "city.") {
				cityName = cityName[5:]
			}

			cityToCountry[cityName] = v.Country
		}
	}

	s, _ := cityToCountry[city]
	return s
}

func mapCompanyToDlc(companyName string, cityName string, game Game) (Dlc, error) {
	for _, d := range AllDLCs[game] {
		if res := readCompanyFile(game, d); res != nil {
			if company, ok := res[companyName]; ok && internal.Contains(company.Cities, cityName) {
				return d, nil
			}
		}
	}

	return None, errors.New("company not found")
}

func mapCargoToDlc(cargoName string, game Game) (Dlc, error) {
	for _, d := range AllDLCs[game] {
		if res := readSimpleJsonArr("cargoes", game, d); res != nil && internal.Contains(res, cargoName) {
			return d, nil
		}
	}

	return None, errors.New("trailer not found")
}

func mapTrailerVariantToDlc(trailerVariant string, game Game) (Dlc, error) {
	for _, d := range AllDLCs[game] {
		if res := readTrailerFile(game, d); res != nil && internal.Contains(res.Variants, trailerVariant) {
			return d, nil
		}
	}

	return None, errors.New("trailer not found")
}

func mapTrailerDefToDlc(trailerDef string, targetCountry string, sourceCountry string, game Game) (Dlc, error) {
	for _, d := range AllDLCs[game] {
		if res := readTrailerFile(game, d); res != nil {
			if def, ok := res.Definitions[trailerDef]; ok {
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

func readTrailerFile(g Game, d Dlc) *trailerFile {
	name := fmt.Sprintf("./data/%s/trailers_%s.json", g.ToString(), d.ToString())

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

func readCompanyFile(g Game, d Dlc) map[string]*companyFile {
	name := fmt.Sprintf("./data/%s/companies_%s.json", g.ToString(), d.ToString())

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

func readSimpleJsonArr(prefix string, g Game, d Dlc) []string {
	path := fmt.Sprintf("./data/%s/%s_%s.json", g.ToString(), prefix, d.ToString())

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
