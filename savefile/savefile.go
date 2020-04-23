package savefile

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"sort"

	"ets2-sync/decryptor"
	"ets2-sync/dlc_mapper"
	"ets2-sync/internal"
	. "ets2-sync/savefile/internal"
	. "ets2-sync/savefile/internal/sections"
)

type SaveFile struct {
	source              []byte
	offset              int
	lineEndingFormat    string
	AvailableCompanies  []string
	AvailableCargoTypes []string
	configSections      []ConfigSection
	companies           map[string]*CompanyConfigSection
	dlc                 dlc_mapper.Dlc
}

func NewSaveFile(br *bytes.Reader, game dlc_mapper.Game) (*SaveFile, error) {

	if br == nil || br.Size() < 4 {
		return nil, errors.New("invalid source input")
	}

	decrypted, err := tryDecrypt(br)

	if err != nil {
		return nil, err
	}

	r := &SaveFile{source: decrypted, companies: map[string]*CompanyConfigSection{}}
	r.parseConfig(decrypted, game)

	return r, nil
}

func (s *SaveFile) AddOffer(offer ApplicableOffer) error {
	if v, ok := s.companies[offer.SourceCompany]; ok {

		job := JobOffer{}
		_, _ = internal.MapToObject(offer, &job)

		v.Jobs[offer.Id] = &job

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

func (s *SaveFile) ExportOffers() []ApplicableOffer {
	var arr []ApplicableOffer

	for _, k := range s.companies {
		if k.Jobs == nil {
			continue
		}

		for _, j := range k.Jobs {
			job := ApplicableOffer{}
			_, _ = internal.MapToObject(j, &job)

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

		decrypted, err := decryptor.NewSiiDecryptor(initVector).Decrypt(encryptedData)

		reader = bytes.NewReader(decrypted) // replace reader with decrypted arr

		flReader, err := zlib.NewReader(reader)

		if err != nil {
			return nil, err
		}

		return ioutil.ReadAll(flReader)
	}

	return nil, nil // todo
}
