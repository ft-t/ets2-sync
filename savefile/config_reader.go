package savefile

import (
	"bufio"
	"bytes"
	"strings"
)

func (s *SaveFile) parseConfig(decrypted []byte) {
	if string(decrypted[8:10]) == "\r\n" {
		s.lineEndingFormat = "\r\n"
	} else {
		s.lineEndingFormat = "\n"
	}

	reader := bufio.NewReader(bytes.NewReader(decrypted))

	i := 0

	companies := make(map[*CompanyConfigSection][]string, 0)
	offers := make(map[string]*JobOffer)

	var currentSection IConfigSection

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
				currentSection = &JobOfferConfigSection{name: parsed[0], nameValue: parsed[2]}
				continue
			}
			if parsed[0] == "company" {
				currentSection = &CompanyConfigSection{name: parsed[0], nameValue: parsed[2]}
				continue
			}

			currentSection = &RawConfigSection{name: parsed[0], nameValue: parsed[2]}
			continue
		}

		if len(parsed) == 1 && parsed[len(parsed)-1] == "}" { // end section
			if currentSection == nil {
				break
			}

			s.configSections = append(s.configSections, currentSection)

			if currentSection.Name() == "job_offer_data" {
				m := currentSection.(*JobOfferConfigSection)
				offers[m.nameValue] = m.Offer
				m.FillOfferData("id", m.nameValue)
			}
			currentSection = nil
			continue
		}

		if currentSection == nil {
			continue // should not happen
		}

		currentSection.AppendLine(line)

		if currentSection.Name() == "job_offer_data" {
			currentSection.(*JobOfferConfigSection).FillOfferData(parsed[0], parsed[1])
		}

		if currentSection.Name() == "economy" && len(parsed) > 0 {
			if strings.Contains(parsed[0], "companies[") {
				s.availableCompanies = append(s.availableCompanies, parsed[1])
			}
			if strings.Contains(parsed[0], "transported_cargo_types[") {
				s.availableCargoTypes = append(s.availableCompanies, parsed[1])
			}
		}

		if currentSection.Name() == "player" {
			if parsed[0] == "current_job:" && parsed[1] != "null" {
				s.currentJob = parsed[1]
			}
			if parsed[0] == "selected_job:" && parsed[1] != "null" {
				s.selectedJob = parsed[1]
			}
		}

		if currentSection.Name() == "company" {
			if parsed[0] == "permanent_data:" {
				currentSection.(*CompanyConfigSection).permanentData = parsed[1]
			}
			if parsed[0] == "delivered_trailer:" {
				currentSection.(*CompanyConfigSection).deliveredTrailer = parsed[1]
			}
			if parsed[0] == "delivered_pos:" {
				currentSection.(*CompanyConfigSection).deliveredPos = parsed[1]
			}
			if parsed[0] == "discovered:" {
				currentSection.(*CompanyConfigSection).discovered = parsed[1]
			}
			if parsed[0] == "reserved_trailer_slot:" {
				currentSection.(*CompanyConfigSection).reservedTrailerSlot = parsed[1]
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
				k.Jobs = append(k.Jobs, offer)
			}
		}
	}
}
