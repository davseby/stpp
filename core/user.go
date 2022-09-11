package core

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/rs/xid"
)

const (
	issuer      = "foodie"
	expireAfter = time.Hour * 24
)

var (
	ErrMalformedToken = errors.New("malformed token")
	ErrExpiredToken   = errors.New("expired token")
)

type User struct {
	ID           xid.ID `json:"id"`
	Name         string `json:"name"`
	PasswordHash []byte `json:"-"`
	Admin        bool   `json:"-"`
}

type UserCore struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (uc *UserCore) Validate() error {
	if uc.Name == "" {
		return errors.New("user name cannot be empty")
	}

	return ValidatePassword(uc.Password)
}

func ValidatePassword(pass string) error {
	if len(pass) < 4 {
		return errors.New("user password must be at least 4 characters long")
	}

	return nil
}

func IssueJWT(secret []byte, id xid.ID, admin bool, tstamp time.Time) ([]byte, error) {
	token := jws.NewJWT(jws.Claims{}, crypto.SigningMethodHS256)

	// Core claims, defined in RFC.
	token.Claims().Set("exp", fmt.Sprintf("%d", tstamp.Add(expireAfter).Unix()))
	token.Claims().Set("iat", fmt.Sprintf("%d", tstamp.Unix()))
	token.Claims().Set("iss", issuer)
	token.Claims().Set("sub", id.String())
	token.Claims().Set("aud", issuer)

	// Custom claim for permissions.
	token.Claims().Set("adm", admin)

	return token.Serialize(secret)
}

func ParseJWT(secret, data []byte, tstamp time.Time) (xid.ID, bool, error) {
	token, err := jws.ParseJWT(data)
	if err != nil {
		return xid.NilID(), false, ErrMalformedToken
	}

	err = token.Validate(secret, crypto.SigningMethodHS256)
	if err != nil {
		return xid.NilID(), false, ErrMalformedToken
	}

	sid, ok := token.Claims().Get("sub").(string)
	if !ok {
		return xid.NilID(), false, ErrMalformedToken
	}

	id, err := xid.FromString(sid)
	if err != nil {
		return xid.NilID(), false, ErrMalformedToken
	}

	adm, ok := token.Claims().Get("adm").(bool)
	if !ok {
		return xid.NilID(), false, ErrMalformedToken
	}

	sexp, ok := token.Claims().Get("exp").(string)
	if !ok {
		return xid.NilID(), false, ErrMalformedToken
	}

	exp, err := strconv.Atoi(sexp)
	if err != nil {
		return xid.NilID(), false, ErrMalformedToken
	}

	if int64(exp) < tstamp.Unix() {
		return xid.NilID(), false, ErrExpiredToken
	}

	return id, adm, nil
}
