package gush

import (
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	AUTH          = "_A"
	HB            = "_H"
	END           = "_E"
	READ_TIMEOUT  = 6
	WRITE_TIMEOUT = 3
)

var userMap = make(map[string]*userChannel)

type userChannel struct {
	uid string
	msg chan string
}

func NewGush() {
	ln, err := net.Listen("tcp", ":8888")
	if err != nil {
		Logger.Error(err)
	}

	go regNotifyApi()

	for {
		conn, err := ln.Accept()
		if err != nil {
			Logger.Error(err)
			continue
		}

		channel := make(chan string)

		uc := &userChannel{"", channel}

		go readConn(conn, uc)
		go wirteConn(conn, uc)
	}
}

func readConn(c net.Conn, uc *userChannel) {
	defer func() {
		uc.msg <- END
	}()

	buf := make([]byte, 1024)

	for {
		c.SetReadDeadline(time.Now().Add(READ_TIMEOUT * time.Second))
		n, err := c.Read(buf)

		if err != nil {
			if err != io.EOF {
				Logger.Error("Read error: %s", err)
			}
			Logger.Error("conn error: %s", err)
			break
		}

		if n > 0 {
			m := string(buf[:n])
			routeMsg(m, uc)
		}
	}
}

func wirteConn(c net.Conn, uc *userChannel) {
	defer func() {
		c.Close()
		close(uc.msg)
		delete(userMap, uc.uid)
	}()

	for {
		msg := <-uc.msg
		if msg == END {
			Logger.Error("conn closed.")
			break
		}

		c.SetWriteDeadline(time.Now().Add(WRITE_TIMEOUT * time.Second))
		_, err := c.Write([]byte(msg))

		if err != nil {
			Logger.Error("Write error: %s", err)
			break
		}
	}
}

func notify(uid string, msg string) string {
	u, ok := userMap[uid]
	if ok {
		u.msg <- msg
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

	http.ListenAndServe(":8080", nil)
}

func routeMsg(request string, uc *userChannel) {
	if strings.HasPrefix(request, AUTH) {
		auth(strings.Replace(request, AUTH, "", 1), uc)
	} else if request == HB {
		heartbeat(uc)
	} else {
		if uc.uid != "" {
			//do something
		}
	}
}

func auth(request string, uc *userChannel) {
	// do auth check
	uc.uid = request
	userMap[request] = uc

	uc.msg <- "OK"
}

func heartbeat(uc *userChannel) {
	uc.msg <- "OK"
}
