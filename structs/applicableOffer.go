package structs

import (
	"ets2-sync/db"
	"ets2-sync/dlc"
	"ets2-sync/utils"
	"github.com/mitchellh/hashstructure"
	"strconv"
)

type ApplicableOffer struct {
	Id                 string // nameparam
	Seed               string
	RequiredDlc        dlc.Dlc
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

func NewApplicableOffer(offer *db.DbOffer, id string) ApplicableOffer {
	newOffer := ApplicableOffer{}
	_, _ = utils.MapToObject(offer, &newOffer)

	newOffer.Id = id

	return newOffer
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
