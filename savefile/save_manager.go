package savefile

import (
	"errors"
	"strings"
)

type SaveManager struct {
	file *SaveFile
}

func NewSaveManager(file *SaveFile) (*SaveManager, error) {
	if file == nil {
		return nil, errors.New("invalid save data")
	}

	mgr := &SaveManager{file: file}

	return mgr, nil
}

func (m *SaveManager) ClearOffers() {
	for _, v := range m.file.companies {
		if v.Jobs == nil {
			continue
		}

		for key, _ := range v.Jobs {
			delete(v.Jobs, key)
		}
	}
}

func (m *SaveManager) TryAddOffer(offer *JobOffer) {
	if strings.Contains(offer.TrailerVariant, "schw_")|| strings.Contains(offer.TrailerDefinition, "schw_"){
		return // todo
	}

	if strings.Contains(offer.TrailerVariant, "krone")|| strings.Contains(offer.TrailerDefinition, "krone"){ // krone
		return // todo
	}

	if strings.Contains(offer.TrailerVariant, "dryliner")|| strings.Contains(offer.TrailerDefinition, "dryliner"){ // krone
		return // todo
	}

	if contains(m.file.AvailableCompanies, offer.SourceCompany) &&
		contains(m.file.AvailableCargoTypes, offer.Cargo) {
		if comp, ok := m.file.companies[offer.SourceCompany]; ok {
			 comp.Jobs[offer.Id] = offer
		}
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
