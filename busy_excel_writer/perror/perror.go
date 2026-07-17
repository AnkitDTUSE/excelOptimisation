package perror

import "fmt"

// MiscError formats an error with optional messages.
func MiscError(err error, msgs ...string) error {
	if err == nil {
		return nil
	}
	if len(msgs) > 0 {
		return fmt.Errorf("%s: %w", msgs[0], err)
	}
	return err
}
