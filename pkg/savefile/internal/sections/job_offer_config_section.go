package sections

import (
	. "ets2-sync/pkg/savefile/internal"
	"io"
	"strings"
)

type JobOfferConfigSection struct {
	name      string
	nameValue string // nameless
	Offer     *JobOffer
}

func NewJobOfferConfigSection(name string, nameValue string) *JobOfferConfigSection {
	sect := JobOfferConfigSection{
		name:      name,
		nameValue: nameValue,
	}

	return &sect
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
		s.Offer.Target = strings.Trim(value, "\"")
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