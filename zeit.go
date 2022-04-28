package zeit

import (
	"database/sql/driver"
	"errors"
	"time"
)

const timeLayout = "15:04:05"

var ZeitLayoutError = errors.New(`zeit error: time must be formatted as "15:04:05"`)
var ZeitTypeError = errors.New(`zeit error: type of scanned source cannot be asserted to []byte`)
var ZeitRangeLocationError = errors.New(`zeit range error: location of times in ZeitRange must be equal`)

type Zeit struct {
	time.Time
}

func Now() Zeit {
	return Zeit{time.Now()}
}

func FromTime(t time.Time) Zeit {
	return Zeit{t}
}

func Parse(src string) (Zeit, error) {
	var z Zeit
	t, err := time.Parse(timeLayout, src)

	if err == nil {
		z.Time = t
	}
	return z, err
}

func (z Zeit) IsZero() bool {
	return z.Time.IsZero()
}

func (z Zeit) MarshalJSON() ([]byte, error) {
	return []byte(`"` + z.Format(timeLayout) + `"`), nil
}

func (z *Zeit) UnmarshalJSON(b []byte) error {
	s := string(b)

	if len(s) != 10 {
		return ZeitLayoutError
	}

	var ret time.Time
	var err error

	if ret, err = time.Parse(timeLayout, s[1:9]); err != nil {
		return err
	}

	z.Time = ret

	return nil
}

func (z Zeit) Value() (driver.Value, error) {
	return driver.Value(z.Format(timeLayout)), nil
}

func (z *Zeit) Scan(v interface{}) error {
	var source string

	switch v := v.(type) {
	case string:
		source = v

	case []byte:
		source = string(v)

	default:
		return ZeitTypeError
	}

	var ret time.Time
	var err error

	if ret, err = time.Parse(timeLayout, source); err != nil {
		return err
	}

	z.Time = ret

	return nil
}

type ZeitRange [2]Zeit

func RangeParse(from string, to string) (*ZeitRange, error) {
	z_from, err := Parse(from)

	if err != nil {
		return nil, err
	}

	z_to, err := Parse(from)
	if err != nil {
		return nil, err
	}

	return RangeFromZeit(z_from, z_to)
}

func RangeFromZeit(from_z Zeit, to_z Zeit) (*ZeitRange, error) {
	if from_z.Location().String() != to_z.Location().String() {
		return nil, ZeitRangeLocationError
	}

	return &ZeitRange{from_z, to_z}, nil
}

func RangeFromTime(from_t time.Time, to_t time.Time) (*ZeitRange, error) {
	from_z := FromTime(from_t)
	to_z := FromTime(to_t)

	return RangeFromZeit(from_z, to_z)
}

func (zr ZeitRange) IsZero() bool {
	return zr[0].IsZero() || zr[1].IsZero()
}

func (zr ZeitRange) Within(t time.Time) bool {
	tmp := time.Date(0, 0, 0, t.Hour(), t.Minute(), t.Second(), 0, t.Location())

	return tmp.After(zr[0].Time) && tmp.Before(zr[1].Time)
}
