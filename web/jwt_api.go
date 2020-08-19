package web

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"com.hyweb/gateway/util"
)

//jwtTokenRequest ...
type jwtTokenRequest struct {
	JwtConf  util.JwtConf `json:"jwt_conf"`
	JwtToken string       `json:"jwt_token,omitempty"`
	JwtClaim string       `json:"jwt_claim,omitempty"`
}

// GetJwtToken ...
func GetJwtToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Pragma", "no-cache")
	var requestBody jwtTokenRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Get request body error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &requestBody); err != nil {
		log.Println("parse request body error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	data := make(map[string]interface{})
	if err := json.Unmarshal([]byte(requestBody.JwtClaim), &data); requestBody.JwtClaim != "" && err != nil {
		log.Println("Unmarshal error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	jwtToken, err := util.GetJWT(requestBody.JwtConf, data)
	if err != nil {
		log.Println("GetJWT error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jwtToken))
}

// VerifyJWTToken ...
func VerifyJWTToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Pragma", "no-cache")
	var requestBody jwtTokenRequest

	body, err := base64.RawURLEncoding.DecodeString(r.Header.Get("Authorization")[len("Bearer "):])
	if err != nil {
		log.Println("DecodeString error", err, r.Header.Get("Authorization")[len("Bearer "):])
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &requestBody); err != nil {
		log.Println("parse request body error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	isValid, err := util.VerifyJWT(requestBody.JwtConf, requestBody.JwtToken)
	if err != nil {
		log.Println("VerifyJWT error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isValid {
		w.WriteHeader(http.StatusForbidden)
	}

	w.WriteHeader(http.StatusOK)
}

// RenewJWTToken ...
func RenewJWTToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Pragma", "no-cache")
	var requestBody jwtTokenRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Get request body error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &requestBody); err != nil {
		log.Println("parse request body error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	data := make(map[string]interface{})
	if err := json.Unmarshal([]byte(requestBody.JwtClaim), &data); requestBody.JwtClaim != "" && err != nil {
		log.Println("Unmarshal error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	jwtToken, err := util.RenewJWT(requestBody.JwtConf, requestBody.JwtToken)
	if err != nil {
		log.Println("RenewJWT error:", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jwtToken))
}
