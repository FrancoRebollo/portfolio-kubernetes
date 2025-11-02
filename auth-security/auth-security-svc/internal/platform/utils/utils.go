package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/FrancoRebollo/auth-security-svc/internal/domain"
	jwt "github.com/golang-jwt/jwt"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

func JWTCreate(duration int, credentials domain.Credentials, tokenType string) (string, error) {
	var tokenSeed string
	expiresAt := time.Now().Add(time.Minute * time.Duration(duration)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id_persona":    credentials.IdPersona,
			"api_key":       credentials.ApiKey,
			"canal_digital": credentials.CanalDigital,
			"exp":           expiresAt,
		})

	if tokenType == "ACCESS" {
		tokenSeed = os.Getenv("JWT_ACCESS_SEED")
	}

	if tokenType == "REFRESH" {
		tokenSeed = os.Getenv("JWT_REFRESH_SEED")
	}

	tokenString, err := token.SignedString([]byte(tokenSeed))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CheckJWTAccessToken(tokenJWT string) (*domain.CheckJWT, error) {

	resp := &domain.CheckJWT{
		IdPersona:   0,
		TokenStatus: "token valido",
	}
	claims := jwt.MapClaims{}

	parsedToken, err := jwt.ParseWithClaims(tokenJWT, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			resp.TokenStatus = "firma inválida"
			return resp, nil
		}
		return []byte(os.Getenv("JWT_ACCESS_SEED")), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok && ve.Errors == jwt.ValidationErrorExpired {
			resp.TokenStatus = "token expirado"
			return resp, nil
		} else {
			resp.TokenStatus = err.Error()
			return resp, nil
		}

	}

	claims, bool := parsedToken.Claims.(jwt.MapClaims)

	if !bool {
		resp.TokenStatus = "invalid claims"
		return resp, nil
	}

	expiration, bool := claims["exp"].(float64)

	if !bool {
		resp.TokenStatus = "error verificando token"
		return resp, nil
	}

	expirationTime := time.Unix(int64(expiration), 0)

	if time.Now().After(expirationTime) {
		resp.TokenStatus = "token expirado"
		return resp, nil
	}

	return resp, nil
}

func GetClaimsFromToken(jwtToken string, tokenType string) (jwt.MapClaims, error) {
	var tokenSeed string

	if tokenType == "ACCESS" {
		tokenSeed = os.Getenv("JWT_ACCESS_SEED")
	}

	if tokenType == "REFRESH" {
		tokenSeed = os.Getenv("JWT_REFRESH_SEED")
	}

	claims := jwt.MapClaims{}

	parsedToken, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("firma inválida: %v", token.Header["alg"])
		}
		return []byte(tokenSeed), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok && ve.Errors == jwt.ValidationErrorExpired {
			return nil, fmt.Errorf("token expirado")
		}
		return nil, err
	}

	claims, bool := parsedToken.Claims.(jwt.MapClaims)

	if !bool {
		return nil, fmt.Errorf("invalid claims")
	}

	return claims, nil
}

func GetTokenExpiration(tokenString string, tokenType string) (*time.Time, error) {
	var tokenSeed string

	if tokenType == "ACCESS" {
		tokenSeed = os.Getenv("JWT_ACCESS_SEED")
	}

	if tokenType == "REFRESH" {
		tokenSeed = os.Getenv("JWT_REFRESH_SEED")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSeed), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			expirationTime := time.Unix(int64(exp), 0)
			return &expirationTime, nil
		}
		return nil, fmt.Errorf("token does not contain 'exp' claim")
	}

	return nil, fmt.Errorf("invalid token")
}

func ComparePasswordHash(hashedPassword, password string) error {
	hashBytes := []byte(hashedPassword)
	passwordBytes := []byte(password)

	err := bcrypt.CompareHashAndPassword(hashBytes, passwordBytes)
	if err != nil {
		return fmt.Errorf("invalid password: %v", err)
	}

	return nil
}

func GenerateVerificationCode(length int) (string, error) {
	const charset = "0123456789"
	code := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := range code {
		num, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		code[i] = charset[num.Int64()]
	}
	return string(code), nil
}

func GenerateRandomPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var password string

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password += string(charset[randomIndex.Int64()])
	}

	return password, nil
}

func SendEmail(to string, subject string, body string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	senderEmail := os.Getenv("MAIL_ADDRESS")
	senderPassword := os.Getenv("MAIL_PASSWORD")

	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(smtpHost, smtpPort, senderEmail, senderPassword)
	return d.DialAndSend(m)
}

func HashCredentials(username, password, seed string) (string, error) {
	if len(seed) != 32 {
		return "", errors.New("seed must be 32 characters long")
	}

	h := hmac.New(sha256.New, []byte(seed))
	h.Write([]byte(username + ":" + password))

	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

func GenerateQRCode(username, seed string) (string, string, error) {
	var buf bytes.Buffer

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Thinksoft-autenticacion",
		AccountName: username,
		Secret:      []byte(seed),
	})
	if err != nil {
		return "", "", err
	}

	qrCodeDir := os.Getenv("QR_CODE_PATH")
	if qrCodeDir == "" {
		return "", "", fmt.Errorf("la variable de entorno QR_CODE_PATH no está definida")
	}

	qrFilePath := filepath.Join(qrCodeDir, username+"_qr.png")

	err = qrcode.WriteFile(key.URL(), qrcode.Medium, 256, qrFilePath)
	if err != nil {
		return "", "", err
	}

	qr, err := qrcode.New(key.URL(), qrcode.Medium)
	if err != nil {
		return "", "", err
	}
	err = qr.Write(256, &buf)
	if err != nil {
		return "", "", err
	}

	qrBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	return key.URL(), qrBase64, nil
}

func ValidateCredentialsAndTOTP(totpCode, seed string) (bool, error) {
	// Usa directamente la semilla si ya está en Base32
	valid := totp.Validate(totpCode, seed)

	if !valid {
		return false, errors.New("invalid TOTP code")
	}

	return true, nil
}

func Encrypt(data, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	plaintext := []byte(data)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func EncryptTwo(value, secret string) (string, error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", err
	}

	plainText := []byte(value)

	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plainText))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plainText)

	return base64.RawURLEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(cryptoText, key string) (string, error) {
	fmt.Println("Decrypting" + key)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	fmt.Println(string(ciphertext))

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

func DecryptTwo(value, secret string) (string, error) {
	ciphertext, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return "", fmt.Errorf("decoding base64: %w", err)
	}

	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

func PointerToString(a *string) string {
	if a != nil {
		return *a
	}

	return ""
}
