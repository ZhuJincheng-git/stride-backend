package jwt_test

import (
	"strings"
	"testing"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ZhuJincheng-git/stride-backend/pkg/jwt"
)

func TestGenerateAndParseRoundTrip(t *testing.T) {
	mgr := jwt.New("super-secret", time.Hour, "stride-backend")
	uid := uuid.New()

	tok, err := mgr.Generate(uid)
	require.NoError(t, err)
	require.NotEmpty(t, tok)

	claims, err := mgr.Parse(tok)
	require.NoError(t, err)
	require.Equal(t, uid, claims.UserID)
	require.Equal(t, "stride-backend", claims.Issuer)
	require.Equal(t, uid.String(), claims.Subject)
}

func TestParseRejectsTamperedSignature(t *testing.T) {
	mgr := jwt.New("secret", time.Hour, "stride")
	tok, err := mgr.Generate(uuid.New())
	require.NoError(t, err)

	parts := strings.Split(tok, ".")
	require.Len(t, parts, 3)
	tampered := parts[0] + "." + parts[1] + "." + parts[2] + "x"

	_, err = mgr.Parse(tampered)
	require.Error(t, err)
}

func TestParseRejectsDifferentSecret(t *testing.T) {
	signer := jwt.New("signer-secret", time.Hour, "stride")
	verifier := jwt.New("different-secret", time.Hour, "stride")

	tok, err := signer.Generate(uuid.New())
	require.NoError(t, err)

	_, err = verifier.Parse(tok)
	require.Error(t, err)
}

func TestParseRejectsExpiredToken(t *testing.T) {
	mgr := jwt.New("secret", time.Millisecond, "stride")
	tok, err := mgr.Generate(uuid.New())
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	_, err = mgr.Parse(tok)
	require.Error(t, err, "expired token should not parse")
}

func TestParseRejectsAlgNoneAttack(t *testing.T) {
	mgr := jwt.New("secret", time.Hour, "stride")

	uid := uuid.New()
	claims := jwt.Claims{
		UserID: uid,
		RegisteredClaims: jwtv5.RegisteredClaims{
			Subject:   uid.String(),
			ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	unsigned := jwtv5.NewWithClaims(jwtv5.SigningMethodNone, claims)
	tok, err := unsigned.SignedString(jwtv5.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	_, err = mgr.Parse(tok)
	require.Error(t, err, "alg=none tokens must be rejected")
}
