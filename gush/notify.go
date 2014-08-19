package gush

import (
	"net/http"
)

func notify(uid string, msg string) string {
	u, ok := userMap[uid]
	if ok {
		u.msg <- MSG_PREFIX + msg
		return "OK"
	} else {
		return "FAIL"
	}
}

func regNotifyApi() {
	http.HandleFunc("/notify", func(rw http.ResponseWriter, request *http.Request) {
		uid := request.FormValue("uid")
		msg := request.FormValue("msg")
		mm := notify(uid, msg)
		rw.Write([]byte(mm))
	})

	http.ListenAndServe(":"+config.Port_notify, nil)
}
