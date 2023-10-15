package support

import (
	"time"
)

type Date time.Time

func (d *Date) UnmarshalFormField(s string) error {
	t, err := time.Parse(time.RFC3339, s)

	if err != nil {
		return err
	}

	*d = Date(t)

	return nil
}
