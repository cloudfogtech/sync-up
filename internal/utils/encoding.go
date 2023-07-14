package utils

import (
	"crypto/sha256"
	"fmt"
	uuid "github.com/satori/go.uuid"
)

func UUID() string {
	return uuid.Must(uuid.NewV4(), nil).String()
}

func Sha256(str string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(str)))
}

func EncodePasswordFE(pass, feSalt string) string {
	return Sha256(feSalt + Sha256(feSalt+pass+feSalt) + feSalt)
}

func EncodePasswordBE(pass, beSalt string) string {
	return Sha256(Sha256(beSalt+pass) + beSalt)
}
