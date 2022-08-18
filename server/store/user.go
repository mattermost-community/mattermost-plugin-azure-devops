package store

<<<<<<< HEAD
import "github.com/Brightscout/mattermost-plugin-azure-devops/server/serializers"

type UserStore interface {
	LoadUser(mattermostUserID string) (*serializers.User, error)
	StoreUser(user *serializers.User) error
	DeleteUser(mattermostUserID string) (bool, error)
=======
type User struct {
	MattermostUserID string
	AccessToken      string
	RefreshToken     string
	ExpiresAt        int64
>>>>>>> [MI-1939]: Used epoch time
}

func (s *Store) StoreUser(user *serializers.User) error {
	if err := s.StoreJSON(GetOAuthKey(user.MattermostUserID), user); err != nil {
		return err
	}

	return nil
}

func (s *Store) LoadUser(mattermostUserID string) (*serializers.User, error) {
	user := serializers.User{}
	if err := s.LoadJSON(GetOAuthKey(mattermostUserID), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) DeleteUser(mattermostUserID string) (bool, error) {
	if err := s.Delete(GetOAuthKey(mattermostUserID)); err != nil {
		return false, err
	}

	return true, nil
}
