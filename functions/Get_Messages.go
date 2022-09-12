package lib

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"bitbucket.org/taubyte/go-sdk/database"
	"bitbucket.org/taubyte/go-sdk/event"
	"github.com/valyala/fastjson"
)

//export getMessages
func getMessages(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	err = _getMessages(h)
	if err != nil {
		h.Write([]byte(fmt.Sprintf("Failed: %s", err.Error())))
		return 1
	}

	return 0
}

func _getMessages(h event.HttpEvent) error {
	var secret string
	_database, err := database.New("Chat_Database")
	if err != nil {
		return err
	}

	secret, err = h.Headers().Get("secret")
	if err != nil {
		return err
	}

	var listStatement string
	if secret == "" {
		listStatement = fmt.Sprintf("msg/all")
	} else {
		listStatement = fmt.Sprintf("msg/%s", secret)
	}

	keys, err := _database.List(listStatement)
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		_, err = h.Write([]byte("Currently No Messages Stored"))
		if err != nil {
			return err
		}
		return nil
	}

	// Sort msg timestamps increasingly
	times := make([]int64, 0)
	for _, key := range keys {
		_trim := key[strings.LastIndex(key, "/")+1:]
		time, err := strconv.ParseInt(_trim, 10, 64)
		if err != nil {
			return err
		}

		times = append(times, time)
	}
	// Increasing order
	sort.Slice(times, func(i, j int) bool {
		return times[i] < times[j]
	})

	var messages []string
	var data string
	for idx, time := range times {
		var getStatement string
		if secret == "" {
			getStatement = fmt.Sprintf("msg/all/%d", time)
		} else {
			getStatement = fmt.Sprintf("msg/%s/%d", secret, time)
		}
		dataBytes, err := _database.Get(getStatement)
		if err != nil {
			return err
		}

		jsonParse, err := fastjson.Parse(string(dataBytes))
		if err != nil {
			return err
		}

		// Detect for last loop
		if idx == len(times)-1 {
			data = fmt.Sprintf(`{"msg": %s, "user": %s, "timestamp": %s}`, jsonParse.Get("msg"), jsonParse.Get("user"), jsonParse.Get("timestamp"))
		} else {
			data = fmt.Sprintf(`{"msg": %s, "user": %s, "timestamp": %s},`, jsonParse.Get("msg"), jsonParse.Get("user"), jsonParse.Get("timestamp"))
		}

		messages = append(messages, data)
	}

	formatData := fmt.Sprintf("%v", messages)

	jsonParse, err := fastjson.Parse(formatData)
	if err != nil {
		return err
	}

	jsonObject := jsonParse.String()
	if err != nil {
		return err
	}

	_, err = h.Write([]byte(jsonObject))
	if err != nil {
		return err
	}

	return nil
}
