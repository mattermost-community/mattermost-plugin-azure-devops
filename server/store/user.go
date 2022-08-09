package store

type User struct {
	MattermostUserID string
	AccessToken      string
	RefreshToken     string
	ExpiresIn        string
}

func (s *Store) StoreUser(user *User) error {
	if err := s.StoreJSON(user.MattermostUserID, user); err != nil {
		return err
	}

	return nil
}

func (s *Store) LoadUser(mattermostUserID string) (*User, error) {
	user := User{}
	if err := s.LoadJSON(mattermostUserID, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) DeleteUser(mattermostUserID string) (bool, error) {
	if err := s.Delete(mattermostUserID); err != nil {
		return false, err
	}

	return true, nil
}
