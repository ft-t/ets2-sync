package savefile

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"ets2-sync/dlc_mapper"
	. "ets2-sync/savefile/internal"
	. "ets2-sync/savefile/internal/sections"
)

func (s *SaveFile) parseConfig(decrypted []byte, game dlc_mapper.Game) {
	if string(decrypted[8:10]) == "\r\n" {
		s.lineEndingFormat = "\r\n"
	} else {
		s.lineEndingFormat = "\n"
	}

	reader := bufio.NewReader(bytes.NewReader(decrypted))

	i := 0

	companies := make(map[*CompanyConfigSection][]string, 0)
	offers := make(map[string]*JobOffer)

	var currentSection ConfigSection

	for {
		i++
		line, er := reader.ReadString('\n')

		if er != nil {
			break
		}

		if i == 1 || i == 2 { // skip header + {
			continue
		}
		parsed := strings.Fields(strings.TrimRight(line, "\r\n"))

		if len(parsed) == 0 { // empty string between sections
			continue
		}

		if parsed[len(parsed)-1] == "{" { // opening configuration block
			if parsed[0] == "job_offer_data" {
				currentSection = NewJobOfferConfigSection(parsed[0], parsed[2])
				continue
			}
			if parsed[0] == "company" {
				currentSection = NewCompanyConfigSection(parsed[0], parsed[2])
				continue
			}

			currentSection = NewRawConfigSection(parsed[0],parsed[2])
			continue
		}

		if len(parsed) == 1 && parsed[len(parsed)-1] == "}" { // end section
			if currentSection == nil {
				break
			}

			if currentSection.Name() == "company" {
				s.companies[currentSection.NameValue()] = currentSection.(*CompanyConfigSection)
			}

			if currentSection.Name() == "job_offer_data" {
				m := currentSection.(*JobOfferConfigSection)
				offers[m.NameValue()] = m.Offer
				m.FillOfferData("id", m.NameValue())

				continue // job_offer_data should not present in configSections
			}

			s.configSections = append(s.configSections, currentSection)

			currentSection = nil
			continue
		}

		if currentSection == nil {
			continue // should not happen
		}

		currentSection.AppendLine(line)

		if currentSection.Name() == "job_offer_data" {
			if len(parsed) == 1 {
				fmt.Println("xer")
			}
			currentSection.(*JobOfferConfigSection).FillOfferData(parsed[0], parsed[1])
		}

		if currentSection.Name() == "economy" && len(parsed) > 0 {
			if strings.Contains(parsed[0], "companies[") {
				s.AvailableCompanies = append(s.AvailableCompanies, parsed[1])
			}
			if strings.Contains(parsed[0], "transported_cargo_types[") {
				s.AvailableCargoTypes = append(s.AvailableCargoTypes, fmt.Sprintf("cargo.%s", parsed[1]))
			}
		}

		if currentSection.Name() == "company" {
			if parsed[0] == "permanent_data:" {
				currentSection.(*CompanyConfigSection).PermanentData = parsed[1]
			}
			if parsed[0] == "delivered_trailer:" {
				currentSection.(*CompanyConfigSection).DeliveredTrailer = parsed[1]
			}
			if parsed[0] == "delivered_pos:" {
				currentSection.(*CompanyConfigSection).DeliveredPos = parsed[1]
			}
			if parsed[0] == "discovered:" {
				currentSection.(*CompanyConfigSection).Discovered = parsed[1]
			}
			if parsed[0] == "reserved_trailer_slot:" {
				currentSection.(*CompanyConfigSection).ReservedTrailerSlot = parsed[1]
			}

			if strings.Contains(parsed[0], "job_offer[") {
				sect := currentSection.(*CompanyConfigSection)

				companies[sect] = append(companies[sect], parsed[1])
			}
		}
	}

	for k, v := range companies {
		if v == nil || len(v) == 0 {
			continue
		}
		for _, jobId := range v {
			if offer, ok := offers[jobId]; ok {
				offer.SourceCompany = k.NameValue()
				k.Jobs[jobId] = offer

				s.dlc |= dlc_mapper.GetRequiredDlc(offer.SourceCompany, offer.Target, offer.Cargo, offer.TrailerDefinition,
					offer.TrailerVariant, game)
			}
		}

		for kJob, job := range k.Jobs {
			if job.Target == "\"\"" || job.Target == "" || job.Target == "null" || job.Cargo == "null" || job.Cargo == "cargo.caravan" {
				delete(k.Jobs, kJob)
			}
		}
	}
}
