package store

import (
	"github.com/pkg/errors"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
)

func (s *Store) StoreOAuthState(mattermostUserID, state string) error {
	oAuthKey := GetOAuthKey(mattermostUserID)
	return s.StoreTTL(oAuthKey, []byte(state), constants.TTLSecondsForOAuthState)
}

func (s *Store) VerifyOAuthState(mattermostUserID, state string) error {
	oAuthKey := GetOAuthKey(mattermostUserID)
	storedState, err := s.Load(oAuthKey)
	if err != nil {
		if err == ErrNotFound {
			return errors.New(constants.AuthAttemptExpired)
		}
		return err
	}

	if string(storedState) != state {
		return errors.New(constants.InvalidAuthState)
	}
	return nil
}
