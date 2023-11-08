package token

import (
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/sod-lol/small-cdn/services/redis"
)

type RefreshAndAccessUUID struct {
	RefreshUUID string
	AccessUUID  string
}

// SaveTokenDetail is function for saving TokenDetail
// TODO:its better design that reduce number of parameter
func SaveTokenDetail(redis *redis.Redis, tokendetail *TokenDetails, username string) (bool, error) {

	//convert to UTC object
	at := time.Unix(tokendetail.at.AccessTokenExpire, 0)
	rt := time.Unix(tokendetail.rt.RefreshTokenExpire, 0)
	now := time.Now()

	redisErr := redis.Set(tokendetail.at.AccessTokenUUID, username, at.Sub(now)).Err()

	if redisErr != nil {
		logrus.Errorf("[Error](SaveTokenDetail) Error accured during set access token in redis. err: %v", redisErr)
		return false, redisErr
	}

	redisErr = redis.Set(tokendetail.rt.RefreshTokenUUID, username, rt.Sub(now)).Err()

	if redisErr != nil {
		logrus.Errorf("[Error](SaveTokenDetail) Error accured during set refresh token in redis. err: %v", redisErr)
		return false, redisErr
	}

	binaryResult, err := json.Marshal(RefreshAndAccessUUID{
		RefreshUUID: tokendetail.rt.RefreshTokenUUID,
		AccessUUID:  tokendetail.at.AccessTokenUUID,
	})

	if err != nil {
		logrus.Errorf("[Error](SaveTokenDetail) cannot marshal RefreshAndAccessUUID. err: %v", redisErr)
		return false, err
	}

	redisErr = redis.Set(username, binaryResult, rt.Sub(now)).Err()

	if redisErr != nil {
		logrus.Errorf("[Error](SaveTokenDetail) Error accured during set user given refresh and access in redis. err: %v", redisErr)
		return false, redisErr
	}

	return true, nil
}

// DeleteAccessToken will delete access token based on given accessTokenUUID
func DeleteAccessToken(redis *redis.Redis, accessTokenUUID string) (bool, error) {

	_, err := redis.Delete(accessTokenUUID).Result()

	if err != nil {
		return false, err
	}

	return true, nil
}

// DeleteAccessToken will delete access token based on given accessTokenUUID
func DeleteRefreshToken(redis *redis.Redis, refreshTokenUUID string) (bool, error) {

	_, err := redis.Delete(refreshTokenUUID).Result()

	if err != nil {
		return false, err
	}

	return true, nil
}

// IsAccessTokenStillValid check if access token still valid
// (Note: this function doesn't return associate value with given key)
func IsAccessTokenStillValid(redis *redis.Redis, accessToken string) (bool, error) {

	isValid, err := redis.Contain(accessToken)

	if isValid {
		return true, nil
	}

	return false, err
}

// GetUsernameByAccessToken return username by given accessToken
// (Note: this function may become expensive in futue realse bcs currently username restore esaily bcs it's present in token
// but in future may username cannot restore esaily therfore this function due backward compality may do some expensive opration)
func GetUsernameByAccessToken(redis *redis.Redis, accessToken string) (string, error) {

	username, err := redis.Get(accessToken).Result()

	if err != nil {
		return "", err
	}

	return username, nil

}

func GetRefreshAndAccessUUIDFrom(redis *redis.Redis, username string) (RefreshAndAccessUUID, error) {

	refreshAndAccessUUID, err := redis.Get(username).Result()

	if err != nil {
		return RefreshAndAccessUUID{}, err
	}

	refreshAndAccessUUIDStruct := RefreshAndAccessUUID{}

	err = json.Unmarshal([]byte(refreshAndAccessUUID), &refreshAndAccessUUIDStruct)

	if err != nil {
		return RefreshAndAccessUUID{}, err
	}

	return refreshAndAccessUUIDStruct, nil
}
