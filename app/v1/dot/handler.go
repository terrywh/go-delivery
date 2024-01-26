package dot

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"os"

	"github.com/libp2p/go-libp2p/core/network"
)

type HeaderCreateFile struct {
	Command   uint8
	Resv1     uint8
	NameSize uint16
	Chunk    uint16 // 断点续传分块
	Total    uint16
	DataSize uint64
}

const (
	CommandCreateFile      uint8 = 0x00
	CommandRemoveFile      uint8 = 0x01
	CommandCreateDirectory uint8 = 0x02
	CommandRemoveDirectory uint8 = 0x03
	CommandCreateFileChunk uint8 = 0x04 // 文件区块
	CommandCreateFileFinal uint8 = 0x05 // 文件结束
)

// StreamHandler ...
func StreamHandler(s network.Stream) {
	// RECVING ? 
	// SENDING ?
	r := bufio.NewReader(s)
	// w := bufio.NewWriter(s)
	var err  error

	header := make([]byte, 16)
	if _, err = io.ReadFull(r, header); err != nil {
		log.Println("failed to read header: ", err)
		return
	}
	switch header[0] {
	case CommandCreateFile:
		var desc HeaderCreateFile
		binary.Read(bytes.NewReader(header), binary.BigEndian, &desc)
		name := make([]byte, desc.NameSize)
		if _, err = io.ReadFull(r, name); err != nil {
			log.Println("failed to read name: ", err)
			return
		}
		file, _ := os.Create("/tmp/tmp.dat")
		io.CopyN(file, r, int64(desc.DataSize))
	}

}