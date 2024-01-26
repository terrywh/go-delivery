package api

import (
	"bufio"
	"log"

	"github.com/libp2p/go-libp2p/core/network"
)

// StreamHandler ...
func StreamHandler(s network.Stream) {
	r := bufio.NewReader(s)
	w := bufio.NewWriter(s)
	var err error
	for {
		data, err := r.ReadString(byte('\n'))
		if err != nil {
			log.Println("failed to read: ", err)
			break
		}
		w.WriteString(data)
		w.Flush()
	}
	log.Println("closed: ", err)
}