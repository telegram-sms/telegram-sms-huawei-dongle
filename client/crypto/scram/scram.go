package scram

import (
	"crypto/rand"
	"fmt"
)

type Scram struct {
	// should be 8
	keySize uint

	nonce []byte
}

func (s *Scram) Init() error {
	if s.keySize == 0 {
		s.keySize = 8
	}

	s.nonce = make([]byte, s.keySize*4)
	_, err := rand.Read(s.nonce)
	if err != nil {
		for i := range s.nonce {
			s.nonce[i] = byte(i)
		}
	}

	return fmt.Errorf("not implemented")
}
