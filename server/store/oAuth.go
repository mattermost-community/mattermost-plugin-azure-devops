package store

import (
	"github.com/pkg/errors"
)

func (s *Store) StoreOAuthState(mattermostUserID, state string) error {
	var ttlSeconds int64 = 60
	return s.StoreTTL(mattermostUserID, []byte(state), ttlSeconds)
}

func (s *Store) VerifyOAuthState(mattermostUserID, state string) error {
	storedState, err := s.Load(mattermostUserID)
	if err != nil {
		if err == ErrNotFound {
			return errors.New("authentication attempt expired, please try again")
		}
		return err
	}

	if string(storedState) != state {
		return errors.New("invalid oauth state, please try again")
	}
	return nil
}
