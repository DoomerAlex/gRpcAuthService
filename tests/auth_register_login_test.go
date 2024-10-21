package tests

import (
	"gRpcAuthService/tests/suite"
	ssov1 "github.com/DoomerAlex/gRpcAuthProtos/gen/go/gRpcAuthService"
	"github.com/brianvoe/gofakeit"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	emptyAppID     = 0
	appId          = 1
	appSecret      = "mySecret"
	passDefaultLen = 10
)

func TestRegisterLogin_login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)
	// Генерит рандомный email
	email := gofakeit.Email()
	pass := randomFakePassword()

	respReq, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{Email: email, Password: pass})

	require.NoError(t, err)
	assert.NotEmpty(t, respReq.GetUserId())

	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{Email: email, Password: pass, AppId: appId})

	require.NoError(t, err)

	loginTime := time.Now()

	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, respReq.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appId, int(claims["app_id"].(float64)))

	const deltaSeconds = 1
	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func TestRegisterLogin_DuplicatedRegistration(t *testing.T) {
	ctx, st := suite.New(t)
	// Генерит рандомный email
	email := gofakeit.Email()
	pass := randomFakePassword()

	respReq, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{Email: email, Password: pass})

	require.NoError(t, err)
	assert.NotEmpty(t, respReq.GetUserId())

	resp2Req, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{Email: email, Password: pass})

	require.Error(t, err)
	assert.Empty(t, resp2Req.GetUserId())
	assert.ErrorContains(t, err, "user already exists")
}

func TestRegister_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "Register with Empty Password",
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "password is required",
		},
		{
			name:        "Register with Empty Email",
			email:       "",
			password:    randomFakePassword(),
			expectedErr: "email is required",
		},
		{
			name:        "Register with Empty all parameters",
			email:       "",
			password:    "",
			expectedErr: "email is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    tt.email,
				Password: tt.password,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestLogin_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		appID       int32
		expectedErr string
	}{
		{
			name:        "Login with Empty Password",
			email:       gofakeit.Email(),
			password:    "",
			appID:       appId,
			expectedErr: "password is required",
		},
		{
			name:        "Login with Empty email",
			email:       "",
			password:    randomFakePassword(),
			appID:       appId,
			expectedErr: "email is required",
		},
		{
			name:        "Login with Empty app_id",
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			appID:       emptyAppID,
			expectedErr: "app_id is required",
		},
		{
			name:        "Login with Empty all parameters",
			email:       "",
			password:    "",
			appID:       emptyAppID,
			expectedErr: "email is required",
		},
		{
			name:        "Login with Non-Matching Password",
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			appID:       appId,
			expectedErr: "invalid email or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
				AppId:    tt.appID,
			})

			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
