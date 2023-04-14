package jwt_pack

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
)

type JwtConfig struct {
	PrivatePem []byte // 私钥
	PublicPem  []byte // 公钥
}

type JwtCurator struct {
	private *rsa.PrivateKey // 私钥
	public  *rsa.PublicKey  // 公钥
}

func NewJwtCurator(conf *JwtConfig) (*JwtCurator, error) {
	pri, err := jwt.ParseRSAPrivateKeyFromPEM(conf.PrivatePem)
	if err != nil {
		return nil, err
	}
	pub, err := jwt.ParseRSAPublicKeyFromPEM(conf.PublicPem)
	if err != nil {
		return nil, err
	}
	return &JwtCurator{
		private: pri,
		public:  pub,
	}, nil
}

func (c *JwtCurator) Token(claim func() (*jwt.StandardClaims, error)) (string, error) {
	customerClaim, err := claim()
	if err != nil {
		return "", err
	}
	return jwt.NewWithClaims(jwt.SigningMethodES256, customerClaim).SignedString(c.public)
}

func (c *JwtCurator) Parser(token string) (*jwt.StandardClaims, error) {
	claims, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return c.private, nil
	})
	if err != nil {
		return nil, err
	}
	if claims != nil {
		if cli, ok := claims.Claims.(*jwt.StandardClaims); ok {
			err = cli.Valid()
			if err == nil {
				return cli, nil
			}
		}
	}
	return nil, err
}
