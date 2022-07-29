package store

import (
	"github.com/pkg/errors"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
)

func (s *Store) StoreOAuthState(mattermostUserID, state string) error {
	return s.StoreTTL(mattermostUserID, []byte(state), constants.TTLSecondsForOAuthState)
}

func (s *Store) VerifyOAuthState(mattermostUserID, state string) error {
	storedState, err := s.Load(mattermostUserID)
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
