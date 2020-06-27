package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/nyks06/backapi"
)

type Tokenizer struct {
	SigningKey string
}

func (t *Tokenizer) Tokenize(args map[string]interface{}) (string, error) {
	params := make(jwt.MapClaims)
	for k, v := range args {
		params[k] = v
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), params)
	tok, err := token.SignedString([]byte(t.SigningKey))
	if err != nil {
		return "", webcore.NewInternalServerError("cannot create a token")
	}

	return tok, nil
}

func (t *Tokenizer) Parse(tokenStr string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.SigningKey), nil
	})

	if err != nil || !token.Valid {
		return nil, webcore.NewInternalServerError("Impossible to parse the given token")
	}

	claims := token.Claims.(jwt.MapClaims)
	return claims, nil
}
