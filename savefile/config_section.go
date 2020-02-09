package savefile

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type IConfigSection interface {
	Name() string
	NameValue() string
	AppendLine(line string)
	Write(w io.Writer, newLine string) (n int64, err error)
}

type CompanyConfigSection struct {
	name                string
	nameValue           string // nameless
	permanentData       string
	deliveredTrailer    string
	deliveredPos        string
	reservedTrailerSlot string
	discovered          string
	raw                 bytes.Buffer
	Jobs                map[string]*JobOffer
}

func (c *CompanyConfigSection) NameValue() string {
	return c.nameValue
}

func (c *CompanyConfigSection) Write(w io.Writer, newLine string) (n int64, err error) {
	_, _ = w.Write([]byte(fmt.Sprintf(" permanent_data: %s%s", c.permanentData, newLine)))       // todo
	_, _ = w.Write([]byte(fmt.Sprintf(" delivered_trailer: %s%s", c.deliveredTrailer, newLine))) // todo
	_, _ = w.Write([]byte(fmt.Sprintf(" delivered_pos: %s%s", c.deliveredPos, newLine)))         // todo
	_, _ = w.Write([]byte(fmt.Sprintf(" job_offer: %d%s", len(c.Jobs), newLine)))                // todo

	index := 0

	for _, j := range c.Jobs {
		_, _ = w.Write([]byte(fmt.Sprintf(" job_offer[%d]: %s%s", index, j.id, newLine)))
		index++
	}

	_, _ = c.raw.WriteTo(w)                                                         // cargo_offer_seeds
	_, _ = w.Write([]byte(fmt.Sprintf(" discovered: %s%s", c.discovered, newLine))) // todo check
	_, _ = w.Write([]byte(fmt.Sprintf(" reserved_trailer_slot: %s%s", c.reservedTrailerSlot, newLine)))

	return 1, nil
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
	id                 string // nameParam
}

type JobOfferConfigSection struct {
	name      string
	nameValue string // nameless
	Offer     *JobOffer
}

func (s *JobOfferConfigSection) NameValue() string {
	return s.nameValue
}

func (s *JobOfferConfigSection) Write(w io.Writer, newLine string) (n int64, err error) {
	_, _ = w.Write([]byte(fmt.Sprintf(" target: %s%s", s.Offer.target, newLine)))
	_, _ = w.Write([]byte(fmt.Sprintf(" expiration_time: %s%s", s.Offer.expirationTime, newLine))) // todo
	_, _ = w.Write([]byte(fmt.Sprintf(" urgency: %s%s", s.Offer.urgency, newLine)))
	_, _ = w.Write([]byte(fmt.Sprintf(" shortest_distance_km: %s%s", s.Offer.shortestDistanceKm, newLine)))
	_, _ = w.Write([]byte(fmt.Sprintf(" ferry_time: %s%s", s.Offer.ferryTime, newLine)))
	_, _ = w.Write([]byte(fmt.Sprintf(" ferry_price: %s%s", s.Offer.ferryPrice, newLine)))
	_, _ = w.Write([]byte(fmt.Sprintf(" cargo: %s%s", s.Offer.cargo, newLine)))
	_, _ = w.Write([]byte(fmt.Sprintf(" company_truck: %s%s", s.Offer.companyTruck, newLine)))
	_, _ = w.Write([]byte(fmt.Sprintf(" trailer_variant: %s%s", s.Offer.trailerVariant, newLine)))
	_, _ = w.Write([]byte(fmt.Sprintf(" trailer_definition: %s%s", s.Offer.trailerDefinition, newLine)))
	_, _ = w.Write([]byte(fmt.Sprintf(" units_count: %s%s", s.Offer.unitsCount, newLine)))
	_, _ = w.Write([]byte(fmt.Sprintf(" fill_ratio: %s%s", s.Offer.fillRatio, newLine)))
	_, _ = w.Write([]byte(fmt.Sprintf(" trailer_place: %s%s", s.Offer.trailerPlace, newLine)))

	return 1, nil
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
	case "id":
		s.Offer.id = value
		break
	}
}

type RawConfigSection struct {
	name      string
	nameValue string // nameless
	raw       bytes.Buffer
}

func (r *RawConfigSection) NameValue() string {
	return r.nameValue
}

func (r *RawConfigSection) Write(w io.Writer, newLine string) (n int64, err error) {
	return r.raw.WriteTo(w)
}

func (r *RawConfigSection) AppendLine(line string) {
	r.raw.WriteString(line)
}

func (r *RawConfigSection) Name() string {
	return r.name
}
