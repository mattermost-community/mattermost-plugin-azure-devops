package store

import (
	"fmt"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
	"github.com/pkg/errors"
)

func (s *Store) StoreOAuthState(mattermostUserID, state string) error {
	var ttlSeconds int64 = 60
	oAuthKey := fmt.Sprintf(constants.OAuthPrefix, mattermostUserID)
	return s.StoreTTL(oAuthKey, []byte(state), ttlSeconds)
}

func (s *Store) VerifyOAuthState(mattermostUserID, state string) error {
	oAuthKey := fmt.Sprintf(constants.OAuthPrefix, mattermostUserID)
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
