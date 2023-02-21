package store

import "github.com/mattermost/mattermost-plugin-azure-devops/server/serializers"

type UserStore interface {
	StoreAzureDevopsUserDetailsWithMattermostUserID(user *serializers.User) error
	LoadAzureDevopsUserIDFromMattermostUser(mattermostUserID string) (string, error)
	LoadAzureDevopsUserDetails(userID string) (*serializers.User, error)
	DeleteUser(mattermostUserID string) (bool, error)
}

func (s *Store) StoreAzureDevopsUserDetailsWithMattermostUserID(user *serializers.User) error {
	if err := s.StoreJSON(GetAzureDevopsUserKey(user.ID), user); err != nil {
		return err
	}

	if err := s.Store(GetOAuthKey(user.MattermostUserID), []byte(user.ID)); err != nil {
		return err
	}

	return nil
}

func (s *Store) LoadAzureDevopsUserIDFromMattermostUser(mattermostUserID string) (string, error) {
	azureDevopsUserID, err := s.Load(GetOAuthKey(mattermostUserID))
	if err != nil {
		return "", err
	}

	return string(azureDevopsUserID), nil
}

func (s *Store) LoadAzureDevopsUserDetails(userID string) (*serializers.User, error) {
	user := serializers.User{}
	if err := s.LoadJSON(GetAzureDevopsUserKey(userID), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) DeleteUser(mattermostUserID string) (bool, error) {
	azureDevopsUserID, err := s.LoadAzureDevopsUserIDFromMattermostUser(mattermostUserID)
	if err != nil {
		return false, err
	}

	if err := s.Delete(GetAzureDevopsUserKey(azureDevopsUserID)); err != nil {
		return false, err
	}

	if err := s.Delete(GetOAuthKey(mattermostUserID)); err != nil {
		return false, err
	}

	return true, nil
}
