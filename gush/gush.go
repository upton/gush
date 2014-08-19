package gush

import (
	"net"
	"strings"
	"time"
)

const (
	AUTH       = "_A:"
	HB         = "_H:"
	END        = "__END"
	MSG_PREFIX = "_M:"
	NEW_LINE   = "\n"
)

var userMap = make(map[string]*UserChannel)

type UserChannel struct {
	uid string
	msg chan string
}

func Run() {
	ln, err := net.Listen("tcp", ":"+config.Port_tcp)
	if err != nil {
		Logger.Error(err)
	}

	go regNotifyApi()

	Logger.Info("start success")
	for {
		conn, err := ln.Accept()
		if err != nil {
			Logger.Error(err)
			continue
		}

		channel := make(chan string)

		uc := &UserChannel{"", channel}

		go readConn(conn, uc)
		go wirteConn(conn, uc)
	}
}

func readConn(c net.Conn, uc *UserChannel) {
	defer func() {
		uc.msg <- END
	}()

	buf := make([]byte, 1024)

	for {
		c.SetReadDeadline(time.Now().Add(time.Duration(config.Read_timeout) * time.Second))
		n, err := c.Read(buf)

		if err != nil {
			Logger.Error("conn error: ", err)
			break
		}

		if n > 0 {
			m := string(buf[:n])
			routeMsg(m, uc)
		}
	}
}

func wirteConn(c net.Conn, uc *UserChannel) {
	defer func() {
		c.Close()
		close(uc.msg)
		delete(userMap, uc.uid)
	}()

	for {
		msg := <-uc.msg
		if msg == END {
			Logger.Warn("conn closed.")
			break
		}

		c.SetWriteDeadline(time.Now().Add(time.Duration(config.Write_timeout) * time.Second))
		_, err := c.Write([]byte(msg + NEW_LINE))

		if err != nil {
			Logger.Error("Write error: ", err)
			break
		}
	}
}

func routeMsg(request string, uc *UserChannel) {
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

func auth(request string, uc *UserChannel) {
	// do auth check
	uc.uid = request
	userMap[request] = uc

	uc.msg <- AUTH + "OK"
}

func heartbeat(uc *UserChannel) {
	uc.msg <- HB + "OK"
}
