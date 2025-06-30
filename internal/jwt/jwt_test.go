package jwt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGener(t *testing.T) {
	token, err := GenerateJWT("trent", 60)
	fmt.Println(token)

	assert.Equal(t, err, nil, "Errorless token generation")

	var username string
	username, err = VerifyJWT(token)

	assert.Equal(t, username, "trent", "Ensure token username can be parsed")
	assert.Equal(t, err, nil, "Ensure token generation is error free")
}

func TestHeaderParsing(t *testing.T) {
	token, _ := GenerateJWT("trent", 60)
	token2, _ := GenerateJWT("tr3nt", 30)

	assert.Equal(t, true, func() bool {
		ok, _ := ParseAuthHeader("Bearer " + token)
		return ok
	}(), "valid auth header (1)")

	assert.Equal(t, false, func() bool {
		ok, _ := ParseAuthHeader("Bearer")
		return ok
	}(), "invalid auth header (2)")

	assert.Equal(t, false, func() bool {
		ok, _ := ParseAuthHeader("Bearer " + token + " def")
		return ok
	}(), "invalid auth header (3)")

	assert.Equal(t, true, func() bool {
		ok, _ := ParseAuthHeader("Bearer " + token2)
		return ok
	}(), "valid auth header (2)")

	assert.Equal(t, false, func() bool {
		ok, _ := ParseAuthHeader("")
		return ok
	}(), "invalid auth header (4)")

	assert.Equal(t, false, func() bool {
		ok, _ := ParseAuthHeader("skfljksf")
		return ok
	}(), "invalid auth header (5)")

}
