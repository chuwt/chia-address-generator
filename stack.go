package chia_address_generator

type AddressStack struct {
	hashes [][]byte
	step int32
}

// NewAddressStack use only for address generator, so the len(hashes) is 14
func NewAddressStack() *AddressStack {
	return &AddressStack{
		hashes: make([][]byte, 14),
		step:   -1,
	}
}

func (stack *AddressStack) Append(hash []byte) {
	var newHash = make([]byte, 32)
	copy(newHash, hash)
	stack.step += 1
	stack.hashes[stack.step] = newHash
}

func (stack *AddressStack) Pop() []byte {
	if stack.step == -1 {
		return nil
	}
	hash := stack.hashes[stack.step]
	stack.step-=1
	return hash
}
