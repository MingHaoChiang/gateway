package util

import (
	"log"
	"time"

	"github.com/SermoDigital/jose/jws"
)

//JwtConf ...
type JwtConf struct {
	Method string `json:"method"`
	Key    string `json:"key"`
	Issuer string `json:"issuer"`
	Expire int64  `json:"expire"`
}

//GetJWT ...
func GetJWT(conf JwtConf, data map[string]interface{}) (token string, err error) {
	payload := jws.Claims{}
	for k, v := range data {
		payload.Set(k, v)
	}
	now := time.Now()
	payload.SetIssuer(conf.Issuer)
	payload.SetIssuedAt(now)
	payload.SetExpiration(now.Add(time.Duration(conf.Expire) * time.Minute))
	jwtObj := jws.NewJWT(payload, jws.GetSigningMethod(conf.Method))
	tokenBytes, err := jwtObj.Serialize([]byte(conf.Key))
	if err != nil {
		log.Println("Serialize err:", conf.Key)
		return
	}
	token = string(tokenBytes)
	return
}

//VerifyJWT ...
func VerifyJWT(conf JwtConf, token string) (ret bool, err error) {
	jwtObj, err := jws.ParseJWT([]byte(token))
	if err != nil {
		log.Println("ParseJWT err:", err, token)
		return
	}
	err = jwtObj.Validate([]byte(conf.Key), jws.GetSigningMethod(conf.Method))
	if err == nil {
		ret = true
		return
	}
	log.Println("Validate err:", err, conf, token)
	return
}

//RenewJWT ...
func RenewJWT(conf JwtConf, token string) (newToken string, err error) {
	isValid, err := VerifyJWT(conf, token)
	if !isValid {
		log.Println("VerifyJWT err:", err, token)
		return
	}
	jwtObj, err := jws.ParseJWT([]byte(token))
	if err != nil {
		log.Println("ParseJWT err:", err)
		return
	}
	payload := jws.Claims{}
	for k, v := range jwtObj.Claims() {
		payload.Set(k, v)
	}
	now := time.Now()
	payload.SetIssuedAt(now)
	payload.SetExpiration(now.Add(time.Duration(conf.Expire) * time.Minute))
	jwtObj = jws.NewJWT(payload, jws.GetSigningMethod(conf.Method))
	tokenBytes, err := jwtObj.Serialize([]byte(conf.Key))
	if err != nil {
		log.Println("Serialize err:", err)
		return
	}
	newToken = string(tokenBytes)
	return
}
