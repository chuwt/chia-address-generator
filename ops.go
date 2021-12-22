package chia_address_generator

import (
	"crypto/sha256"
	bls "github.com/chuwt/chia-bls-go"
	"math/big"
)

const (
	defaultHiddenHashStr = "711d6c4e32c92e53179b199484cf8c897542bc57f2b22582799f9d657eec4699"
)

var (
	defaultHiddenHash = []byte{
		113, 29, 108, 78, 50, 201, 46, 83, 23, 155, 25, 148, 132, 207, 140, 137, 117,
		66, 188, 87, 242, 178, 37, 130, 121, 159, 157, 101, 126, 236, 70, 153,
	}

	defaultArg = []byte{
		115, 237, 167, 83, 41, 157, 125, 72, 51, 57, 216, 8, 9, 161, 216, 5,
		83, 189, 164, 2, 255, 254, 91, 254, 255, 255, 255, 255, 0, 0, 0, 1,
	}

	bigArg = new(big.Int).SetBytes(defaultArg)
)

func opSha256(pkBytes, hiddenHash []byte) []byte {
	h := sha256.New()
	h.Write(pkBytes)
	h.Write(hiddenHash)
	return h.Sum(nil)
}

func opPubKeyForExp(arg []byte) bls.PublicKey {
	i0 := bigIntFromBytes(arg)
	i0.Mod(i0, bigArg)
	return bls.KeyFromBytes(i0.Bytes()).GetPublicKey()
}

func opPointAdd(pk, exPk bls.PublicKey) []byte {
	return pk.Add(exPk).Bytes()
}

func bigIntFromBytes(buf []byte) *big.Int {
	var x = new(big.Int)
	if len(buf) == 0 {
		return x.SetBytes(buf)
	}

	if (0x80 & buf[0]) == 0 { // positive number
		return x.SetBytes(buf)
	}

	for i := range buf {
		buf[i] = ^buf[i]
	}

	return x.SetBytes(buf).Add(x, big.NewInt(1)).Neg(x)
}