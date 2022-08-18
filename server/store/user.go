package store

type User struct {
	MattermostUserID string
	AccessToken      string
	RefreshToken     string
	ExpiresAt        int64
}

func (s *Store) StoreUser(user *User) error {
	if err := s.StoreJSON(GetOAuthKey(user.MattermostUserID), user); err != nil {
		return err
	}

	return nil
}

func (s *Store) LoadUser(mattermostUserID string) (*User, error) {
	user := User{}
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
