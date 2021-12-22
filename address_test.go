package chia_address_generator

import (
	"encoding/hex"
	bls "github.com/chuwt/chia-bls-go"
	"testing"
)

var (
	testPkBytes, _ = hex.DecodeString("aaf079d607cabb95c0039c51317cd6e84e66bb6b5c9aecf8fdc4f0ba97c7f2ec8bb2b1831ad3ea0ba8f701a26177e43e")

	testAddress = "xch1knrllhacj7j2m7xqt64klys3kfalewr5p94dg9cpxfygpr70secqdnnl9r"
)

func TestGenerateAddress(t *testing.T) {
	pk, err := bls.NewPublicKey(testPkBytes)
	if err != nil {
		t.Fatal(err)
	}

	addr1, err := NewAddressFromPkBytes(testPkBytes, "xch")
	if err != nil {
		t.Fatal(err)
	}

	addr2, err := NewAddressFromPK(pk, "xch")
	if err != nil {
		t.Fatal(err)
	}

	addr3, err := NewAddressFromPKHex(pk.Hex(), "xch")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(testAddress == addr1 && addr1 == addr2 && addr2 == addr3, addr1)
}

func TestPH2Addr(t *testing.T) {
	prefix, puzzleHash, err := GetPuzzleHashFromAddress(testAddress)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("prefix:", prefix)
	t.Log("puzzleHash", hex.EncodeToString(puzzleHash))

	address, err := GetAddressFromPuzzleHash(puzzleHash, prefix)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(address == testAddress, address)
}
