package record

import "encoding/json"

func decodeEvent(data []byte) (Event, error) {
	var e Event
	if len(data) >= 8 {
		data = data[8:]
	}
	if err := json.Unmarshal(data, &e); err != nil {
		return Event{}, err
	}
	return e, nil
}

func marshalEvent(e Event) ([]byte, error) {
	return json.Marshal(e)
}
