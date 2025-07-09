package authgo

import (
	"bytes"
	"context"
	"fmt"

	jwtsso "github.com/lunyashon/auth/internal/lib/jwt"
	"github.com/lunyashon/auth/internal/lib/passauth"
	validate "github.com/lunyashon/auth/internal/services/validation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errors = map[codes.Code]struct {
	message string
}{
	codes.InvalidArgument:  {},
	codes.PermissionDenied: {},
	codes.Internal:         {},
	codes.NotFound:         {},
	codes.Unauthenticated:  {},
}

func (a *AuthData) ChangePassword(
	ctx context.Context,
	accessToken,
	oldPassword,
	newPassword string,
) error {

	userId, err := jwtsso.ValidateAccessToken(
		accessToken,
		"sso.auth.change.password",
		a.Yaml.NameSSOService,
		a.KeysStore.PublicKey,
	)
	if err != nil {
		if _, ok := errors[status.Code(err)]; ok {
			return err
		}
		a.Log.ErrorContext(
			ctx,
			"ERROR change password",
			"message", err.Error(),
		)
		return status.Errorf(codes.Internal, "internal server error")
	}

	param, err := a.DB.User.GetParamByUserId(ctx, userId)
	if err != nil {
		return err
	}

	rb := passauth.ExecAuthService(&passauth.RealBcrypt{})

	fmt.Println(param)

	if err := rb.VerifyPassword(oldPassword, param.Password); err != nil {
		if _, ok := errors[status.Code(err)]; ok {
			var errMsg bytes.Buffer
			errMsg.WriteString(status.Convert(err).Message())
			errMsg.WriteString(": ")
			errMsg.WriteString("old password")
			return status.Errorf(status.Code(err), "%v", errMsg.String())
		}
		a.Log.ErrorContext(
			ctx,
			"ERROR change password",
			"accessToken", accessToken,
			"message", err.Error(),
		)
		return status.Errorf(codes.Internal, "internal server error")
	}

	if err := validate.PassValidateWithNew(ctx, oldPassword, newPassword); err != nil {
		return err
	}

	return a.password(ctx, newPassword, userId, rb)
}

func (a *AuthData) ResetPassword(
	ctx context.Context,
	token string,
	password string,
) error {
	rb := passauth.ExecAuthService(&passauth.RealBcrypt{})

	userId, err := a.DB.Token.CheckForgotToken(ctx, token)
	if err != nil {
		return err
	}

	param, err := a.DB.User.GetParamByUserId(ctx, userId)
	if err != nil {
		return err
	}

	if err := rb.VerifyPassword(password, param.Password); err != nil {
		if _, ok := errors[status.Code(err)]; !ok {
			a.Log.ErrorContext(
				ctx,
				"ERROR reset password",
				"message", err.Error(),
			)
			return status.Errorf(codes.Internal, "internal server error")
		}
	} else {
		return status.Errorf(codes.InvalidArgument, "old and new passwords are the same")
	}

	if err := validate.PassValidateWithNew(ctx, password, password); err != nil {
		return err
	}

	if err := a.password(ctx, password, userId, rb); err != nil {
		return err
	}

	if err := a.DB.Token.DeleteForgotToken(ctx, token); err != nil {
		return err
	}

	return nil
}

func (a *AuthData) password(
	ctx context.Context,
	password string,
	userId int,
	rb *passauth.PassAuthService,
) error {

	hashPass, err := rb.GeneratePassword([]byte(password))
	if err != nil {
		if _, ok := errors[status.Code(err)]; ok {
			return err
		}
		a.Log.ErrorContext(
			ctx,
			"ERROR change password",
			"message", err.Error(),
		)
		return status.Errorf(codes.Internal, "internal server error")
	}

	if err := a.DB.User.ChangePassword(ctx, userId, string(hashPass)); err != nil {
		return err
	}

	return nil
}
