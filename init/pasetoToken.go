package init

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/o1egl/paseto"
	"time"
)

type PasetoMaker struct {
	Paseto *paseto.V2
	Key    []byte
}

type payload struct {
	userName     string
	currentTime  time.Time
	randomString string
}

func NewPasetoMaker(key string) *PasetoMaker {
	return &PasetoMaker{
		Paseto: paseto.NewV2(),
		Key:    []byte(key),
	}
}

func (maker *PasetoMaker) CreateToken(userName string) (string, error) {
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)

	randomString := base64.URLEncoding.EncodeToString(randomBytes)

	newPayload := &payload{
		currentTime:  time.Now(),
		randomString: randomString,
		userName:     userName,
	}

	return maker.Paseto.Encrypt(maker.Key, newPayload, nil)
}

func (maker *PasetoMaker) VerifyToken(token, userName string) (bool, error) {
	payload := &payload{}

	err := maker.Paseto.Decrypt(token, maker.Key, payload, nil)

	if err != nil {
		return false, err
	}

	if payload.userName == userName {
		return true, nil
	}

	return false, err
}
