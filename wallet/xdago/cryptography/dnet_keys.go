package cryptography

type DnetKeys struct {
	Prv          []byte
	Pub          []byte
	Sect0Encoded []byte
	Sect0        []byte
}

func NewDnetKeys() *DnetKeys {
	return &DnetKeys{
		Prv:          make([]byte, 1024),
		Pub:          make([]byte, 1024),
		Sect0Encoded: make([]byte, 512),
		Sect0:        make([]byte, 512),
	}
}
