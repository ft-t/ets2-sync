package sections

import (
	"bytes"
	. "ets2-sync/pkg/savefile/internal"
	"fmt"
	"io"
	"sort"
	"strings"
)

type CompanyConfigSection struct {
	name                string
	nameValue           string // nameless
	PermanentData       string
	DeliveredTrailer    string
	DeliveredPos        string
	ReservedTrailerSlot string
	Discovered          string
	raw                 bytes.Buffer
	Jobs                map[string]*JobOffer
}

func NewCompanyConfigSection(name string, nameValue string) *CompanyConfigSection {
	sect := CompanyConfigSection{
		name:      name,
		nameValue: nameValue,
		raw:       bytes.Buffer{},
		Jobs:      map[string]*JobOffer{},
	}

	return &sect
}
func (c *CompanyConfigSection) NameValue() string {
	return c.nameValue
}

func (c *CompanyConfigSection) Write(w io.Writer, newLine string) (n int64, err error) {
	_, _ = w.Write([]byte(fmt.Sprintf(" permanent_data: %s%s", c.PermanentData, newLine)))
	_, _ = w.Write([]byte(fmt.Sprintf(" delivered_trailer: %s%s", c.DeliveredTrailer, newLine)))
	_, _ = w.Write([]byte(fmt.Sprintf(" delivered_pos: %s%s", c.DeliveredPos, newLine)))
	_, _ = w.Write([]byte(fmt.Sprintf(" job_offer: %d%s", len(c.Jobs), newLine)))

	var jobIds []string

	for k, _ := range c.Jobs {
		jobIds = append(jobIds, k)
	}

	sort.Strings(jobIds)

	index := 0

	for _, id := range jobIds {
		_, _ = w.Write([]byte(fmt.Sprintf(" job_offer[%d]: %s%s", index, id, newLine)))
		index++
	}

	_, _ = c.raw.WriteTo(w) // cargo_offer_seeds
	_, _ = w.Write([]byte(fmt.Sprintf(" discovered: %s%s", c.Discovered, newLine)))
	_, _ = w.Write([]byte(fmt.Sprintf(" reserved_trailer_slot: %s%s", c.ReservedTrailerSlot, newLine)))

	return 1, nil
}

func (c *CompanyConfigSection) Name() string {
	return c.name
}

func (c *CompanyConfigSection) AppendLine(line string) {
	if strings.Contains(line, "cargo_offer_seeds") {
		c.raw.WriteString(line)
	}
}
