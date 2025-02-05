//go:build !solution

package retryupdate

import (
	"errors"

	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

func UpdateValue(c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) error {
	var oldValue *string
	var oldVersion uuid.UUID
	var authErr *kvapi.AuthError
	var conflictErr *kvapi.ConflictError
	purchased := false
	for !purchased {
		getRes, err := c.Get(&kvapi.GetRequest{Key: key})
		if errors.Is(err, kvapi.ErrKeyNotFound) {
			purchased = true
			continue
		}
		if err == nil {
			purchased = true
			oldValue = &getRes.Value
			oldVersion = getRes.Version
			continue
		}
		if errors.As(err, &authErr) {
			return err
		}
	}
	updated, err := updateFn(oldValue)
	if err != nil {
		return err
	}
	newVer := uuid.Must(uuid.NewV4())
	isWritten := false
	for !isWritten {
		_, err := c.Set(&kvapi.SetRequest{Key: key, Value: updated, OldVersion: oldVersion, NewVersion: newVer})
		if errors.Is(err, kvapi.ErrKeyNotFound) {
			oldVersion = uuid.UUID{}
			updated, err = updateFn(nil)
			if err != nil {
				return err
			}
			continue
		}
		if err == nil || errors.As(err, &authErr) {
			return err
		}
		if errors.As(err, &conflictErr) {
			if conflictErr.ExpectedVersion == newVer {
				return nil
			}
			return UpdateValue(c, key, updateFn)
		}
	}
	return nil
}
