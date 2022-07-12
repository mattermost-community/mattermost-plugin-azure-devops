package store

import "github.com/pkg/errors"

type User struct {
	MattermostUserID string
	AccessToken      string
	RefreshToken     string
	ExpiresIn        string
}

func (s *Store) StoreUser(user *User) error {
	err := s.StoreJSON(user.MattermostUserID, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) LoadUser(mattermostUserID string) (*User, error) {
	user := User{}
	err := s.LoadJSON(mattermostUserID, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) DeleteUser(mattermostUserID string) bool {
	err := s.Delete(mattermostUserID)
	if err != nil {
		errors.Wrap(err, err.Error())
		return false
	}

	return true
}
