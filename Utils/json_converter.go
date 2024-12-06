package Utils

import (
	"encoding/json"
)

func MustDeMarshallBytesToMsg(jsonData *[]byte) Message {
	var msg Message
	err := json.Unmarshal(*jsonData, &msg)
	if err != nil {
		panic("Failed to unmarshall from Json to Message")
	}	
	return msg
}

func MustMarshalMsgToBytes(msg *Message) []byte {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		panic("Failed to marshall from Message to Json")
	}
	return jsonData
}
