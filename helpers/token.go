package helpers

import (
	"authentication-authorization/config"
	"context"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`

	jwt.StandardClaims
}

var jwtKey []byte

func SetJWTKey(key string) {
	jwtKey = []byte(key)
}

func GetJWTKey() []byte {
	return []byte(jwtKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
	// Use the dynamically set JWT key here
	secretKey := GetJWTKey() // This retrieves the key set in SetJWTKey in main.go

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func VerifyPassword(foundpwd, pwd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(foundpwd), []byte(pwd))
	return err == nil, err
}

func HashPassword(pwd *string) *string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*pwd), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	hashedPwd := string(bytes)
	return &hashedPwd
}

func GenerateToken(email, userID, role string) (string, string) {
	log.Printf("JWT Key %v Type: %T", jwtKey, jwtKey)

	tokenExpiry := time.Now().Add(5 * time.Minute).Unix()
	refreshTokenExpiry := time.Now().Add(15 * time.Minute).Unix()

	claims := &Claims{
		UserId: userID,
		Email:  email,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiry,
		},
	}
	refreshClaims := &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshTokenExpiry,
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedAccessToken, err := accessToken.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}
	return signedAccessToken, signedRefreshToken
}

func UpdateAllToken(signedToken, signedRefToken, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	userCollection := config.OpenConnection("users")

	updateObj := bson.D{
		{"$set", bson.D{
			{"token", signedToken},
			{"refresh_token", signedRefToken},
			{"updated_at", time.Now()},
		}},
	}

	filter := bson.M{"user_id": userID}
	_, err := userCollection.UpdateOne(ctx, filter, updateObj)
	return err
}
