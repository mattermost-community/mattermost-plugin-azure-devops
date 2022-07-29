package store

import (
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/pkg/errors"
)

type Store struct {
	api plugin.API
}

func NewStore(api plugin.API) *Store {
	return &Store{
		api,
	}
}

func (s *Store) Load(key string) ([]byte, error) {
	data, appErr := s.api.KVGet(key)
	if appErr != nil {
		return nil, errors.WithMessage(appErr, "failed plugin KVGet")
	}
	if data == nil {
		return nil, nil
	}
	return data, nil
}

func (s *Store) Store(key string, data []byte) error {
	if appErr := s.api.KVSet(key, data); appErr != nil {
		return errors.WithMessagef(appErr, "failed plugin KVSet %q", key)
	}
	return nil
}

func (s *Store) StoreTTL(key string, data []byte, ttlSeconds int64) error {
	appErr := s.api.KVSetWithExpiry(key, data, ttlSeconds)
	if appErr != nil {
		return errors.WithMessagef(appErr, "failed plugin KVSet (ttl: %vs) %q", ttlSeconds, key)
	}
	return nil
}

func (s *Store) StoreWithOptions(key string, value []byte, opts model.PluginKVSetOptions) (bool, error) {
	success, appErr := s.api.KVSetWithOptions(key, value, opts)
	if appErr != nil {
		return false, errors.WithMessagef(appErr, "failed plugin KVSet (ttl: %vs) %q", opts.ExpireInSeconds, key)
	}
	return success, nil
}

func (s *Store) Delete(key string) error {
	appErr := s.api.KVDelete(key)
	if appErr != nil {
		return errors.WithMessagef(appErr, "failed plugin KVDelete %q", key)
	}
	return nil
}
