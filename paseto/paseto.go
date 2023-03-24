package paseto

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

type PasetoInterface interface {
	CreateToken(string) (string, error)
	VerifyToken(token string) error
}

type Payload struct {
	UserName     string
	CurrentTime  time.Time
	RandomString string
}

func NewPasetoMaker(key string) PasetoInterface {
	return &PasetoMaker{
		Paseto: paseto.NewV2(),
		Key:    []byte(key),
	}
}

func (maker *PasetoMaker) CreateToken(userName string) (string, error) {
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)

	randomString := base64.URLEncoding.EncodeToString(randomBytes)

	newPayload := &Payload{
		CurrentTime:  time.Now(),
		RandomString: randomString,
		UserName:     userName,
	}

	return maker.Paseto.Encrypt(maker.Key, newPayload, nil)
}

func (maker *PasetoMaker) VerifyToken(token string) error {
	payload := &Payload{}
	err := maker.Paseto.Decrypt(token, maker.Key, payload, nil)
	if err != nil {
		return err
	}

	return nil
}
