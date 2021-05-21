package jwt

import (
	"context"
	"encoding/json"
	"fmt"
	"micro/client/redis"
	"micro/config"
	"micro/model"
	zapLogger "micro/pkg/logger"
	"micro/pkg/token"
	"strings"
	"time"

	jjwt "github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
)

// Generate new jwt token and store into redis DB
func (j *micro) Generate(ctx context.Context, model interface{}) (*model.JWT, error) {
	logger = zapLogger.GetZapLogger(config.Confs.Debug())

	td, err := j.GenerateJWT()
	if err != nil {
		return nil, err
	}

	if err := j.genRefJWT(td); err != nil {
		return nil, err
	}

	if err := j.store(ctx, model, td); err != nil {
		return nil, err
	}

	return td, nil
}

// generate JWT tokens
func (j *micro) GenerateJWT() (*model.JWT, error) {
	// create new jwt
	td := &model.JWT{}
	td.AtExpires = time.Now().Add(time.Duration(config.Confs.Get().Redis.UserDuration)).Unix()
	td.RtExpires = time.Now().Add(time.Duration(config.Confs.Get().Redis.UserDuration)).Unix()
	td.AccessUUID = token.Generate(30)
	td.RefreshUUID = token.Generate(60)

	// New MapClaims for access token
	atClaims := jjwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["uuid"] = td.AccessUUID
	atClaims["exp"] = td.AtExpires
	at := jjwt.NewWithClaims(jjwt.SigningMethodHS256, atClaims)

	var err error
	td.AccessToken, err = at.SignedString([]byte(config.Confs.Get().JWT.Secret))
	if err != nil {
		zapLogger.Prepare(logger).
			Development().
			Level(zap.ErrorLevel).
			Commit(err.Error())
		return nil, err
	}
	return td, nil
}

// generate refresh tokens
func (j *micro) genRefJWT(td *model.JWT) error {
	// New MapClaims for refresh access token
	rtClaims := jjwt.MapClaims{}
	rtClaims["uuid"] = td.RefreshUUID
	rtClaims["exp"] = td.RtExpires
	rt := jjwt.NewWithClaims(jjwt.SigningMethodHS256, rtClaims)

	var err error
	td.RefreshToken, err = rt.SignedString([]byte(config.Confs.Get().JWT.RSecret))
	if err != nil {
		zapLogger.Prepare(logger).
			Development().
			Level(zap.ErrorLevel).
			Commit(err.Error())
		return err
	}
	return nil
}

// store uuid into redis
func (j *micro) store(ctx context.Context, model interface{}, td *model.JWT) error {
	bt, err := json.Marshal(model)
	if err != nil {
		zapLogger.Prepare(logger).
			Development().
			Level(zap.ErrorLevel).
			Commit(err.Error())
		return fmt.Errorf("can not marshal data: %s", model)
	}
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	now := time.Now()

	// make map for store in redis
	if err := redis.Storage.Set(ctx, td.AccessUUID, string(bt), at.Sub(now)); err != nil {
		return err
	}
	return nil
}

// Get jwt token from redis
func (j *micro) Get(ctx context.Context, token string, response interface{}) error {
	if err := redis.Storage.Get(ctx, token, &response); err != nil {
		return err
	}
	return nil
}

// Verify a token
func (j *micro) Verify(tk string) (string, error) {
	strArr := strings.Split(tk, " ")
	if len(strArr) != 2 {
		return "", fmt.Errorf("invalid JWT token")
	}
	token, err := jjwt.Parse(strArr[1], func(token *jjwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jjwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Confs.Get().JWT.Secret), nil
	})
	if err != nil {
		return "", err
	}
	if _, ok := token.Claims.(jjwt.Claims); !ok && !token.Valid {
		return "", err
	}

	claims, ok := token.Claims.(jjwt.MapClaims)
	if ok && token.Valid {
		AccessUUID, ok := claims["uuid"].(string)
		if !ok {
			return "", fmt.Errorf("error in claims uuid from client")
		}

		return AccessUUID, nil
	}
	return "", err
}
