package badger3

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/colbee1/jobq"
)

const (
	keyFieldSeparator = "\x1f"

	prefixKeyTopicQueue   = "tq:"
	prefixKeyTopicInfo    = "ti:"
	prefixKeyDelayedQueue = "dq:"
)

//-----------------------------------------------------------------------
//	Topic Info
//-----------------------------------------------------------------------

type keyTopicInfo struct {
	Topic jobq.Topic
}

func newKeyTopicInfo(topic jobq.Topic) keyTopicInfo {
	return keyTopicInfo{Topic: topic}
}

func (k keyTopicInfo) Info() []byte {
	key := fmt.Sprintf("%st=%s", prefixKeyTopicInfo, k.Topic)

	return []byte(key)
}

type modelTopicInfo struct {
	DateCreated time.Time `json:"dateCreated"`
}

func (mti *modelTopicInfo) Decode(data []byte) error {
	n := modelTopicInfo{}
	if err := json.Unmarshal(data, &n); err != nil {
		return err
	}
	*mti = n

	return nil
}

func (mti *modelTopicInfo) Encode() ([]byte, error) {
	return json.Marshal(mti)
}

//-----------------------------------------------------------------------
//	Topic Item
//-----------------------------------------------------------------------

type keyTopicItem struct {
	Topic  jobq.Topic
	Weight jobq.Weight
	JobID  jobq.ID
}

func newKeyTopicItem(topic jobq.Topic, w jobq.Weight, jid jobq.ID) keyTopicItem {
	return keyTopicItem{Topic: topic, Weight: w, JobID: jid}
}

func (k keyTopicItem) Item() []byte {
	key := fmt.Sprintf("%st=%s\x1fw=%020d\x1fid=%020d",
		prefixKeyTopicQueue,
		k.Topic,
		int(k.Weight)+jobq.WeightMax,
		k.JobID,
	)

	return []byte(key)
}

func (k keyTopicItem) Prefix() []byte {
	key := fmt.Sprintf("%st=%s", prefixKeyTopicQueue, k.Topic)

	return []byte(key)
}

func (k *keyTopicItem) Decode(data []byte) error {
	s := string(data)
	parts := strings.Split(s, keyFieldSeparator)
	if len(parts) != 3 {
		return fmt.Errorf("expected 3 fields in topic item key, got: %d", len(parts))
	}

	topic, weight, jid := parts[0], parts[1], parts[2]

	prefix := prefixKeyTopicQueue + "t="
	if !strings.HasPrefix(topic, prefix) {
		return fmt.Errorf("wrong prefix for topic in topic item key")
	}
	topic = strings.TrimPrefix(topic, prefix)

	prefix = "w="
	if !strings.HasPrefix(weight, prefix) {
		return fmt.Errorf("wrong prefix for weight in topic item key")
	}
	weight = strings.TrimPrefix(weight, prefix)

	prefix = "id="
	if !strings.HasPrefix(jid, prefix) {
		return fmt.Errorf("wrong prefix for id in topic item key")
	}
	jid = strings.TrimPrefix(jid, prefix)

	k.Topic = jobq.NewTopicFromString(topic)
	if v, err := jobq.NewJobIDFromString(jid); err != nil {
		return err
	} else {
		k.JobID = v
	}
	if v, err := jobq.NewWeightFromString(weight); err != nil {
		return err
	} else {
		k.Weight = v - jobq.WeightMax
	}

	return nil
}

//-----------------------------------------------------------------------
//	Delayed Item
//-----------------------------------------------------------------------

type keyDelayedItem struct {
	At    time.Time
	JobID jobq.ID
}

func newKeyDelayedItem(delayedAt time.Time, jid jobq.ID) keyDelayedItem {
	return keyDelayedItem{At: delayedAt, JobID: jid}
}

func (k keyDelayedItem) Item() []byte {
	key := fmt.Sprintf("%sts=%020d\x1fid=%020d", prefixKeyDelayedQueue, k.At.Unix(), k.JobID)

	return []byte(key)
}

func (k keyDelayedItem) Prefix() []byte {
	key := fmt.Sprintf("%sts=", prefixKeyDelayedQueue)

	return []byte(key)
}

func (k *keyDelayedItem) Decode(data []byte) error {
	s := string(data)
	parts := strings.Split(s, keyFieldSeparator)
	if len(parts) != 2 {
		return fmt.Errorf("expected 2 fields in topic item key, got: %d", len(parts))
	}

	ts, jid := parts[0], parts[1]

	prefix := prefixKeyDelayedQueue + "ts="
	if !strings.HasPrefix(ts, prefix) {
		return fmt.Errorf("wrong prefix in topic item key")
	}
	ts = strings.TrimPrefix(ts, prefix)

	prefix = "id="
	if !strings.HasPrefix(jid, prefix) {
		return fmt.Errorf("wrong prefix in topic item key")
	}
	jid = strings.TrimPrefix(jid, prefix)

	id, err := jobq.NewJobIDFromString(jid)
	if err != nil {
		return err
	}

	tsn, err := strconv.ParseUint(ts, 10, 64)
	if err != nil {
		return err
	}

	k.At = time.Unix(int64(tsn), 0)
	k.JobID = id

	return nil
}
