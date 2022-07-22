package store

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
)

var ErrNotFound = errors.New("not found")

// Ensure makes sure the initial value for a key is set to the value provided, if it does not already exists
// Returns the value set for the key in kv-store and error
func (s *Store) Ensure(key string, newValue []byte) ([]byte, error) {
	value, err := s.Load(key)
	switch err {
	case nil:
		return value, nil
	case ErrNotFound:
		break
	default:
		return nil, err
	}

	if err = s.Store(key, newValue); err != nil {
		return nil, err
	}

	// Load again in case we lost the race to another server
	value, err = s.Load(key)
	if err != nil {
		return newValue, nil
	}
	return value, nil
}

// LoadJSON loads a json value stored in the KVStore using StoreJSON
// unmarshalling it to an interface using json.Unmarshal
func (s *Store) LoadJSON(key string, v interface{}) (returnErr error) {
	data, err := s.Load(key)
	if err != nil {
		return err
	}

	if data == nil {
		return nil
	}

	return json.Unmarshal(data, v)
}

// StoreJSON stores a json value from an interface to KVStore after marshaling it using json.Marshal
func (s *Store) StoreJSON(key string, v interface{}) (returnErr error) {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return s.Store(key, data)
}

// AtomicModifyWithOptions modifies the value for a key in KVStore, only if the initial value is not changed while attempting to modify it.
// To avoid race conditions, we retry the modification multiple times after waiting for a fixed wait interval.
// input: kv store key and a modify function which takes initial value and returns final value with PluginKVSetOptions and error.
// returns: nil for a successful update and error if the update is unsuccessful or the retry limit reached.
func (s *Store) AtomicModifyWithOptions(key string, modify func(initialValue []byte) ([]byte, *model.PluginKVSetOptions, error)) error {
	currentAttempt := 0
	for {
		initialBytes, appErr := s.Load(key)
		if appErr != nil && appErr != ErrNotFound {
			return errors.Wrap(appErr, "unable to read initial value")
		}

		newValue, opts, err := modify(initialBytes)
		if err != nil {
			return errors.Wrap(err, "modification error")
		}

		// No modifications have been done. No reason to hit the plugin API.
		if bytes.Equal(initialBytes, newValue) {
			return nil
		}

		if opts == nil {
			opts = &model.PluginKVSetOptions{}
		}
		opts.Atomic = true
		opts.OldValue = initialBytes
		success, setError := s.StoreWithOptions(key, newValue, *opts)
		if setError != nil {
			return errors.Wrap(setError, "problem writing value")
		}
		if success {
			return nil
		}

		currentAttempt++
		if currentAttempt >= constants.AtomicRetryLimit {
			return errors.New("reached write attempt limit")
		}

		time.Sleep(constants.AtomicRetryWait)
	}
}

// AtomicModify calls AtomicModifyWithOptions with nil PluginKVSetOptions
// to atomically modify a value in KVStore and set it for an indefinite time
// See AtomicModifyWithOptions for more info
func (s *Store) AtomicModify(key string, modify func(initialValue []byte) ([]byte, error)) error {
	return s.AtomicModifyWithOptions(key, func(initialValue []byte) ([]byte, *model.PluginKVSetOptions, error) {
		b, err := modify(initialValue)
		return b, nil, err
	})
}
