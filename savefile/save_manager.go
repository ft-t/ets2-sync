package savefile

import (
	"errors"
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

func (m *SaveManager) ClearJobs() {
	for _, v := range m.file.companies {
		if v.Jobs == nil {
			continue
		}

		for key, _ := range v.Jobs {
			//for i, s := range m.file.configSections {
			//	if s.
			//}

			delete(v.Jobs, key)
		}
	}
}
