package chia_address_generator

import (
	"crypto/sha256"
	"encoding/hex"
	bls "github.com/chuwt/chia-bls-go"
	"go.uber.org/atomic"
	"strings"
)

func GetAddressFromPuzzleHash(ph []byte, prefix string) (string, error) {
	bits, err := ConvertBits(ph, 8, 5, true)
	if err != nil {
		return "", nil
	}
	return Encode(prefix, bits)
}

func GetPuzzleHashFromAddress(address string) (string, []byte, error) {
	prefix, data, err := Decode(address)
	if err != nil {
		return "", nil, err
	}
	puzzleHash, err := ConvertBits(data, 5, 8, false)
	if err != nil {
		return "", nil, err
	}
	return prefix, puzzleHash, nil
}

func NewAddressFromPkBytes(pkBytes []byte, prefix string) (string, error) {
	pk, err := bls.NewPublicKey(pkBytes)
	if err != nil {
		return "", err
	}

	return NewAddressFromPK(pk, prefix)
}

func NewAddressFromPK(pk bls.PublicKey, prefix string) (string, error) {
	// generate synthetic key
	syntheticKeyBytes := opPointAdd(
		pk,
		opPubKeyForExp(
			opSha256(
				pk.Bytes(),
				defaultHiddenHash,
			),
		),
	)

	bits, err := ConvertBits(genAddress(syntheticKeyBytes), 8, 5, true)
	if err != nil {
		return "", nil
	}
	return Encode(prefix, bits)
}

func NewAddressFromPKHex(pkHex, prefix string) (string, error) {
	if strings.HasPrefix(pkHex, "0x") {
		pkHex = pkHex[2:]
	}
	pkBytes, err := hex.DecodeString(pkHex)
	if err != nil {
		return "", err
	}
	return NewAddressFromPkBytes(pkBytes, prefix)
}

func genAddress(pkBytes []byte) []byte {
	programString := newProgramString(pkBytes)

	var (
		hash      = sha256.New()
		hashStack = NewAddressStack()

		hashBuf = make([]byte, 32)     // 32
		tmpBuf  = make([]byte, 1+2*32) // 1 + 32 + 32

		text    = ""
		firstFF = atomic.NewInt32(0)
	)

	for i := len(programString) - 2; i >= 0; i -= 2 {
		if firstFF.CAS(1, 2) {
			text = programString[i : i+96+2] // 2(len) + 96(hash)
		} else {
			text = programString[i : i+2]
		}
		switch text {
		case "ff":
			p0 := hashStack.Pop()
			p1 := hashStack.Pop()
			copy(tmpBuf[:1], []byte{02})
			copy(tmpBuf[1:1+32], p0)
			copy(tmpBuf[1+32:], p1)
			hash.Write(tmpBuf)
			hashBuf = hash.Sum(nil)
			hash.Reset()
			hashStack.Append(hashBuf)
			if firstFF.CAS(0, 1) {
				i -= 96
			}
		default:
			switch text {
			case "80":
				hash.Write([]byte{01})
			default:
				if len(text) != 2 {
					copy(tmpBuf[:1], []byte{01})
					decodeText, _ := hex.DecodeString(text[2:])
					copy(tmpBuf[1:49], decodeText)
					hash.Write(tmpBuf[:49])
				} else {
					copy(tmpBuf[:1], []byte{01})
					decodeText, _ := hex.DecodeString(text)
					copy(tmpBuf[1:2], decodeText)
					hash.Write(tmpBuf[:2])
				}
			}
			hashBuf = hash.Sum(nil)
			hashStack.Append(hashBuf)
			hash.Reset()
		}
	}
	return hashStack.hashes[0]
}

func newProgramString(pkBytes []byte) string {
	return "" +
		"ff02ffff01ff02ffff01ff02ffff03ff0bffff01ff02ffff03ffff09ff05ffff1dff0bffff1effff0bff0bffff02ff06fff" +
		"f04ff02ffff04ff17ff8080808080808080ffff01ff02ff17ff2f80ffff01ff088080ff0180ffff01ff04ffff04ff04ffff" +
		"04ff05ffff04ffff02ff06ffff04ff02ffff04ff17ff80808080ff80808080ffff02ff17ff2f808080ff0180ffff04ffff0" +
		"1ff32ff02ffff03ffff07ff0580ffff01ff0bffff0102ffff02ff06ffff04ff02ffff04ff09ff80808080ffff02ff06ffff" +
		"04ff02ffff04ff0dff8080808080ffff01ff0bffff0101ff058080ff0180ff018080ffff04ffff01b0" +
		hex.EncodeToString(pkBytes) +
		"ff018080"
}
