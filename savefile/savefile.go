package savefile

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io"
	"io/ioutil"
)

type SaveFile struct {
	source           []byte
	offset           int
	lineEndingFormat string
}

func NewSaveFile(br *bytes.Reader) (*SaveFile, error) {

	if br == nil || br.Size() < 4 {
		return nil, errors.New("invalid source input")
	}

	decrypted, err := tryDecrypt(br)

	if err != nil {
		return nil, err
	}

	r := &SaveFile{source: decrypted}

	if string(decrypted[8:10]) == "\r\n" {
		r.lineEndingFormat = "\r\n"
	} else {
		r.lineEndingFormat = "\n"
	}

	return r, nil
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