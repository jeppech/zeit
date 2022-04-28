package zeit

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

const timeLayout = "15:04:05"

var (
	ErrZeitLayout        = errors.New(`zeit error: time must be formatted as "` + timeLayout + `"`)
	ErrZeitType          = errors.New(`zeit error: type of scanned source cannot be asserted to []byte`)
	ErrZeitRangeLocation = errors.New(`zeit range error: location of times in ZeitRange must be equal`)
	ErrZeitRangeDuration = errors.New(`zeit range error: duration is larger than 24 hours`)
)

type Zeit struct {
	time.Time
}

// Now will return a new Zeit instant in UTC
func Now() Zeit {
	return Zeit{time.Now().UTC()}
}

// NowInLoc will return a new Zeit instant in the given time.Location
func NowInLoc(loc *time.Location) Zeit {
	return Zeit{time.Now().In(loc)}
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

// ParseInloc parses a string in a given time.Location
func ParseInLoc(src string, loc *time.Location) (Zeit, error) {
	var z Zeit
	t, err := time.Parse(timeLayout, src)

	if err == nil {
		z.Time = t.In(loc)
	}
	return z, err
}

func (z Zeit) Add(d time.Duration) Zeit {
	z.Time = z.Time.Add(d)
	return z
}

// Equal asserts if the stringified version are equal
func (z Zeit) Equal(t Zeit) bool {
	return z.String() == t.String()
}

func (z Zeit) Before(t Zeit) bool {
	return z.Time.Before(t.Time)
}

func (z Zeit) After(t Zeit) bool {
	return z.Time.After(t.Time)
}

// IsZero reports if the underlying time.Time instant is zero
func (z Zeit) IsZero() bool {
	return z.Time.IsZero()
}

// String implements the fmt.Stringer interface
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

// RangeFromZeit creates a ZeitRange instant from two Zeit instants
func RangeFromZeit(from_z Zeit, to_z Zeit) (ZeitRange, error) {
	if from_z.Location().String() != to_z.Location().String() {
		return ZeitRange{}, ErrZeitRangeLocation
	}

	if to_z.Sub(from_z.Time) >= 24*time.Hour {
		return ZeitRange{}, ErrZeitRangeDuration
	}

	return ZeitRange{from_z, to_z}, nil
}

// RangeFromTime creates a ZeitRange instant from two time.Time instants
func RangeFromTime(from_t time.Time, to_t time.Time) (ZeitRange, error) {
	from_z := FromTime(from_t)
	to_z := FromTime(to_t)

	return RangeFromZeit(from_z, to_z)
}

// RangeParseInLoc create a ZeitRange from two strings in the given time.Location
func RangeParseInLoc(from string, to string, loc *time.Location) (ZeitRange, error) {
	var z ZeitRange
	z_from, err := ParseInLoc(from, loc)

	if err != nil {
		return z, err
	}

	z_to, err := ParseInLoc(from, loc)
	if err != nil {
		return z, err
	}

	return RangeFromZeit(z_from, z_to)
}

// RangeParse create a ZeitRange from two strings in UTC
//
// ex. RangeParse("10:15:00", "23:00:30")
func RangeParse(from string, to string) (ZeitRange, error) {
	var z ZeitRange
	z_from, err := Parse(from)

	if err != nil {
		return z, err
	}

	z_to, err := Parse(to)
	if err != nil {
		return z, err
	}

	return RangeFromZeit(z_from, z_to)
}

// ToZeit returns the two individual Zeit instants
func (zr ZeitRange) ToZeit() (Zeit, Zeit) {
	return zr[0], zr[1]
}

// Split takes a Duration and splits a ZeitRange into
// as many ZeitRange instants as posible
func (zr ZeitRange) Split(d time.Duration) []ZeitRange {
	z1, z2 := zr.ToZeit()

	diff := z2.Sub(z1.Time)
	slots := int(diff / d)

	z_range := make([]ZeitRange, slots)
	start := z1

	for i := 0; i < slots; i++ {
		next := start.Add(d)
		z_range[i] = ZeitRange{start, next}
		start = next
	}

	return z_range
}

// SplitExcept will split the ZeitRange into
// as many ZeitRanges as posible, without overlapping
// any of except []*ZeitRange
func (zr ZeitRange) SplitExcept(d time.Duration, except []ZeitRange) []ZeitRange {
	z_range := zr.Split(d)
	z_range2 := make([]ZeitRange, 0, len(z_range))

	fmt.Println(len(z_range))
	var is_excp bool

	for _, zr2 := range z_range {
		is_excp = false
		for _, excp := range except {
			z1, z2 := zr2.ToZeit()
			if !excp.Within(z1) && !excp.Within(z2) {
				continue
			}
			is_excp = true
			break
		}

		if !is_excp {
			z_range2 = append(z_range2, zr2)
		}
	}
	return z_range2
}

// Add given duration to both Zeit instants
func (zr ZeitRange) Add(d time.Duration) (ZeitRange, error) {
	z1, z2 := zr.ToZeit()

	return RangeFromZeit(z1.Add(d), z2.Add(d))
}

// Within reports if a give time.Time instant is within
// the range of the two Zeit instants in the ZeitRange
//
// note: the locat
func (zr ZeitRange) Within(z Zeit) bool {
	return z.After(zr[0]) && z.Before(zr[1])
}

// String implements the fmt.Stringer interface
func (zr ZeitRange) String() string {
	return fmt.Sprintf("%s, %s", zr[0].Format(timeLayout), zr[1].Format(timeLayout))
}

// Duration returns the duration between the two Zeit instants
func (zr ZeitRange) Duration() time.Duration {
	return zr[1].Sub(zr[0].Time)
}

// IsZero reports if either of the two underlying time.Time instants are zero
func (zr ZeitRange) IsZero() bool {
	return zr[0].IsZero() || zr[1].IsZero()
}
