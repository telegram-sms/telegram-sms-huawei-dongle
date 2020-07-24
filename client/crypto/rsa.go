package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"log"
	"math/big"
)

func pkcs1Type2(src []byte, blockBytes int) []byte {
	expectedRandom := blockBytes - len(src) - 3
	random := make([]byte, expectedRandom)
	_, _ = rand.Read(random)

	padded := bytes.NewBuffer(nil)
	padded.Grow(blockBytes)
	padded.Write([]byte{0x00, 0x02}) // type 2?
	padded.Write(random)
	padded.WriteByte(0x00)
	padded.Write(src)

	return padded.Bytes()
}

func unPKCS1Type2(src []byte) []byte {
	size := len(src)
	// one of the first 2 bytes needs to be 0x02
	if size < 4 || (src[0] != 0x02 && src[1] != 0x02) {
		// too small or not this type of padding
		return nil
	}

	for i := size - 1; i > 0; i-- {
		if src[i] == 0 {
			return src[i+1 : size]
		}
	}

	return nil
}

func EncryptHuaweiRSA(input []byte, pubKey *rsa.PublicKey) []byte {
	if len(input) == 0 || pubKey == nil {
		log.Fatal("could not do rsa with empty input or empty key")
		return nil
	}

	result := bytes.NewBuffer(nil)
	e := &big.Int{}
	e.SetInt64(int64(pubKey.E))
	encrypted := &big.Int{}
	plain := &big.Int{}

	b64 := []byte(base64.StdEncoding.EncodeToString(input))
	maxSize := len(b64)
	for i := 0; i < maxSize; i += 245 {
		end := i + 245
		if end > maxSize {
			end = maxSize
		}
		// TODO: Derive block size instead of hard code it.
		block := pkcs1Type2(b64[i:end], 256)
		plain.SetBytes(block)
		encrypted.Exp(plain, e, pubKey.N)
		result.Write(encrypted.Bytes())
	}

	return result.Bytes()
}

func DecryptHuaweiRSA(encrypted string, privKey *rsa.PrivateKey) []byte {
	blob, _ := hex.DecodeString(encrypted)

	c := &big.Int{}
	m := &big.Int{}

	c.SetBytes(blob)
	m.Exp(c, privKey.D, privKey.N)

	return unPKCS1Type2(m.Bytes())
}
