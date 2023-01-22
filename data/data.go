// Package data provides base types and structure for data interactions.
// This currently includes Task and TimeEntry data.
package data

import (
	"errors"
	"fmt"
	"github.com/timshannon/bolthold"
)

func translateNotFound(err error, id uint64) error {
	if errors.Is(err, bolthold.ErrNotFound) {
		return fmt.Errorf("%w: ID %d not found", ErrIDNotFound, id)
	}
	return err
}
