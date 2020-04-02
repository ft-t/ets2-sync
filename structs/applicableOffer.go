package structs

import (
	"ets2-sync/db"
	"ets2-sync/dlc"
	"ets2-sync/utils"
)

type ApplicableOffer struct {
	Id                 string // nameparam
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
