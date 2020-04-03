package internal

import (
	"ets2-sync/internal"
	"ets2-sync/structs"
	"fmt"
	"io"
)

type JobOffer struct {
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
	Id                 string // nameParam
}

func NewJobOffer(offer structs.ApplicableOffer) *JobOffer {
	job := JobOffer{}
	_, _ = internal.MapToObject(offer, &job)

	return &job
}

func (j *JobOffer) Write(w io.Writer, newLine string) {
	_, _ = w.Write([]byte(fmt.Sprintf(" target: \"%s\"%s", j.Target, newLine)))
	_, _ = w.Write([]byte(fmt.Sprintf(" expiration_time: %s%s", "86400000", newLine)))
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