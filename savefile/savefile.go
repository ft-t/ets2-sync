package savefile

import (
	"bytes"
	"compress/zlib"
	"errors"
	dlc2 "ets2-sync/dlc"
	"ets2-sync/structs"
	"ets2-sync/utils"
	"fmt"
	"io"
	"io/ioutil"
	"sort"
)

type SaveFile struct {
	source              []byte
	offset              int
	lineEndingFormat    string
	AvailableCompanies  []string
	AvailableCargoTypes []string
	configSections      []IConfigSection
	companies           map[string]*CompanyConfigSection
	dlc                 dlc2.Dlc
}

func NewSaveFile(br *bytes.Reader) (*SaveFile, error) {

	if br == nil || br.Size() < 4 {
		return nil, errors.New("invalid source input")
	}

	decrypted, err := tryDecrypt(br)

	if err != nil {
		return nil, err
	}

	r := &SaveFile{source: decrypted, companies: map[string]*CompanyConfigSection{}}
	r.parseConfig(decrypted)

	return r, nil
}

func (s *SaveFile) AddOffer(offer structs.ApplicableOffer) error {
	if v, ok := s.companies[offer.SourceCompany]; ok {
		v.Jobs[offer.Id] = newJobOffer(offer)

		return nil
	}

	return errors.New("can not find company")
}

func (s *SaveFile) ClearOffers() {
	for _, v := range s.companies {
		if v.Jobs == nil {
			continue
		}

		for key, _ := range v.Jobs {
			delete(v.Jobs, key)
		}
	}
}

func (s *SaveFile) ExportOffers() []structs.ApplicableOffer {
	var arr []structs.ApplicableOffer

	for _, k := range s.companies {
		if k.Jobs == nil {
			continue
		}

		for _, j := range k.Jobs {
			job := structs.ApplicableOffer{}
			_, _ = utils.MapToObject(j, &job)

			arr = append(arr, job)
		}
	}

	return arr
}

func (s *SaveFile) Write(w io.Writer) (n int, err error) {
	n, err = w.Write([]byte(fmt.Sprintf("SiiNunit%s", s.lineEndingFormat)))

	if err != nil {
		return n, err
	}

	n, err = w.Write([]byte(fmt.Sprintf("{%s", s.lineEndingFormat)))

	if err != nil {
		return n, err
	}

	writeHeader := func(name string, nameValue string) {
		_, _ = w.Write([]byte(fmt.Sprintf("%s : %s {%s", name, nameValue, s.lineEndingFormat)))
	}
	writeEnd := func() {
		_, _ = w.Write([]byte(fmt.Sprintf("}%s", s.lineEndingFormat)))
		_, _ = w.Write([]byte(fmt.Sprintf("%s", s.lineEndingFormat)))
	}

	for _, k := range s.configSections {
		writeHeader(k.Name(), k.NameValue())
		_, _ = k.Write(w, s.lineEndingFormat) // write struct
		writeEnd()

		if comp, ok := k.(*CompanyConfigSection); ok {
			if comp.Jobs != nil {
				var jobIds []string
				for k, _ := range comp.Jobs {
					jobIds = append(jobIds, k)
				}

				sort.Strings(jobIds)

				for _, j := range jobIds {
					writeHeader("job_offer_data", comp.Jobs[j].Id)
					comp.Jobs[j].Write(w, s.lineEndingFormat)
					writeEnd()
				}
			}
		}
	}

	n, err = w.Write([]byte(fmt.Sprintf("}%s", s.lineEndingFormat)))

	return 0, nil // todo
}

func tryDecrypt(reader *bytes.Reader) ([]byte, error) {
	buff := make([]byte, 4)
	_, err := reader.Read(buff)

	if err != nil {
		return nil, err
	}

	header := string(buff)

	if header == "BSII" {
		return nil, errors.New("binary save format is unsupported")
	}

	if header == "SiiN" { // already plain text
		_, err = reader.Seek(0, io.SeekStart)

		if err != nil {
			return nil, err
		}

		return ioutil.ReadAll(reader)
	}

	if header == "ScsC" { // encrypted
		_, _ = io.CopyN(ioutil.Discard, reader, 32) // skip header related data 32 + 4 bytes

		initVector := make([]byte, 16)
		_, err = reader.Read(initVector)

		if err != nil {
			return nil, err
		}

		_, err := reader.Seek(4, io.SeekCurrent) // skip 4 bytes of length for unpacked data

		if err != nil {
			return nil, err
		}

		encryptedData := make([]byte, reader.Len())

		_, err = reader.Read(encryptedData)

		if err != nil {
			return nil, err
		}

		decrypted, err := decryptSii(encryptedData, []byte{0x2A, 0x5F, 0xCB, 0x17, 0x91, 0xD2, 0x2F, 0xB6, 0x02, 0x45, 0xB3, 0xD8, 0x36,
			0x9E, 0xD0, 0xB2, 0xC2, 0x73, 0x71, 0x56, 0x3F, 0xBF, 0x1F, 0x3C, 0x9E, 0xDF, 0x6B, 0x11, 0x82, 0x5A, 0x5D, 0x0A}, initVector)

		reader = bytes.NewReader(decrypted) // replace reader with decrypted arr

		flReader, err := zlib.NewReader(reader)

		if err != nil {
			return nil, err
		}

		return ioutil.ReadAll(flReader)
	}

	return nil, nil // todo
}
