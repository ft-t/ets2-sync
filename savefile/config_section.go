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
		_, _ = w.Write([]byte(fmt.Sprintf(" job_offer[%d]: %s%s", index, j.Id, newLine)))
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
	SourceCompany      string
	Target             string
	ExpirationTime     string
	Urgency            string
	ShortestDistanceKm string
	FerryTime          string
	FerryPrice         string
	Cargo              string
	CompanyTruck       string
	TrailerVariant     string
	TrailerDefinition  string
	UnitsCount         string
	FillRatio          string
	TrailerPlace       string
	Id                 string // nameParam
}

func (j *JobOffer) Write(w io.Writer, newLine string){
	_, _ = w.Write([]byte(fmt.Sprintf(" target: %s%s", j.Target, newLine)))
	_, _ = w.Write([]byte(fmt.Sprintf(" expiration_time: %s%s", j.ExpirationTime, newLine)))
	_, _ = w.Write([]byte(fmt.Sprintf(" urgency: %s%s", j.Urgency, newLine)))
	_, _ = w.Write([]byte(fmt.Sprintf(" shortest_distance_km: %s%s", j.ShortestDistanceKm, newLine)))
	_, _ = w.Write([]byte(fmt.Sprintf(" ferry_time: %s%s", j.FerryTime, newLine)))
	_, _ = w.Write([]byte(fmt.Sprintf(" ferry_price: %s%s", j.FerryPrice, newLine)))
	_, _ = w.Write([]byte(fmt.Sprintf(" cargo: %s%s", j.Cargo, newLine)))
	_, _ = w.Write([]byte(fmt.Sprintf(" company_truck: %s%s", j.CompanyTruck, newLine)))
	_, _ = w.Write([]byte(fmt.Sprintf(" trailer_variant: %s%s", j.TrailerVariant, newLine)))
	_, _ = w.Write([]byte(fmt.Sprintf(" trailer_definition: %s%s", j.TrailerDefinition, newLine)))
	_, _ = w.Write([]byte(fmt.Sprintf(" units_count: %s%s", j.UnitsCount, newLine)))
	_, _ = w.Write([]byte(fmt.Sprintf(" fill_ratio: %s%s", j.FillRatio, newLine)))
	_, _ = w.Write([]byte(fmt.Sprintf(" trailer_place: %s%s", j.TrailerPlace, newLine)))
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
	panic("should not be called")
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
		s.Offer.Target = value
		break
	case "expiration_time":
		s.Offer.ExpirationTime = value
		break
	case "urgency":
		s.Offer.Urgency = value
		break
	case "shortest_distance_km":
		s.Offer.ShortestDistanceKm = value
		break
	case "ferry_time":
		s.Offer.FerryTime = value
		break
	case "ferry_price":
		s.Offer.FerryPrice = value
		break
	case "cargo":
		s.Offer.Cargo = value
		break
	case "company_truck":
		s.Offer.CompanyTruck = value
		break
	case "trailer_variant":
		s.Offer.TrailerVariant = value
		break
	case "trailer_definition":
		s.Offer.TrailerDefinition = value
		break
	case "units_count":
		s.Offer.UnitsCount = value
		break
	case "fill_ratio":
		s.Offer.FillRatio = value
		break
	case "trailer_place":
		s.Offer.TrailerPlace = value
		break
	case "id":
		s.Offer.Id = value
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
