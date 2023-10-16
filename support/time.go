package support

import (
	"errors"
	"fmt"
	"time"
)

type Time time.Time

func (d *Time) UnmarshalFormField(s string) error {
	t, err := time.Parse(time.RFC3339, s)

	if err != nil {
		return err
	}

	*d = Time(t)

	return nil
}

type Date time.Time

func (d *Date) UnmarshalJSON(data []byte) error {
	// remove double-quotes
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New(fmt.Sprintf("cannot parse \"%v\": missing double-quotes", data))
	}

	s := string(data[1 : len(data)-1])

	t, err := time.Parse(time.DateOnly, s)

	if err != nil {
		return err
	}

	*d = Date(t)

	return nil
}

func (d Date) String() string {
	return time.Time(d).String()
}
