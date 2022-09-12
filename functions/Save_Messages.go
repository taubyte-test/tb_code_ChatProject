package lib

import (
	"fmt"

	"bitbucket.org/taubyte/go-sdk/database"
	"bitbucket.org/taubyte/go-sdk/event"
	"github.com/valyala/fastjson"
)

//export saveMessage
func saveMessage(e event.Event) uint32 {
	database, err := database.New("Chat_Database")
	if err != nil {
		return 1
	}

	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	msg, err := h.Headers().Get("message")
	if err != nil {
		return 1
	}

	user, err := h.Headers().Get("user")
	if err != nil {
		return 1
	}

	timestamp, err := h.Headers().Get("timestamp")
	if err != nil {
		return 1
	}

	secret, err := h.Headers().Get("secret")
	if err != nil {
		return 1
	}

	data := fmt.Sprintf(`{"msg": "%s", "user": "%s", "timestamp": %s}`, msg, user, timestamp)

	var entry string
	if secret == "" {
		entry = fmt.Sprintf("msg/all/%s", timestamp)
	} else {
		entry = fmt.Sprintf("msg/%s/%s", secret, timestamp)
	}

	v, err := fastjson.Parse(data)
	if err != nil {
		return 1
	}

	jsonObject, err := v.Object()
	if err != nil {
		return 1
	}

	err = database.Put(entry, []byte(jsonObject.String()))
	if err != nil {
		return 1
	}

	_, err = h.Write([]byte(jsonObject.String()))
	if err != nil {
		return 1
	}

	return 0
}
