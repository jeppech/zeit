package zeit

import (
	"database/sql/driver"
	"time"
)

// Value implements the database/sql Valuer interface
func (z Zeit) Value() (driver.Value, error) {
	return driver.Value(z.Format(timeLayout)), nil
}

// Scan implements the database/sql Scanner interface
func (z *Zeit) Scan(v interface{}) error {
	var source string

	switch v := v.(type) {
	case string:
		source = v

	case []byte:
		source = string(v)

	default:
		return ErrZeitType
	}

	var ret time.Time
	var err error

	if ret, err = time.Parse(timeLayout, source); err != nil {
		return err
	}

	z.Time = ret

	return nil
}
