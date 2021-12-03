package main

import (
	"github.com/goccy/go-json"
	"strconv"
	"time"
)

// SCATime is time.Time but it marshals and unmarshals to Unix Millis
// This will be moved to its own library soon
type SCATime struct {
	time.Time
}

func (st SCATime) MarshalJSON() ([]byte, error) {
	um := st.UnixMilli()
	str := strconv.FormatInt(um, 10)
	return []byte(str), nil
}

// UnmarshalJSON parses buf as an int64
// Postel's law: be conservative in what you send, be liberal in what you accept
func (st *SCATime) UnmarshalJSON(buf []byte) error {
	var tStr string

	err := json.Unmarshal(buf, &tStr)
	if err != nil {
		return err
	}

	t, err := strconv.ParseInt(tStr, 0, 64)
	if err != nil {
		return err
	}

	// This is bad, but I don't know of a better solution
	realT := time.UnixMilli(t)
	realTBuf, err := realT.GobEncode()
	if err != nil {
		return err
	}

	return st.GobDecode(realTBuf)
}