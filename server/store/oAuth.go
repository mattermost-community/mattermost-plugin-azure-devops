package store

import (
	"github.com/pkg/errors"
)

func (s *Store) StoreOAuthState(mattermostUserID, state string) error {
	var ttlSeconds int64 = 60
	oAuthKey := GetOAuthKey(mattermostUserID)
	return s.StoreTTL(oAuthKey, []byte(state), ttlSeconds)
}

func (s *Store) VerifyOAuthState(mattermostUserID, state string) error {
	oAuthKey := GetOAuthKey(mattermostUserID)
	storedState, err := s.Load(oAuthKey)
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
