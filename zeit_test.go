package zeit

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

const (
	t_bad_length = "083000"
	t_bad_char   = "14:30:AB"
	t_bad_value  = "14:30:69"

	t_from = "08:30:00"
	t_to   = "17:00:00"

	t2_from = "09:00:00"
	t2_to   = "11:30:00"

	t3_from = "14:00:00"
	t3_to   = "14:30:00"

	t4_from = "15:00:00"
	t4_to   = "16:00:00"
)

func TestParseBadZeit(t *testing.T) {
	var tests = []struct {
		val      string
		want_err error
	}{
		{t_from, nil},
		{t_to, nil},
		{t_bad_length, ErrZeitLayout},
	}

	for _, v := range tests {
		_, err := Parse(v.val)

		if err != v.want_err {
			t.Errorf(`expected error of '%s', got '%s'`, v.want_err, err)
		}
	}

	invalid_t := []string{t_bad_char, t_bad_value}

	for _, v := range invalid_t {
		if _, err := Parse(v); err == nil {
			t.Errorf("should fail on parsing %s", t_bad_char)
		}
	}
}

type ZeitTime struct {
	From Zeit `json:"from"`
	To   Zeit `json:"to"`
}

var zeit_list = []byte(`[
	{"from": "` + t_from + `", "to":"` + t_to + `"},
	{"from": "` + t2_from + `", "to":"` + t2_to + `"},
	{"from": "` + t3_from + `", "to":"` + t3_to + `"},
	{"from": "` + t4_from + `", "to":"` + t4_to + `"}
]`)

func TestZeitJson(t *testing.T) {
	zt := make([]ZeitTime, 0)

	if err := json.Unmarshal(zeit_list, &zt); err != nil {
		t.Error(err)
	}
}

func TestZeitRange(t *testing.T) {
	zr, err := ParseRange(t_from, t_to) // 08:30:00 - 17:00:00

	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		val   time.Duration
		wants int
	}{
		// {1 * time.Hour, 8},
		{2 * time.Hour, 4},
		// {4 * time.Hour, 2},
		// {5 * time.Hour, 1},
		// {6 * time.Hour, 1},
		// {7 * time.Hour, 1},
		// {510 * time.Minute, 1},
		// {511 * time.Minute, 0},
	}

	for _, v := range tests {
		zr_tmp := zr.Split(v.val)
		if len(zr_tmp) != v.wants {
			t.Errorf("expected %d ZeitRange instants, got %d", v.wants, len(zr_tmp))
			t.Error(zr_tmp)
		}
	}
	// 08:30:00
	// zr2, _ := ParseRange(t2_from, t2_to) // 09:00:00 - 11:30:00
	// zr3, _ := ParseRange(t3_from, t3_to) // 14:00:00 - 14:30:00
	// zr4, _ := ParseRange(t4_from, t4_to) // 15:00:00 - 16:00:00
	// 17:00:00
	// _exclude := []ZeitRange{zr2, zr3, zr4}
	exclude_empt := []ZeitRange{}

	arr := zr.Exclude(exclude_empt)

	offset := 1 * time.Hour
	for _, zr5 := range arr {
		for _, v := range tests {
			zr6 := zr5.SplitOffset(offset, v.val)
			for _, zr7 := range zr6 {
				fmt.Printf("%s - %s\n", zr7[0], zr7[1])
			}
		}
	}

	// for _, v := range tests {
	// 	zr_tmp := zr.Segment(excepts)
	// 	for _, zr5 := range zr_tmp {
	// 		fmt.Printf("%s - %s\n", zr5[0], zr5[1])
	// 	}
	// 	fmt.Println("---")
	// }
}
