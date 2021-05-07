// Package go-oui provides functions to work with MAC and OUI's
package ouidb

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
)

var db *OuiDB

func lookup(t *testing.T, mac, org string) {
	if db == nil {
		t.Fatal("database not initialized")
	}
	v, err := db.Lookup(mac)
	if err != nil {
		t.Fatalf("parse: %s: %s", mac, err.Error())
	}
	if v != org {
		t.Fatalf("lookup: input %s, expect %q, got %q", mac, org, v)
	}
	//t.Logf("%s => %s\n", mac, v)
}

func string48(b [6]byte) string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		b[0], b[1], b[2], b[3], b[4], b[5])
}

func string24(b [3]byte) string {
	return fmt.Sprintf("%02x:%02x:%02x:aa:bb:cc", b[0], b[1], b[2])
}

func invalid(t *testing.T, mac string) {
	if db == nil {
		t.Fatal("database not initialized")
	}
	v, err := db.Lookup(mac)
	if err == nil {
		t.Fatalf("didn't fail on invalid %s, got %s", mac, v)
	}
}

func TestInitialization(t *testing.T) {
	var err error
	db, err = NewFromFile("oui.txt")
	if err != nil {
		t.Fatalf("can't load database file oui.txt: %s", err)
	}

	file, err := os.Open("oui.txt")
	if err != nil {
		t.Fatalf("can't open database file oui.txt: %s", err)
	}
	defer file.Close()

	db, err = NewFromReader(file)
	if err != nil {
		t.Fatalf("can't load database file oui.txt: %s", err)
	}

	data, err := ioutil.ReadFile("oui.txt")
	if err != nil {
		t.Fatalf("can't read database file oui.txt: %s", err)
	}

	db, err = New(data)
	if err != nil {
		t.Fatalf("can't load database file oui.txt: %s", err)
	}
}

func TestMissingDBFile(t *testing.T) {
	_, err := NewFromFile("bad-file")
	if err == nil {
		t.Fatal("didn't return err on missing file")
	}
}

func TestInvalidDBFile(t *testing.T) {
	_, err := NewFromFile("ouidb_test.go")
	if err == nil {
		t.Fatal("didn't return err on bad file")
	}
}

func TestLookup24(t *testing.T) {
	lookup(t, "60:03:08:a0:ec:a6", "Apple")
}

func TestLookup36(t *testing.T) {
	lookup(t, "00:1B:C5:00:E1:55", "VigorEle")
}

func TestLookup40(t *testing.T) {
	lookup(t, "20-52-45-43-56-aa", "Receive")
}

func TestLookupUnknown(t *testing.T) {
	lookup(t, "ff:ff:00:a0:ec:a6", "")
}

func TestFormatSingleZero(t *testing.T) {
	lookup(t, "0:25:9c:42:0:62", "Cisco-Li")
}

func TestFormatUppercase(t *testing.T) {
	lookup(t, "0:25:9C:42:C2:62", "Cisco-Li")
}

func TestInvalidMAC1(t *testing.T) {
	invalid(t, "00:25-:9C:42:C2:62")
}

func TestLookupAll48(t *testing.T) {
	for _, b := range db.blocks48 {
		lookup(t, string48(b.oui), b.Organization())
	}
}

func TestLookupAll24(t *testing.T) {
	for _, b := range db.blocks24 {
		lookup(t, string24(b.oui), b.Organization())
	}
}

func BenchmarkAll(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b := db.blocks24[rand.Intn(len(db.blocks24))]
		db.Lookup(string24(b.oui))
	}
}
