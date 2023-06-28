package internal_error

import "errors"

var (
	ErrInternalError     = errors.New("internal error")
	ErrGmailNotFound     = errors.New("gmail not found")
	ErrUsernameNotFound  = errors.New("username not found")
	ErrUserIDNotFound    = errors.New("user ID not found")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrIncompleteProfile = errors.New("incomplete profile")
	ErrUnverified        = errors.New("unverified account")

	ErrRecommendationLimitReached = errors.New("recommendation limit reached")
	ErrRecommendationNotReady     = errors.New("recommendation not ready")
	ErrNoMatchingProfileFound     = errors.New("no matching profile found")
	ErrStillInDiscussion          = errors.New("cannot make any action while in discussion state")
	ErrNoDiscussionInProgress     = errors.New("user has no discussion in progress")
)
