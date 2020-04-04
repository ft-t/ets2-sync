package savefile

import (
	"strconv"

	"ets2-sync/dlc_mapper"
	"github.com/mitchellh/hashstructure"
)

type ApplicableOffer struct {
	Id                 string // nameparam
	RequiredDlc        dlc_mapper.Dlc
	SourceCompany      string
	Target             string
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
}

func (o *ApplicableOffer) CalculateHash() string {
	hash, err := hashstructure.Hash(struct {
		S string
		T string
		C string
	}{o.SourceCompany, o.Target, o.Cargo}, nil)

	if err != nil {
		return ""
	}

	return strconv.FormatUint(hash, 10)
}
