// Package go-oui provides functions to work with MAC and OUI's
package ouidb

import (
	"bufio"
	"encoding/hex"
	"errors"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// https://code.wireshark.org/review/gitweb?p=wireshark.git;a=blob_plain;f=manuf
// Bigger than we need, not too big to worry about overflow
const big = 0xFFFFFF

var ErrInvalidMACAddress = errors.New("invalid MAC address")

type HardwareAddr net.HardwareAddr

func parseMAC(s string) ([6]byte, error) {
	var hw [6]byte

	oct := strings.FieldsFunc(s, func(r rune) bool { return r == ':' || r == '-' })

	_, err := hex.Decode(hw[:], []byte(strings.Join(oct, "")))
	if err != nil {
		return hw, err
	}

	return hw, nil
}

// Mask returns the result of masking the address with mask.
func (address HardwareAddr) Mask(mask []byte) []byte {
	n := len(address)
	if n != len(mask) {
		return nil
	}
	out := make([]byte, n)
	for i := 0; i < n; i++ {
		out[i] = address[i] & mask[i]
	}
	return out
}

type t2 struct {
	T3    map[byte]t2
	Block *AddressBlock
}

type OuiDb struct {
	hw   [6]byte
	mask int

	Blocks []AddressBlock

	t map[int]t2
}

// New returns a new OUI database loaded from the specified file.
func New(file string) *OuiDb {
	db := &OuiDb{}
	if err := db.Load(file); err != nil {
		return nil
	}
	return db
}

// Lookup finds the OUI the address belongs to
func (m *OuiDb) lookup(address [6]byte) *AddressBlock {
	a := macToUint64(address)
	for _, block := range m.Blocks {
		o := macToUint64(block.Oui)
		m := maskToUint64(block.Mask)

		if a &m == o {
			return &block
		}
	}

	return nil
}

// VendorLookup obtains the vendor organization name from the MAC address s.
func (m *OuiDb) VendorLookup(s string) (string, error) {
	addr, err := parseMAC(s)
	if err != nil {
		return "", err
	}
	block := m.lookup(addr)
	if block == nil {
		return "", ErrInvalidMACAddress
	}
	return block.Organization, nil
}

func byteIndex(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}

func (m *OuiDb) Load(path string) error {
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

		block := AddressBlock{}

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

		if i := byteIndex(s, '/'); i < 0 {
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

		//fmt.Println("OUI:", block.Oui, block.Mask, err)

		m.Blocks = append(m.Blocks, block)

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

func CIDRMask(ones, bits int) []byte {
	l := bits / 8
	m := make([]byte, l)

	n := uint(ones)
	for i := 0; i < l; i++ {
		if n >= 8 {
			m[i] = 0xff
			n -= 8
			continue
		}
		m[i] = ^byte(0xff >> n)
		n = 0
	}

	return (m)
}

// oui, mask, organization
type AddressBlock struct {
	Oui          [6]uint8
	Mask         uint8
	Organization string
}

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

