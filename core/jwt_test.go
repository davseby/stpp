package core

import (
	"foodie/server/apierr"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_IssueJWT(t *testing.T) {
	tstamp := time.Date(2000, time.April, 15, 10, 0, 0, 0, time.UTC)

	id, err := xid.FromString("73s3r876i1e72n4h3d40")
	require.NoError(t, err)

	sjwt, err := IssueJWT([]byte{1, 2, 3}, id, true, tstamp.Add(time.Minute))
	require.NoError(t, err)

	assert.Equal(
		t,
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImZvb2RpZSIsImV4cCI6Ijk1NTc5NjQ2MCIsImlhdCI6Ijk1NTc5Mjg2MCIsImlzcyI6ImZvb2RpZSIsInN1YiI6IjczczNyODc2aTFlNzJuNGgzZDQwIn0.0O2UJyLoKD3q-3DcbWfBvRhCY8OF00Ko462JUUHiaNU",
		string(sjwt),
	)
}

func Test_ParseJWT(t *testing.T) {
	tstamp := time.Date(2000, time.April, 15, 10, 0, 0, 0, time.UTC)
	secret := []byte{1, 2, 3}

	id, err := xid.FromString("73s3r876i1e72n4h3d40")
	require.NoError(t, err)

	tests := map[string]struct {
		Data   []byte
		ID     xid.ID
		Admin  bool
		Tstamp time.Time
		Error  *apierr.Error
	}{
		"Invalid JWT": {
			Data:  []byte("UUHiaNU"),
			Error: apierr.Unauthorized(),
		},
		"Invalid secret": {
			Data:  []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImZvb2RpZSIsImV4cCI6Ijk1NTc5NjQ2MCIsImlhdCI6Ijk1NTc5Mjg2MCIsImlzcyI6ImZvb2RpZSIsInN1YiI6IjczczNyODc2aTFlNzJuNGgzZDQwIn0.eJQz7KkyHVTqXFI5AYxSWqlXjyKod_QEALV6IklsLko"),
			Error: apierr.Unauthorized(),
		},
		"Missing user id": {
			Data:  []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImZvb2RpZSIsImV4cCI6Ijk1NTc5NjQ2MCIsImlhdCI6Ijk1NTc5Mjg2MCIsImlzcyI6ImZvb2RpZSJ9.CQfQcBIRhLBgitNAO_pNgsj5_B1yjEQIw8MXntcltLY"),
			Error: apierr.Internal(),
		},
		"Invalid user id format": {
			Data:  []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImZvb2RpZSIsImV4cCI6Ijk1NTc5NjQ2MCIsImlhdCI6Ijk1NTc5Mjg2MCIsImlzcyI6ImZvb2RpZSIsInN1YiI6IjMzMyJ9.aRsCCtwPS876E_ZPp_6q-Le-BXET3UCtPbWrlooPatA"),
			Error: apierr.Internal(),
		},
		"Missing admin flag": {
			Data:  []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJmb29kaWUiLCJleHAiOiI5NTU3OTY0NjAiLCJpYXQiOiI5NTU3OTI4NjAiLCJpc3MiOiJmb29kaWUiLCJzdWIiOiI3M3Mzcjg3NmkxZTcybjRoM2Q0MCJ9.UD3GEnD08vLySHKV1r1JDfnOytJd2OI5PxC8UOtudEc"),
			Error: apierr.Internal(),
		},
		"Missing expiration time": {
			Data:  []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImZvb2RpZSIsImlhdCI6Ijk1NTc5Mjg2MCIsImlzcyI6ImZvb2RpZSIsInN1YiI6IjczczNyODc2aTFlNzJuNGgzZDQwIn0.oUfq0jz3ukjoWKs5kRsIgnNOmpxDk2UMgB--EizjegA"),
			Error: apierr.Internal(),
		},
		"Invalid expiration time format": {
			Data:  []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImZvb2RpZSIsImV4cCI6Imtld2tlcSIsImlhdCI6Ijk1NTc5Mjg2MCIsImlzcyI6ImZvb2RpZSIsInN1YiI6IjczczNyODc2aTFlNzJuNGgzZDQwIn0.9DdHFlZ4m8sA1L4bL2Sj4dZQPdVBKqpH90ZbOevlS8I"),
			Error: apierr.Internal(),
		},
		"Expired token": {
			Data:   []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImZvb2RpZSIsImV4cCI6Ijk1NTc5NjQ2MCIsImlhdCI6Ijk1NTc5Mjg2MCIsImlzcyI6ImZvb2RpZSIsInN1YiI6IjczczNyODc2aTFlNzJuNGgzZDQwIn0.0O2UJyLoKD3q-3DcbWfBvRhCY8OF00Ko462JUUHiaNU"),
			Tstamp: tstamp.Add(2 * time.Hour),
			Error:  apierr.Unauthorized(),
		},
		"Successfully parsed JWT": {
			Data:   []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImZvb2RpZSIsImV4cCI6Ijk1NTc5NjQ2MCIsImlhdCI6Ijk1NTc5Mjg2MCIsImlzcyI6ImZvb2RpZSIsInN1YiI6IjczczNyODc2aTFlNzJuNGgzZDQwIn0.0O2UJyLoKD3q-3DcbWfBvRhCY8OF00Ko462JUUHiaNU"),
			ID:     id,
			Admin:  true,
			Tstamp: tstamp,
		},
	}

	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			id, admin, err := ParseJWT(secret, test.Data, test.Tstamp)
			if test.Error != nil {
				assert.Equal(t, xid.NilID(), id)
				assert.False(t, admin)
				assert.Equal(t, test.Error, err)
				return
			}

			require.Nil(t, err)
			assert.Equal(t, test.ID, id)
			assert.Equal(t, test.Admin, admin)
		})
	}
}
