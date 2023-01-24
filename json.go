package zeit

import "time"

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
