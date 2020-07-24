package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

var E = "010001"
var N = "bcb844f8e429e3ed1d8581cab8cf82b0aac8d269023f908f67df488777710bb0426d5d63ec10c30cdf10fc34a7203c8282f43d07626f" +
	"5103bdeb96241cbf848ba7ea1165f6ae908ddc16361ee38a2101753859e7b547c849c51d813ab57d259c0f79f98377a83233925462509d7a" +
	"15dfb17d3e4182f19433f9116b6047b1655881c20ba8cbd277646ce3b221962382ec8e42fb3aff17348f94eedddd600e427fd90dcbe74ec7" +
	"4dfdc85fa5972225c6d97bb2029f0c51d667916660c81d3ae0de5037ef80e33e0f0eb5220be21d8a3bf2ce8f98e784aa27f0f03c9c612a90" +
	"fa10f662687bff503c9df1ed0f7d96fb7ea5a5ec520f2260e126bbad5c76692615a5"
var D = "b56135893161c19ac7e0e519fdfe1351d1132a879a8d9556ff326ef72429165ed5b95f25066225d55d1f6a070109ce9e715664c1902e" +
	"04e35fc9e987d3c98e8edb57f058db7a739ca48704853394329cc018e4effa1f7fb4c72ad065a8c11b409eef508cb698858763808eed842d" +
	"2e90cc79df37ffae480e9bb7ce47bf201492115128d14dc8fbdcc9ba0726cce91e4d985060c5fc18795735fbf4d77dbadf2988e95684723d" +
	"33342a6b324c6b82610daae7fcd933b22e5cd50949ecf85bfc0baae84b66271cba4fbb34c54376d4b64e11ce271e6076ce7cacea8e6bdd7d" +
	"8469b2c497348ad3e862b197f70b12c9ca61f0ef01fff13581e79274c69441ffbd11"

func generateKey(t *testing.T) {
	_, _, _ = E, N, D
	bitSize := 2048

	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	assert.Nilf(t, err, "failed to generate a key")

	fmt.Printf("E: %06x\n", key.E)
	fmt.Printf("N: %s\n", hex.EncodeToString(key.N.Bytes()))
	fmt.Printf("D: %s\n", hex.EncodeToString(key.D.Bytes()))
}

func TestDummy(t *testing.T) {
	if false {
		generateKey(t)
	}
}

func TestDoHuaWeiRSA(t *testing.T) {
	// const rsa = new RSAKey();
	// rsa.setPublic(g_encPublickey.n, g_encPublickey.e);
	// const crypted = rsa.encrypt('test')
	// NOTE: one of the crypted data is the value below
	crypted := "9864fcdd2d34d6743bbf526ca908bc4a8d81fade782db10a4317f44ea7ef8994061421b4d0291205440566e52d8245d813a7c" +
		"7d1c827f994fa3be69bed7185ec40a86db29d32ff8816b4d980c1ba367d24eccec5f9cbe374b7a60b8a723f9842f8db7e5d1b970dbb4" +
		"03d14bf0547ecaa838d525a584b398e4458c2fe284a2e2b2fb94bb33c162a6a5bb26665ed577f8fe5bc2d2af3e18652ae0555b2767e3" +
		"48bc77f56a4a1637faced8e4049756578c24b03b7f1a6ad24cc301af28c0a339b254c38847f928e7afbb8eea9cf307abb3043d220d2a" +
		"1641b186b8ef867595ba5b686603c6fa548733e5a664d2b245cc8225a9343ccc01262ac7d60df0c63344b9a"

	var e = &big.Int{}
	var d = &big.Int{}
	var n = &big.Int{}
	e.SetString(E, 16)
	d.SetString(D, 16)
	n.SetString(N, 16)

	privKey := &rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: n,
			E: int(e.Int64()),
		},
		D: d,
	}
	result := DecryptHuaweiRSA(crypted, privKey)
	assert.Equalf(t, "test", string(result), "it should be able to decrypt messages encoded in web ui")
}

func TestDoHuaWeiRSA_Encrypt(t *testing.T) {
	var e = &big.Int{}
	var d = &big.Int{}
	var n = &big.Int{}
	e.SetString(E, 16)
	d.SetString(D, 16)
	n.SetString(N, 16)

	privKey := &rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: n,
			E: int(e.Int64()),
		},
		D: d,
	}

	encrypted := EncryptHuaweiRSA([]byte("hello"), &privKey.PublicKey)
	encryptedStr := hex.EncodeToString(encrypted)
	fmt.Println(encryptedStr)

	result := DecryptHuaweiRSA(encryptedStr, privKey)
	fmt.Println(string(result))

	//assert.Equalf(t, "test", string(result), "it should be able to decrypt messages encoded in web ui")
}

func TestSomething(t *testing.T) {
	var n = &big.Int{}
	n.SetString(N, 16)
	pubKey := &rsa.PublicKey{
		N: n,
		E: 0x10001,
	}
	input := `<?xml version="1.0" encoding="UTF-8"?><request><Username>admin</Username><Password>AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA</Password><password_type>4</password_type></request>`
	output := EncryptHuaweiRSA([]byte(input), pubKey)
	fmt.Println(hex.EncodeToString(output))
}
