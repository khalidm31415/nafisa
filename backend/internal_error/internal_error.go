package internal_error

import "errors"

var (
	ErrInternalError    = errors.New("internal error")
	ErrUsernameNotFound = errors.New("username not found")
	ErrUserIDNotFound   = errors.New("user ID not found")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrUnverified       = errors.New("unverified account")

	ErrRecommendationLimitReached = errors.New("recommendation limit reached")
	ErrRecommendationNotReady     = errors.New("recommendation not ready")
	ErrNoMatchingProfileFound     = errors.New("no matching profile found")
	ErrStillInDiscussion          = errors.New("cannot make any action while in discussion state")
)
