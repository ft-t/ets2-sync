package savefile

import (
	"errors"
	"ets2-sync/dlc"
	"ets2-sync/utils"
	"fmt"
	"strings"
)

type SaveManager struct {
	file         *SaveFile
	supportedDlc dlc.Dlc
}

func NewSaveManager(file *SaveFile, supportedDlc dlc.Dlc) (*SaveManager, error) {
	if file == nil {
		return nil, errors.New("invalid save data")
	}

	mgr := &SaveManager{file: file, supportedDlc: supportedDlc}

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
	tot := offer.MapToDlc()

	if m.file.clientSupportedDlc&tot == tot {
		fmt.Println(tot)
	}

	if strings.Contains(offer.TrailerVariant, "schw_") || strings.Contains(offer.TrailerDefinition, "schw_") {
		return // todo
	}

	if strings.Contains(offer.TrailerVariant, "krone") || strings.Contains(offer.TrailerDefinition, "krone") { // krone
		return // todo
	}

	if strings.Contains(offer.TrailerVariant, "dryliner") || strings.Contains(offer.TrailerDefinition, "dryliner") { // krone
		return // todo
	}

	if utils.Contains(m.file.AvailableCompanies, offer.SourceCompany) &&
		utils.Contains(m.file.AvailableCargoTypes, offer.Cargo) {
		if comp, ok := m.file.companies[offer.SourceCompany]; ok {
			comp.Jobs[offer.Id] = offer
		}
	}
}
