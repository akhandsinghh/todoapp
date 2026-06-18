package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Claims struct {
	UserID int64  `json:"sub"`
	Email  string `json:"email"`
	Exp    int64  `json:"exp"`
}

func SignToken(userID int64, email, secret string) (string, error) {
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	claims := Claims{UserID: userID, Email: email, Exp: time.Now().Add(24 * time.Hour).Unix()}
	h, _ := json.Marshal(header)
	c, _ := json.Marshal(claims)
	unsigned := enc(h) + "." + enc(c)
	sig := sign(unsigned, secret)
	return unsigned + "." + sig, nil
}

func VerifyToken(token, secret string) (Claims, error) {
	var claims Claims
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return claims, errors.New("invalid token")
	}
	unsigned := parts[0] + "." + parts[1]
	if !hmac.Equal([]byte(sign(unsigned, secret)), []byte(parts[2])) {
		return claims, errors.New("invalid signature")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return claims, err
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return claims, err
	}
	if claims.Exp < time.Now().Unix() {
		return claims, errors.New("token expired")
	}
	if claims.UserID == 0 {
		if id, err := strconv.ParseInt(fmt.Sprint(claims.UserID), 10, 64); err == nil {
			claims.UserID = id
		}
	}
	return claims, nil
}

func enc(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

func sign(data, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	return enc(mac.Sum(nil))
}
