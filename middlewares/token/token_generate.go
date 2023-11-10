package token

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/satori/go.uuid"
)

type AccessTokenDetail struct {
	AccessToken       string
	AccessTokenUUID   string
	AccessTokenExpire int64

	//below filed's is extra filed that we extract from token
	Username string
	UserID   uint
}

type RefreshTokenDetail struct {
	RefreshToken       string
	RefreshTokenUUID   string
	RefreshTokenExpire int64

	//below filed's is extra filed that we extract from token
	Username string
	UserID   uint
}

type TokenDetails struct {
	at *AccessTokenDetail
	rt *RefreshTokenDetail
}

// CreateToken is function that generate token base on given username and return generated token
func CreateToken(username string, userID uint) (*TokenDetails, error) {

	accessSecretKey := os.Getenv("ACCESS_SECRET_KEY") //you can pass keys and gain better performance bcs every time you are reading from env
	if len(accessSecretKey) == 0 {
		return nil, errors.New("cannot get access secret key from env")
	}

	refreshSecretKey := os.Getenv("REFRESH_SECRET_KEY")
	if len(refreshSecretKey) == 0 {
		return nil, errors.New("cannot get refresh secret key from env")
	}

	//TODO: its better to determine and set access and refresh expire time from env variable

	//creating access
	accessClaims := jwt.MapClaims{}
	accessClaims["username"] = username
	accessClaims["userID"] = userID

	accessUUID := uuid.NewV4().String()
	accessClaims["access_uuid"] = accessUUID

	accessExpire := time.Now().Add(20 * time.Hour).Unix()
	accessClaims["expire_time"] = time.Now().Add(20 * time.Hour).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	signedAccess, _ := at.SignedString([]byte(accessSecretKey))

	accessTokenDetail := &AccessTokenDetail{
		AccessToken:       signedAccess,
		AccessTokenUUID:   accessUUID,
		AccessTokenExpire: accessExpire,
	}

	//creating refresh
	refreshClaims := jwt.MapClaims{}
	refreshClaims["username"] = username
	refreshClaims["userID"] = userID

	refreshUUID := uuid.NewV4().String()
	refreshClaims["refresh_uuid"] = refreshUUID

	refreshExpire := time.Now().Add(2 * time.Hour).Unix()
	refreshClaims["expire_time"] = time.Now().Add(7 * 24 * time.Hour).Unix()

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	signedRefresh, _ := rt.SignedString([]byte(refreshSecretKey))

	refreshTokenDetail := &RefreshTokenDetail{
		RefreshToken:       signedRefresh,
		RefreshTokenUUID:   refreshUUID,
		RefreshTokenExpire: refreshExpire,
	}

	td := &TokenDetails{
		at: accessTokenDetail,
		rt: refreshTokenDetail,
	}

	return td, nil
}

func (t *TokenDetails) GetAccessToken() string {
	return t.at.AccessToken
}

func (t *TokenDetails) GetRefreshToken() string {
	return t.rt.RefreshToken
}

// CreateTokenBasedOnRefreshToken will renew access and refresh token but expire time of refresh token
// remain same because we want after speific period user login in to his account
// Note: this function will not delete previous refresh token uuid from database
// therefore you should delete previous refresh token uuid from database by yourself
func CreateTokenBasedOnRefreshToken(rt *RefreshTokenDetail) (*TokenDetails, error) {

	newToken, err := CreateToken(rt.Username, rt.UserID)

	if err != nil {
		return nil, err
	}

	remainingRefreshTokenExpireTime := time.Unix(rt.RefreshTokenExpire, 0)
	now := time.Now()

	newToken.rt.RefreshTokenExpire = int64(remainingRefreshTokenExpireTime.Sub(now))

	return newToken, nil

}

// ExtractRefreshTokenFrom given refreshToken stirng
func ExtractRefreshTokenFrom(refreshToken string) (*RefreshTokenDetail, error) {

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected singing method. err: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("REFRESH_SECRET_KEY")), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token is expired")
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, fmt.Errorf("token validation faild")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		refreshUUID, ok := claims["refresh_uuid"].(string)

		if !ok {
			return nil, fmt.Errorf("token is invalid because cannot retrieve refresh uuid")
		}

		username, ok := claims["username"].(string)

		if !ok {
			return nil, fmt.Errorf("token is invalid because cannot retrieve username")
		}

		userID, ok := claims["userID"].(float64)

		if !ok {
			return nil, fmt.Errorf("token is invalid because cannot retrieve userID")
		}

		//the reason that we convert expire time to float 64 its kinda weird
		// it seems when we retrive json from browser its convert to float 64
		//because js only support 64 floating points so we should first conver to float64
		//then after that we should convert to int64
		//more information: https://stackoverflow.com/a/29690346
		_, ok = claims["expire_time"].(float64)

		if !ok {
			return nil, fmt.Errorf("token is invalid because cannot retrieve expire time")
		}

		refreshExpireTime := int64(claims["expire_time"].(float64)) //convert float64 to int64

		return &RefreshTokenDetail{
			RefreshToken:       refreshToken,
			RefreshTokenUUID:   refreshUUID,
			RefreshTokenExpire: refreshExpireTime,
			Username:           username,
			UserID:             uint(userID),
		}, nil
	}

	return nil, fmt.Errorf("token is invalid")

}

// //ExtractAccessTokenFrom given accessToken stirng
func ExtractAccessTokenFrom(accessToken string) (*AccessTokenDetail, error) {

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected singing method. err: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("ACCESS_SECRET_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, fmt.Errorf("token validation faild")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)

		if !ok {
			return nil, fmt.Errorf("token is invalid because cannot retrieve refresh uuid")
		}

		username, ok := claims["username"].(string)

		if !ok {
			return nil, fmt.Errorf("token is invalid because cannot retrieve username")
		}

		fmt.Println(claims)

		userID, ok := claims["userID"].(float64)
		if !ok {
			return nil, fmt.Errorf("token is invalid because cannot retrieve userID")
		}

		//the reason that we convert epire time to float 64 its kinda weird
		// it seems when we retrive json from browser its convert to float 64
		//because js only support 64 floating points so we should first conver to float64
		//then after that we should convert to int64
		//more information: https://stackoverflow.com/a/29690346
		_, ok = claims["expire_time"].(float64)

		if !ok {
			return nil, fmt.Errorf("token is invalid because cannot retrieve expire time")
		}

		accessExpireTime := int64(claims["expire_time"].(float64))

		return &AccessTokenDetail{
			AccessToken:       accessToken,
			AccessTokenUUID:   accessUUID,
			AccessTokenExpire: accessExpireTime,
			Username:          username,
			UserID:            uint(userID),
		}, nil
	}

	return nil, fmt.Errorf("token is invalid")

}
