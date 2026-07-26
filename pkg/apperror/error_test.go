package apperror_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/ZhuJincheng-git/stride-backend/pkg/apperror"
	"github.com/stretchr/testify/require"
)

func TestNewAndError(t *testing.T) {
	err := apperror.New(apperror.CodeNotFound, "user not found")
	require.Equal(t, apperror.CodeNotFound, err.Code)
	require.Equal(t, "user not found", err.Message)
	require.Equal(t, "not_found: user not found", err.Error())
}

func TestNewfFormats(t *testing.T) {
	err := apperror.Newf(apperror.CodeInvalidArgument, "field %s is required", "title")
	require.Equal(t, "field title is required", err.Message)
}

func TestWrapKeepsCauseAccessibleViaErrorsIs(t *testing.T) {
	cause := fmt.Errorf("connect timeout")
	err := apperror.Wrap(apperror.CodeInternal, "talking to db", cause)

	require.True(t, errors.Is(err, cause), "wrapped error should chain via Unwrap")
	require.Contains(t, err.Error(), "connect timeout", "string form should include cause")
}

func TestStatusCode(t *testing.T) {
	cases := []struct {
		code apperror.Code
		want int
	}{
		{apperror.CodeInvalidArgument, http.StatusBadRequest},
		{apperror.CodeUnauthenticated, http.StatusUnauthorized},
		{apperror.CodePermissionDenied, http.StatusForbidden},
		{apperror.CodeNotFound, http.StatusNotFound},
		{apperror.CodeInternal, http.StatusInternalServerError},
		{apperror.Code("totally_make_up"), http.StatusInternalServerError},
	}
	for _, tc := range cases {
		t.Run(string(tc.code), func(t *testing.T) {
			require.Equal(t, tc.want, apperror.New(tc.code, "x").StatusCode())
		})
	}
}

func TestAsAppErrorExtractsFromWrappedChain(t *testing.T) {
	inner := apperror.New(apperror.CodePermissionDenied, "nope")
	wrapped := fmt.Errorf("call to service failed: %w", inner)

	got, ok := apperror.AsAppError(wrapped)
	require.True(t, ok)
	require.Same(t, inner, got)
}

func TestAsAppErrorReturnsFalseForNonAppError(t *testing.T) {
	got, ok := apperror.AsAppError(fmt.Errorf("random"))
	require.False(t, ok)
	require.Nil(t, got)

	got, ok = apperror.AsAppError(nil)
	require.False(t, ok)
	require.Nil(t, got)
}
