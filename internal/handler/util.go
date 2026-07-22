package handler

import (
	"github.com/ZhuJincheng-git/stride-backend/pkg/apperror"
)

func bindingError(err error) error {
	if err == nil {
		return nil
	}
	return apperror.New(apperror.CodeInvalidArgument, err.Error())
}