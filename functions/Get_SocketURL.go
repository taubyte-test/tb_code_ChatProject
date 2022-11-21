package lib

import (
	"bitbucket.org/taubyte/go-sdk/event"
	"bitbucket.org/taubyte/go-sdk/globals/u32"
	"bitbucket.org/taubyte/go-sdk/pubsub"
)

//export getsocketurl
func getsocketurl(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	err = _getsocketurl(h)
	if err != nil {
		h.Write([]byte(err.Error()))
		h.Return(500)

		return 1
	}

	return 0
}

func _getsocketurl(h event.HttpEvent) error {
	channel, err := pubsub.Channel("chatChannel")
	if err != nil {
		return err
	}

	url, err := channel.WebSocket().Url()
	if err != nil {
		return err
	}

	u, err := u32.GetOrCreate("chatUsers")
	if err == nil {

		u.Set(u.Value() + 1)
	}

	_, err = h.Write([]byte(url.Path))
	if err != nil {
		return err
	}

	return nil
}
