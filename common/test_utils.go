package common

import "net"

func GetRandomPort() string {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}

	defer func() {
		err = listener.Close()
		if err != nil {
			panic(err)
		}
	}()
	_, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		panic(err)
	}

	return port
}
