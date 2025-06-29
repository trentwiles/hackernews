package jwt

import (
	"testing"
	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestGener(t *testing.T) {
	token, err := generateJWT("trent", 60)
	fmt.Println(token)

	assert.Equal(t, err, nil, "Errorless token generation")

	var username string
	username, err = verifyJWT(token)

	assert.Equal(t, username, "trent", "Ensure token username can be parsed")
	assert.Equal(t, err, nil, "Ensure token generation is error free")
}