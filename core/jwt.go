package core

import (
	"fmt"
	"foodie/server/apierr"
	"strconv"
	"time"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/rs/xid"
)

const (
	// jwtIssuer specifies the jwtIssuer of the JWT tokens.
	jwtIssuer = "foodie"

	// jwtExpiresAfter specifies how long the JWT tokens should last.
	jwtExpiresAfter = time.Hour * 24
)

// IssueJWT issues a jwt with the specified parameters. Returns a
// serialized token and a serialization error, if any.
func IssueJWT(secret []byte, id xid.ID, admin bool, tstamp time.Time) ([]byte, error) {
	token := jws.NewJWT(jws.Claims{}, crypto.SigningMethodHS256)

	// core claims, defined in RFC.
	token.Claims().Set("exp", fmt.Sprintf("%d", tstamp.Add(jwtExpiresAfter).Unix()))
	token.Claims().Set("iat", fmt.Sprintf("%d", tstamp.Unix()))
	token.Claims().Set("iss", jwtIssuer)
	token.Claims().Set("sub", id.String())
	token.Claims().Set("aud", jwtIssuer)

	// custom claim for permissions.
	token.Claims().Set("adm", admin)

	data, err := token.Serialize(secret)
	if err != nil {
		// unlikely to happen.
		return nil, err
	}

	return data, nil
}

// ParseJWT parses a jwt and validates the data.
func ParseJWT(secret, data []byte, tstamp time.Time) (xid.ID, bool, *apierr.Error) {
	token, err := jws.ParseJWT(data)
	if err != nil {
		return xid.NilID(), false, apierr.Unauthorized()
	}

	err = token.Validate(secret, crypto.SigningMethodHS256)
	if err != nil {
		return xid.NilID(), false, apierr.Unauthorized()
	}

	sid, ok := token.Claims().Get("sub").(string)
	if !ok {
		return xid.NilID(), false, apierr.Internal()
	}

	id, err := xid.FromString(sid)
	if err != nil {
		return xid.NilID(), false, apierr.Internal()
	}

	adm, ok := token.Claims().Get("adm").(bool)
	if !ok {
		return xid.NilID(), false, apierr.Internal()
	}

	sexp, ok := token.Claims().Get("exp").(string)
	if !ok {
		return xid.NilID(), false, apierr.Internal()
	}

	exp, err := strconv.Atoi(sexp)
	if err != nil {
		return xid.NilID(), false, apierr.Internal()
	}

	if int64(exp) < tstamp.Unix() {
		return xid.NilID(), false, apierr.Unauthorized()
	}

	return id, adm, nil
}
