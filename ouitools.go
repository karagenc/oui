// Package go-oui provides functions to work with MAC and OUI's
package ouidb

import (
	"bufio"
	"encoding/hex"
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var ErrInvalidMACAddress = errors.New("invalid MAC address")


// Helper functions

func macToUint64(address [6]byte) uint64 {
	var a uint64
	for _, x := range address {
		a <<= 8
		a |= uint64(x)
	}
	return a
}

func maskToUint64(mask uint8) uint64 {
	return ^(uint64(1)<<(48-mask) - 1)
}

func parseMAC(s string) ([6]byte, error) {
	var hw [6]byte

	oct := strings.FieldsFunc(s,
		func(r rune) bool { return r == ':' || r == '-' })

	_, err := hex.Decode(hw[:], []byte(strings.Join(oct, "")))
	if err != nil {
		return hw, err
	}

	return hw, nil
}


// oui, mask, organization
type addressBlock struct {
	Oui          [6]byte
	Mask         byte
	Organization string
}

type OuiDB struct {
	blocks []addressBlock
}

func (m *OuiDB) load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return (err)
	}

	fieldsRe := regexp.MustCompile(`^(\S+)\t+(\S+)(\s+#\s+(\S.*))?`)

	re := regexp.MustCompile(`((?:(?:[0-9a-zA-Z]{2})[-:]){2,5}(?:[0-9a-zA-Z]{2}))(?:/(\w{1,2}))?`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" || text[0] == '#' || text[0] == '\t' {
			continue
		}

		block := addressBlock{}

		// Split input text into address, short organization name
		// and full organization name
		fields := fieldsRe.FindAllStringSubmatch(text, -1)
		addr := fields[0][1]
		if fields[0][4] != "" {
			block.Organization = fields[0][4]
		} else {
			block.Organization = fields[0][2]
		}

		matches := re.FindAllStringSubmatch(addr, -1)
		if len(matches) == 0 {
			continue
		}

		s := matches[0][1]

		if i := strings.IndexByte(s, '/'); i < 0 {
			block.Oui, err = parseMAC(s)
			block.Mask = 24 // len(block.Oui) * 8
		} else {
			var mask int
			block.Oui, err = parseMAC(s[:i])
			mask, err = strconv.Atoi(s[i+1:])
			block.Mask = uint8(mask)
		}

		if err != nil {
			continue
		}

		m.blocks = append(m.blocks, block)

		// create smart map
		for i := len(block.Oui) - 1; i >= 0; i-- {
			_ = block.Oui[i]

		}

		// fmt.Printf("BLA %v %v ALB", m.hw, m.mask)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// New returns a new OUI database loaded from the specified file.
func New(file string) *OuiDB {
	db := &OuiDB{}
	if err := db.load(file); err != nil {
		return nil
	}
	return db
}

func (m *OuiDB) blockLookup(address [6]byte) *addressBlock {
	a := macToUint64(address)
	for _, block := range m.blocks {
		o := macToUint64(block.Oui)
		m := maskToUint64(block.Mask)

		if a&m == o {
			return &block
		}
	}

	return nil
}

// Lookup obtains the vendor organization name from the MAC address s.
func (m *OuiDB) Lookup(s string) (string, error) {
	addr, err := parseMAC(s)
	if err != nil {
		return "", err
	}
	block := m.blockLookup(addr)
	if block == nil {
		return "", ErrInvalidMACAddress
	}
	return block.Organization, nil
}

