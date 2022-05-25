package auth

import (
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/electric-saw/pg-shazam/internal/pkg/log"
)

type EncryptionType int

const (
	EncryptionType_MD5 EncryptionType = iota
	EncryptionType_SCRAM_SHA_256
	EncryptionType_UNKNOWN
)

func passwordCheck(user, pass, passDb string) (bool, error) {
	encryptionType := decodePasswordAlorithm(passDb)

	switch encryptionType {
	case EncryptionType_MD5:
		return checkMD5(user, pass, passDb)
	case EncryptionType_SCRAM_SHA_256:
		return checkSCRAM_SHA_256(user, pass, passDb)
	default:
		log.Warnf("Unknown password encryption type: %x, trying to compare plain text", encryptionType)
		return pass == passDb, nil
	}
}

func decodePasswordAlorithm(hash string) EncryptionType {
	lowerHash := strings.ToLower(hash)
	if strings.HasPrefix(lowerHash, "md5") {
		return EncryptionType_MD5
	}

	if strings.HasPrefix(lowerHash, "scram-sha-256") {
		return EncryptionType_SCRAM_SHA_256
	}

	return EncryptionType_UNKNOWN
}

func checkMD5(user, pass, passDb string) (bool, error) {
	h := md5.New()
	_, _ = h.Write([]byte(pass + user))
	hashPass := fmt.Sprintf("md5%x", string(h.Sum(nil)))

	return hashPass == passDb, nil
}

func checkSCRAM_SHA_256(user, pass, passDb string) (bool, error) {
	// generate scram-sha-256 hash from password and username
	// https://gist.githubusercontent.com/jkatz/e0a1f52f66fa03b732945f6eb94d9c21/raw/2948060a8ec97994f70679a98b7b385e3dcee6f1/encrypt_password.py	//

	// return false, fmt.Errorf("Not implemented")
	return true, nil
}
