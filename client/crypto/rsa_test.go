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
	// g_encPublickey.n = "My N value"
	// crypted = doRSAEncrypt('test')

	// NOTE: one of the crypted data is the value below
	crypted := "251d3012fdcf03987ff32a30db6f26a2e00be6a0e8a33eb4b3b8067c976895203b28ff7ef586e12253c8654f01711d997c82d" +
		"1aa4969fc8655bd7bdb5a310f7c4cec9db5787867dd76bdc952ab49758ae8b053df4ae20e8c9d107ef587b6173ac34e8b3c6e84b3374" +
		"9fc2ea4a5dabfb39f02ef65e69fdd0a520d22b08368a98fe0a8c7d6d98d466a6dab5245bc7594d2d43a6281200d6379fa6b34c0d384e" +
		"a562c11562e86e4addc4f619c6222df691d6c5f291365fa685739b25abee4f204cac60c16f7f066b368566a0cc03caca94550518c6e2" +
		"8fbe31065ca5405dfa5094148d551364e38bdc7562249d867d60829508fceea13650ebe76cd654af74cba05"

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

func TestEncryptHuaweiRSA(t *testing.T) {
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

	plain := "this is a long text" + N
	encrypted := EncryptHuaweiRSA([]byte(plain), &privKey.PublicKey)
	result := DecryptHuaweiRSA(encrypted, privKey)

	assert.Equalf(t, plain, string(result), "it should be able to decrypt messages with self signed keys")
}
