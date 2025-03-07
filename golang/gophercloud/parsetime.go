package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gophercloud/gophercloud"
)

type Heartbeat struct {
	HeartbeatTimestamp time.Time `json:"-"`
}

type JSONRFC3339ZNoTNoZOVN time.Time

func (jt *JSONRFC3339ZNoTNoZOVN) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}

	t, err := time.Parse(gophercloud.RFC3339ZNoTNoZ, s)
	if err != nil {
		// if fail try to parse with different format
		t, err = time.Parse("2006-01-02 15:04:05.999999-07:00", s)
		if err != nil {
			return err
		}
	}
	*jt = JSONRFC3339ZNoTNoZOVN(t)
	return nil
}

// UnmarshalJSON helps to convert the timestamps into the time.Time type.
func (r *Heartbeat) UnmarshalJSON(b []byte) error {
	type tmp Heartbeat
	var s struct {
		tmp
		HeartbeatTimestamp JSONRFC3339ZNoTNoZOVN `json:"heartbeat_timestamp"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Heartbeat(s.tmp)

	r.HeartbeatTimestamp = time.Time(s.HeartbeatTimestamp)

	return nil
}

func main() {
	// Test input with different timestamp formats
	json1 := `{"heartbeat_timestamp":"2023-02-23 14:12:07.471000+00:00"}`
	json2 := `{"heartbeat_timestamp":"2023-02-23 14:12:07.471"}`

	var heartbeat1, heartbeat2 Heartbeat
	if err := json.Unmarshal([]byte(json1), &heartbeat1); err != nil {
		fmt.Println("Error unmarshalling json1:", err)
	} else {
		fmt.Println("Heartbeat1 Timestamp:", heartbeat1.HeartbeatTimestamp)
	}

	if err := json.Unmarshal([]byte(json2), &heartbeat2); err != nil {
		fmt.Println("Error unmarshalling json2:", err)
	} else {
		fmt.Println("Heartbeat2 Timestamp:", heartbeat2.HeartbeatTimestamp)
	}
}
