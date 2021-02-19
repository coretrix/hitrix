package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"strconv"
	"strings"
	"time"
)

type JWT struct {
}

func (t *JWT) EncodeJWT(secret string, headers, payload map[string]string) (string, error) {
	algo, ok := headers["algo"]

	if !ok {
		return "", fmt.Errorf("cannot create JWT")
	}

	hashData := t.createHash(secret, algo)

	head, err := json.Marshal(headers)

	if err != nil {
		return "", err
	}

	payl, err := json.Marshal(payload)

	if err != nil {
		return "", err
	}

	h := base64.URLEncoding.EncodeToString(head)
	pl := base64.URLEncoding.EncodeToString(payl)

	token := fmt.Sprintf("%s.%s", h, pl)

	_, err = hashData.Write([]byte(token))

	if err != nil {
		return "", err
	}

	newHash := hashData.Sum(nil)
	hs := base64.URLEncoding.EncodeToString(newHash)

	return fmt.Sprintf("%s.%s", token, hs), nil
}

func (t *JWT) VerifyJWT(secret, jwt string, expire int64) error {
	jwtToken := strings.Split(jwt, ".")

	if len(jwtToken) != 3 {
		return fmt.Errorf("token not valid need to be from three parts")
	}

	header := make(map[string]string)

	h, err := base64.URLEncoding.DecodeString(jwtToken[0])

	if err != nil {
		return err
	}

	err = json.Unmarshal(h, &header)

	if err != nil {
		return err
	}

	algo, ok := header["algo"]

	if !ok {
		return fmt.Errorf("missing algo in header part")
	}

	jwtType, ok := header["type"]

	if !ok || jwtType != "JWT" {
		return fmt.Errorf("missing type in header part")
	}

	mhash := t.createHash(secret, algo)

	_, err = mhash.Write([]byte(fmt.Sprintf("%s.%s", jwtToken[0], jwtToken[1])))

	if err != nil {
		return err
	}

	jwtTokenValue, _ := base64.URLEncoding.DecodeString(jwtToken[2])
	valid := hmac.Equal(mhash.Sum(nil), jwtTokenValue)

	if !valid {
		return fmt.Errorf("token not valid")
	}

	return t.checkTime(jwtToken[1], expire)
}

func (t *JWT) checkTime(exp string, expire int64) error {
	payload := make(map[string]string)

	payl, err := base64.URLEncoding.DecodeString(exp)

	if err != nil {
		return err
	}

	err = json.Unmarshal(payl, &payload)
	if err != nil {
		return err
	}

	expireTime, ok := payload["exp"]

	if !ok {
		return fmt.Errorf("token expire time not valid")
	}

	expTime, _ := strconv.ParseInt(expireTime, 10, 64)
	expTime = time.Now().Unix() - expTime

	valid := expTime < expire
	if !valid {
		return fmt.Errorf("token time not valid %d %d", expTime, expire)
	}

	return nil
}

func (t *JWT) createHash(secret, algo string) hash.Hash {
	switch algo {
	case "HS256":
		return hmac.New(sha256.New, []byte(secret))
	default:
		return hmac.New(sha256.New, []byte(secret))
	}
}
