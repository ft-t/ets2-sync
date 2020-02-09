package savefile

import (
	"bytes"
	"strings"
)

type IConfigSection interface {
	Name() string
	AppendLine(line string)
	//	Write() string // todo
}

type CompanyConfigSection struct {
	name          string
	nameValue     string // nameless
	permanentData string
	raw           bytes.Buffer
	Jobs          []*JobOffer
}

func (c *CompanyConfigSection) Name() string {
	return c.name
}

func (c *CompanyConfigSection) AppendLine(line string) {
	if strings.Contains(line, "cargo_offer_seeds") {
		c.raw.WriteString(line)
	}
}

type JobOffer struct {
	target             string
	expirationTime     string
	urgency            string
	shortestDistanceKm string
	ferryTime          string
	ferryPrice         string
	cargo              string
	companyTruck       string
	trailerVariant     string
	trailerDefinition  string
	unitsCount         string
	fillRatio          string
	trailerPlace       string
}

type JobOfferConfigSection struct {
	name      string
	nameValue string // nameless
	Offer     *JobOffer
}

func (s *JobOfferConfigSection) Name() string {
	return s.name
}

func (s *JobOfferConfigSection) AppendLine(line string) {
	return
}

func (s *JobOfferConfigSection) FillOfferData(fieldName string, value string) {
	if s.Offer == nil {
		s.Offer = &JobOffer{}
	}

	switch strings.Trim(fieldName, ":") {
	case "target":
		s.Offer.target = value
		break
	case "expiration_time":
		s.Offer.expirationTime = value
		break
	case "urgency":
		s.Offer.urgency = value
		break
	case "shortest_distance_km":
		s.Offer.shortestDistanceKm = value
		break
	case "ferry_time":
		s.Offer.ferryTime = value
		break
	case "ferry_price":
		s.Offer.ferryPrice = value
		break
	case "cargo":
		s.Offer.cargo = value
		break
	case "company_truck":
		s.Offer.companyTruck = value
		break
	case "trailer_variant":
		s.Offer.trailerVariant = value
		break
	case "trailer_definition":
		s.Offer.trailerDefinition = value
		break
	case "units_count":
		s.Offer.unitsCount = value
		break
	case "fill_ratio":
		s.Offer.fillRatio = value
		break
	case "trailer_place":
		s.Offer.trailerPlace = value
		break
	}
}

type RawConfigSection struct {
	name      string
	nameValue string // nameless
	raw       bytes.Buffer
}

func (r *RawConfigSection) AppendLine(line string) {
	r.raw.WriteString(line)
}

func (r *RawConfigSection) Name() string {
	return r.name
}
