# chia-address-generator
> This repo is a hack way to generate an address from `publicKey`.
> So it's not a good enough way to use it in prod, use it just for personal

## Methods
- generate address from pkBytes
```
func NewAddressFromPkBytes(pkBytes []byte, prefix string) (string, error)
```

- generate address from pkHex
```
func NewAddressFromPKHex(pkHex, prefix string) (string, error)
```

- generate address from pk
```
func NewAddressFromPK(pk bls.PublicKey, prefix string) (string, error)
```


## libs
1. use [https://github.com/chuwt/chia-bls-go](https://github.com/chuwt/chia-bls-go) to decode key
2. modify [https://github.com/btcsuite/btcutil/tree/master/bech32](https://github.com/btcsuite/btcutil/tree/master/bech32) to decode bech32

## Buy me coffee
- ETH: `0xdAdf173d0029dfABb64807686b04a1A1Bf6dc79e`
