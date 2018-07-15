// Author: catlittlechen@gmail.com
package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"time"
)

func toBytes(value int64) []byte {
	result := make([]byte, 0, 8)
	shift := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, s := range shift {
		result = append(result, byte((value>>s)&0xFF))
	}
	return result
}

func password(secret string, value []byte) uint32 {

	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	hmacSha1 := hmac.New(sha1.New, key)
	hmacSha1.Write(value)
	hash := hmacSha1.Sum(nil)

	offset := hash[len(hash)-1] & 0x0F
	number := binary.BigEndian.Uint32(hash[offset : offset+4])
	number &= 0x7fffffff
	password := number % 1000000

	return password
}

var secret = flag.String("s", "secret", "secret key to use")

func main() {
	flag.Parse()

	second := time.Now().Unix()
	pwd := password(*secret, toBytes(second/30))
	fmt.Printf("%06d", pwd)
}
