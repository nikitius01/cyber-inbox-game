package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"cybersecurity-game/backend/internal/users"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Service struct {
	users     users.Repository
	jwtSecret []byte
}

func NewService(repo users.Repository, secret string) *Service {
	return &Service{users: repo, jwtSecret: []byte(secret)}
}

func (s *Service) Register(username, email, password string) (users.User, string, error) {
	if !validText(username, 3, 32) || !validText(email, 5, 120) || !validText(password, 8, 128) || !strings.Contains(email, "@") {
		return users.User{}, "", errors.New("invalid registration data")
	}
	user := users.User{
		ID:           randomID("usr"),
		Username:     strings.TrimSpace(username),
		Email:        strings.ToLower(strings.TrimSpace(email)),
		PasswordHash: hashPassword(password),
		Stats: users.Stats{
			SolvedByCategory: map[string]int{"easy": 0, "medium": 0, "hard": 0, "nightmare": 0, "AI": 0},
		},
		CreatedAt: time.Now().UTC(),
	}
	if err := s.users.Create(user); err != nil {
		return users.User{}, "", err
	}
	token, err := s.IssueToken(user.ID)
	return user, token, err
}

func (s *Service) Login(email, password string) (users.User, string, error) {
	user, err := s.users.FindByEmail(strings.ToLower(strings.TrimSpace(email)))
	if err != nil || !verifyPassword(password, user.PasswordHash) {
		return users.User{}, "", ErrInvalidCredentials
	}
	token, err := s.IssueToken(user.ID)
	return user, token, err
}

func (s *Service) IssueToken(userID string) (string, error) {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	claims, _ := json.Marshal(map[string]any{"sub": userID, "exp": time.Now().Add(24 * time.Hour).Unix()})
	body := base64.RawURLEncoding.EncodeToString(claims)
	signature := s.sign(header + "." + body)
	return header + "." + body + "." + signature, nil
}

func (s *Service) ValidateToken(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 || !hmac.Equal([]byte(s.sign(parts[0]+"."+parts[1])), []byte(parts[2])) {
		return "", errors.New("bad token")
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}
	var claims struct {
		Sub string `json:"sub"`
		Exp int64  `json:"exp"`
	}
	if err := json.Unmarshal(raw, &claims); err != nil {
		return "", err
	}
	if claims.Sub == "" || time.Now().Unix() > claims.Exp {
		return "", errors.New("expired token")
	}
	return claims.Sub, nil
}

func (s *Service) sign(value string) string {
	mac := hmac.New(sha256.New, s.jwtSecret)
	mac.Write([]byte(value))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func validText(value string, minLen, maxLen int) bool {
	value = strings.TrimSpace(value)
	return len(value) >= minLen && len(value) <= maxLen
}

func randomID(prefix string) string {
	buf := make([]byte, 12)
	_, _ = rand.Read(buf)
	return fmt.Sprintf("%s_%s", prefix, hex.EncodeToString(buf))
}

func hashPassword(password string) string {
	salt := make([]byte, 16)
	_, _ = rand.Read(salt)
	hash := derivePasswordKey([]byte(password), salt, 120000, 32)
	return hex.EncodeToString(salt) + ":" + hex.EncodeToString(hash)
}

func verifyPassword(password, stored string) bool {
	parts := strings.Split(stored, ":")
	if len(parts) != 2 {
		return false
	}
	salt, err := hex.DecodeString(parts[0])
	if err != nil {
		return false
	}
	expected, err := hex.DecodeString(parts[1])
	if err != nil {
		return false
	}
	actual := derivePasswordKey([]byte(password), salt, 120000, len(expected))
	return subtle.ConstantTimeCompare(actual, expected) == 1
}

func derivePasswordKey(password, salt []byte, iterations, keyLen int) []byte {
	hashLen := 32
	blocks := (keyLen + hashLen - 1) / hashLen
	out := make([]byte, 0, blocks*hashLen)
	for block := 1; block <= blocks; block++ {
		mac := hmac.New(sha256.New, password)
		mac.Write(salt)
		mac.Write([]byte{byte(block >> 24), byte(block >> 16), byte(block >> 8), byte(block)})
		u := mac.Sum(nil)
		t := append([]byte(nil), u...)
		for i := 1; i < iterations; i++ {
			mac = hmac.New(sha256.New, password)
			mac.Write(u)
			u = mac.Sum(nil)
			for j := range t {
				t[j] ^= u[j]
			}
		}
		out = append(out, t...)
	}
	return out[:keyLen]
}
