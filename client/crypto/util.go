package crypto

import (
	"math/big"
	"strconv"
)

func Hex2Int(input string) int {
	result, _ := strconv.ParseInt(input, 16, 31)
	return int(result)
}

func HexToBigInt(input string) *big.Int {
	result := &big.Int{}
	result.SetString(input, 16)
	return result
}
