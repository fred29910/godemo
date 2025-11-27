package main

import (
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
)

type base58Encoder struct{}

func (enc base58Encoder) Encode(u uuid.UUID) string {
	return base58.Encode(u[:])
}

func (enc base58Encoder) Decode(s string) (uuid.UUID, error) {
	return uuid.FromBytes(base58.Decode(s))
}

func main() {
	enc := base58Encoder{}
	fmt.Println(shortuuid.NewWithEncoder(enc)) // 6R7VqaQHbzC1xwA5UueGe6
}
