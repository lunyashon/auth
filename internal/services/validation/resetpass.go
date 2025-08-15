package validate

import (
	"bytes"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errors = map[codes.Code]bool{
	codes.InvalidArgument: true,
}

func PassValidateWithNew(ctx context.Context, oldPassword, newPassword string) error {
	if err := checkPassword(oldPassword); err != nil {
		if _, ok := errors[status.Code(err)]; ok {
			var errMsg bytes.Buffer
			errMsg.WriteString(status.Convert(err).Message())
			errMsg.WriteString(": ")
			errMsg.WriteString("old password")
			return status.Errorf(status.Code(err), "%v", errMsg.String())
		}
		return status.Errorf(codes.Internal, "internal server error")
	}
	if err := passwordRegisterValidate(newPassword); err != nil {
		if _, ok := errors[status.Code(err)]; ok {
			var errMsg bytes.Buffer
			errMsg.WriteString(status.Convert(err).Message())
			errMsg.WriteString(": ")
			errMsg.WriteString("new password")
			return status.Errorf(status.Code(err), "%v", errMsg.String())
		}
		return status.Errorf(codes.Internal, "internal server error")
	}
	if oldPassword == newPassword {
		return status.Errorf(codes.InvalidArgument, "old and new passwords are the same")
	}
	return nil
}

func ValidateNewPassword(ctx context.Context, newPassword string) error {
	if err := passwordRegisterValidate(newPassword); err != nil {
		if _, ok := errors[status.Code(err)]; ok {
			var errMsg bytes.Buffer
			errMsg.WriteString(status.Convert(err).Message())
			errMsg.WriteString(": ")
			errMsg.WriteString("new password")
			return status.Errorf(status.Code(err), "%v", errMsg.String())
		}
		return status.Errorf(codes.Internal, "internal server error")
	}
	return nil
}

func PassValidateWithoutNew(ctx context.Context, oldPassword string) error {
	if err := passwordRegisterValidate(oldPassword); err != nil {
		if _, ok := errors[status.Code(err)]; ok {
			var errMsg bytes.Buffer
			errMsg.WriteString(status.Convert(err).Message())
			errMsg.WriteString(": ")
			errMsg.WriteString("old password")
			return status.Errorf(status.Code(err), "%v", errMsg.String())
		}
		return status.Errorf(codes.Internal, "internal server error")
	}
	return nil
}
