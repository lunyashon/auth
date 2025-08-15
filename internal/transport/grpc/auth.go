package gRPCHudler

import (
	"context"
	"strings"

	"github.com/lunyashon/auth/internal/services/authgo"

	sso "github.com/lunyashon/protobuf/auth/gen/go/sso/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ServerAPI struct {
	sso.UnimplementedAuthServer
	auth *authgo.AuthData
}

func Register(gRPCServer *grpc.Server, auth *authgo.AuthData) {
	sso.RegisterAuthServer(gRPCServer, &ServerAPI{auth: auth})
}

// Realization gRPC method Login (get JWT token)
// Return gRPC login data (token and services) response or error
func (s *ServerAPI) Login(
	ctx context.Context,
	data *sso.LoginRequest,
) (*sso.LoginResponse, error) {

	device, ip, err := getDeviceAndIP(ctx)
	if err != nil {
		return nil, err
	}

	tokens, err := s.auth.LoginUser(ctx, data, device, ip)
	if err != nil {
		return nil, err
	}
	return &sso.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

// Realization gRPC method Register (register user)
// Return gRPC id new user or error
func (s *ServerAPI) Register(
	ctx context.Context,
	data *sso.RegisterRequest,
) (*sso.RegisterResponse, error) {

	id, err := s.auth.RegisterUser(ctx, data)
	if err != nil {
		return nil, err
	}
	return &sso.RegisterResponse{
		UserId: id,
	}, nil
}

// Realization gRPC method CreateToken (create token)
// Return gRPC token or error
func (s *ServerAPI) CreateToken(
	ctx context.Context,
	data *sso.TokenRequest,
) (*sso.TokenResponse, error) {

	token, err := s.auth.RegisterToken(ctx, data)
	return &sso.TokenResponse{
		Token: token,
	}, err
}

func (s *ServerAPI) GetServices(
	ctx context.Context,
	data *sso.ServicesRequest,
) (*sso.ServicesResponse, error) {
	services, err := s.auth.GetServices(ctx, data)
	if err != nil {
		return nil, err
	}
	return &sso.ServicesResponse{
		Services: services,
	}, nil
}

// Realization gRPC method logout user and off refresh and access token
// Return success (bool) or error
func (s *ServerAPI) Logout(
	ctx context.Context,
	data *sso.LogoutRequest,
) (*sso.LogoutResponse, error) {

	accessToken, err := checkAuth(ctx)
	if err != nil {
		return nil, err
	}

	err = s.auth.OnceLogoutUser(ctx, accessToken, data.RefreshToken)
	return &sso.LogoutResponse{
		Success: err == nil,
	}, err
}

func (s *ServerAPI) MassLogout(
	ctx context.Context,
	data *sso.MassLogoutRequest,
) (*sso.MassLogoutResponse, error) {

	accessToken, err := checkAuth(ctx)
	if err != nil {
		return nil, err
	}

	err = s.auth.MassLogoutUser(ctx, accessToken, data.RefreshToken)
	return &sso.MassLogoutResponse{
		Success: err == nil,
	}, err
}

// Realization gRPC method ValidateToken (validate token)
// Return gRPC user id or error
func (s *ServerAPI) ValidateToken(
	ctx context.Context,
	data *sso.ValidateRequest,
) (*sso.ValidateResponse, error) {

	id, err := s.auth.ValidateToken(ctx, data)
	return &sso.ValidateResponse{
		UserID: int64(id),
	}, err
}

// Realization gRPC method RefreshToken (refresh token)
// Return gRPC refresh token or error
func (s *ServerAPI) RefreshToken(
	ctx context.Context,
	data *sso.RefreshRequest,
) (*sso.RefreshResponse, error) {
	return &sso.RefreshResponse{}, nil
}

// Realization gRPC method UpdateAccessToken (update access token)
// Return gRPC access token or error
func (s *ServerAPI) UpdateAccessToken(
	ctx context.Context,
	data *sso.AccessTokenRequest,
) (*sso.AccessTokenResponse, error) {
	accessToken, err := s.auth.UpdateAccessToken(ctx, data)
	return &sso.AccessTokenResponse{
		AccessToken: accessToken,
	}, err
}

// Realization gRPC method ChangePassword (change password)
// Return gRPC success (bool) or error
func (s *ServerAPI) ChangePassword(
	ctx context.Context,
	data *sso.PasswordRequest,
) (*sso.PasswordResponse, error) {

	accessToken, err := checkAuth(ctx)
	if err != nil {
		return nil, err
	}

	err = s.auth.ChangePassword(
		ctx,
		accessToken,
		data.OldPassword,
		data.NewPassword,
	)

	return &sso.PasswordResponse{
		Success: err == nil,
	}, err
}

// Realization gRPC method ForgotPassword (forgot password)
// Return gRPC success (bool) or error
func (s *ServerAPI) ForgotPassword(
	ctx context.Context,
	data *sso.ForgotRequest,
) (*sso.ForgotResponse, error) {
	err := s.auth.ForgotPassword(ctx, data.Email)
	return &sso.ForgotResponse{
		Success: err == nil,
	}, err
}

// Realization gRPC method CheckForgotToken (check forgot token)
// Return gRPC success (bool) or error
func (s *ServerAPI) CheckForgotToken(
	ctx context.Context,
	data *sso.CheckForgotRequest,
) (*sso.CheckForgotResponse, error) {
	err := s.auth.CheckForgotToken(ctx, data.Token)
	return &sso.CheckForgotResponse{
		Success: err == nil,
	}, err
}

// Realization gRPC method ResetPassword (reset password)
// Return gRPC success (bool) or error
func (s *ServerAPI) ResetPassword(
	ctx context.Context,
	data *sso.ResetRequest,
) (*sso.ResetResponse, error) {
	err := s.auth.ResetPassword(ctx, data.Token, data.Password)
	return &sso.ResetResponse{
		Success: err == nil,
	}, err
}

// Realization gRPC method ConfirmEmail (confirm email)
// Return gRPC success (bool) or error
func (s *ServerAPI) ConfirmEmail(
	ctx context.Context,
	data *sso.EmailRequest,
) (*sso.EmailResponse, error) {
	accessToken, err := checkAuth(ctx)
	if err != nil {
		return nil, err
	}
	err = s.auth.ConfirmEmail(ctx, data.Email, accessToken)
	return &sso.EmailResponse{
		Success: err == nil,
	}, err
}

// Realization gRPC method CheckConfirmToken (check confirm token)
// Return gRPC success (bool) or error
func (s *ServerAPI) CheckConfirmToken(
	ctx context.Context,
	data *sso.CheckConfirmRequest,
) (*sso.CheckConfirmResponse, error) {
	accessToken, err := checkAuth(ctx)
	if err != nil {
		return nil, err
	}
	err = s.auth.CheckConfirmToken(ctx, data.Code, accessToken)
	return &sso.CheckConfirmResponse{
		Success: err == nil,
	}, err
}

// NO REALIZATION
// Realization gRPC method GetProfile (get profile)
// Return gRPC profile or error
func (s *ServerAPI) GetProfile(
	ctx context.Context,
	data *sso.ProfileRequest,
) (*sso.ProfileResponse, error) {
	accessToken, err := checkAuth(ctx)
	if err != nil {
		return nil, err
	}

	profile, err := s.auth.GetProfile(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *ServerAPI) GetMiniProfile(
	ctx context.Context,
	data *sso.ProfileRequest,
) (*sso.MiniProfileResponse, error) {
	accessToken, err := checkAuth(ctx)
	if err != nil {
		return nil, err
	}

	profile, err := s.auth.GetMiniProfile(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// NO USE
// Realization gRPC method get JWKS
// Return JWKS tokens or erro
func (s *ServerAPI) GetJWKS(
	ctx context.Context,
	data *sso.Empty,
) (*sso.JWKSResponse, error) {

	jwks := &sso.JWKSResponse{}
	jwks.Keys = append(
		jwks.Keys,
		s.auth.GetJWK(ctx),
	)
	return jwks, nil
}

func checkAuth(ctx context.Context) (string, error) {

	var accessToken string

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return accessToken, status.Errorf(codes.InvalidArgument, "missing metadata")
	}

	if len(md.Get("Authorization")) != 0 {
		token := md.Get("Authorization")[0]
		accessToken = strings.TrimPrefix(token, "Bearer ")
	} else {
		return accessToken, status.Errorf(codes.InvalidArgument, "the required Authorization parameter was not passed")
	}

	if accessToken == "" || accessToken == "undefined" {
		return accessToken, status.Errorf(codes.InvalidArgument, "the required Authorization parameter was not passed")
	}

	return accessToken, nil
}

func getDeviceAndIP(ctx context.Context) (string, string, error) {

	var (
		device, ip string
	)

	if md, ok := metadata.FromIncomingContext(ctx); ok {

		if len(md.Get("cf-connecting-ip")) > 0 {
			ip = md.Get("cf-connecting-ip")[0]
		}

		if len(md.Get("x-real-ip")) > 0 && ip == "" {
			ip = md.Get("x-real-ip")[0]
		}

		if len(md.Get("x-forwarded-for")) > 0 {
			ips := strings.Split(md.Get("x-forwarded-for")[0], ",")
			if len(ips) > 0 {
				ip = strings.TrimSpace(ips[0])
			}
		}

		if len(md.Get("x-client-ip")) > 0 && ip == "" {
			ip = md.Get("x-client-ip")[0]
		}

		if len(md.Get("grpcgateway-user-agent")) > 0 && device == "" {
			device = md.Get("grpcgateway-user-agent")[0]
		}
	}

	return device, ip, nil
}
