// Package go-oui provides functions to work with MAC and OUI's
package ouidb

import (
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
		t.Fatalf("lookup: input %s, expect %s, got %s", mac, org, v)
	}
	t.Logf("%s => %s\n", mac, v)
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
	db = New("oui.txt")
	if db == nil {
		t.Fatal("can't load database file oui.txt")
	}
}

func TestMissingDBFile(t *testing.T) {
	db := New("bad-file")
	if db != nil {
		t.Fatal("didn't return nil on missing file")
	}
}

func TestInvalidDBFile(t *testing.T) {
	db := New("ouidb_test.go")
	if db != nil {
		t.Fatal("didn't return nil on bad file")
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

func TestFormatSingleZero(t *testing.T) {
	lookup(t, "0:25:9c:42:0:62", "Cisco-Li")
}

func TestFormatUppercase(t *testing.T) {
	lookup(t, "0:25:9C:42:C2:62", "Cisco-Li")
}

func TestInvalidMAC1(t *testing.T) {
	invalid(t, "00:25-:9C:42:C2:62")
}
