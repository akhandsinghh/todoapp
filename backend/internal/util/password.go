package util

import (
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
    "strings"
)

func HashPassword(password string) (string, error) {
    salt := make([]byte, 16)
    if _, err := rand.Read(salt); err != nil { return "", err }
    digest := sha256.Sum256([]byte(base64.RawStdEncoding.EncodeToString(salt) + password))
    return base64.RawStdEncoding.EncodeToString(salt) + ":" + hex.EncodeToString(digest[:]), nil
}

func CheckPassword(password, stored string) bool {
    parts := strings.Split(stored, ":")
    if len(parts) != 2 { return false }
    digest := sha256.Sum256([]byte(parts[0] + password))
    return hex.EncodeToString(digest[:]) == parts[1]
}
