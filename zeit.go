package zeit

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

const timeLayout = "15:04:05"

var (
	ErrZeitLayout        = errors.New(`zeit error: time must be formatted as "15:04:05"`)
	ErrZeitType          = errors.New(`zeit error: type of scanned source cannot be asserted to []byte`)
	ErrZeitRangeLocation = errors.New(`zeit range error: location of times in ZeitRange must be equal`)
)

type Zeit struct {
	time.Time
}

// Now will return a new Zeit instant in UTC
func Now() Zeit {
	return Zeit{time.Now()}
}

// FromTime takes a time.Time instant and sources a Zeit instant from that
func FromTime(t time.Time) Zeit {
	return Zeit{t}
}

// Parse a string into a Zeit instant
//
// ex. Parse("23:10:05")
func Parse(src string) (Zeit, error) {
	var z Zeit
	t, err := time.Parse(timeLayout, src)

	if err == nil {
		z.Time = t
	}
	return z, err
}

// IsZero reports if the underlying time.Time instant is zero
func (z Zeit) IsZero() bool {
	return z.Time.IsZero()
}

func (z Zeit) String() string {
	return z.Format(timeLayout)
}

// MarshalJSON implements the encoding/json marshaller interface
func (z Zeit) MarshalJSON() ([]byte, error) {
	return []byte(`"` + z.Format(timeLayout) + `"`), nil
}

// UnmarshalJSON implements the encoding/json unmarshaller interface
func (z *Zeit) UnmarshalJSON(b []byte) error {
	s := string(b)

	if len(s) != 10 {
		return ErrZeitLayout
	}

	var ret time.Time
	var err error

	if ret, err = time.Parse(timeLayout, s[1:9]); err != nil {
		return err
	}

	z.Time = ret

	return nil
}

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

type ZeitRange [2]Zeit

// String implements the fmt.Stringer interface
func (zr ZeitRange) String() string {
	return fmt.Sprintf("%s, %s", zr[0].Format(timeLayout), zr[1].Format(timeLayout))
}

// Duration returns the duration between the two Zeit instants
func (zr ZeitRange) Duration() time.Duration {
	return zr[1].Sub(zr[0].Time)
}

// Split takes a Duration and splits a ZeitRange into
// as many ZeitRange instants as posible
func (zr ZeitRange) Split(d time.Duration) []*ZeitRange {
	return nil
}

// SplitExcept will split the ZeitRange into
// as many ZeitRanges as posible, without overlapping
// any of except []*ZeitRange
func (zr ZeitRange) SplitExcept(d time.Duration, except []*ZeitRange) []*ZeitRange {
	return nil
}

// RangeFromZeit creates a ZeitRange instant from two Zeit instants
func RangeFromZeit(from_z Zeit, to_z Zeit) (*ZeitRange, error) {
	if from_z.Location().String() != to_z.Location().String() {
		return nil, ErrZeitRangeLocation
	}

	return &ZeitRange{from_z, to_z}, nil
}

// RangeFromTime creates a ZeitRange instant from two time.Time instants
// useful if you need to use another Location
func RangeFromTime(from_t time.Time, to_t time.Time) (*ZeitRange, error) {
	from_z := FromTime(from_t)
	to_z := FromTime(to_t)

	return RangeFromZeit(from_z, to_z)
}

// RangeParse create a ZeitRange from two strings
//
// ex. RangeParse("10:15:00", "23:00:30")
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

// IsZero reports if either of the two underlying time.Time instants are zero
func (zr *ZeitRange) IsZero() bool {
	return zr[0].IsZero() || zr[1].IsZero()
}

// Within reports if a give time.Time instant is within
// the range of the two Zeit instants in the ZeitRange
func (zr *ZeitRange) Within(t time.Time) bool {
	tmp := time.Date(0, 0, 0, t.Hour(), t.Minute(), t.Second(), 0, t.Location())

	return tmp.After(zr[0].Time) && tmp.Before(zr[1].Time)
}
