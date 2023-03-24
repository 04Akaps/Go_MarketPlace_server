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
	currentTime  time.Time
	randomString string
}

func NewPasetoMaker(key string) *PasetoMaker {
	return &PasetoMaker{
		Paseto: paseto.NewV2(),
		Key:    []byte(key),
	}
}

func (maker *PasetoMaker) CreateToken() (string, error) {
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)

	randomString := base64.URLEncoding.EncodeToString(randomBytes)

	newPayload := &payload{
		currentTime:  time.Now(),
		randomString: randomString,
	}

	return maker.Paseto.Encrypt(maker.Key, newPayload, nil)
}
