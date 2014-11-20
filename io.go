package ttunnel

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

// readBytes reads a uint16 number of bytes to read from the reader,
// then reads that number of bytes into a buffer to return.  If the
// number of bytes is greater than `limit`, an error is returned.
func readBytes(conn io.Reader, limit uint16) ([]byte, error) {
	var length uint16
	if err := binary.Read(conn, binary.LittleEndian, &length); err != nil {
		return nil, err
	}

	if length > limit {
		return nil, fmt.Errorf("Data too long: %v bytes.", length)
	}

	buf := make([]byte, length)

	if _, err := conn.Read(buf); err != nil {
		return nil, err
	}

	return buf, nil
}

// writeBytes is the counterpart to readBytes. It writes the data
// length to the writer, then writes the data.
func writeBytes(conn io.Writer, data []byte) error {
	err := binary.Write(conn, binary.LittleEndian, uint16(len(data)))
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}

// Copy bytes from c1 to c2 and vice versa.
func copyDuplex(c1, c2 io.ReadWriteCloser) {
	go copy1(c1, c2)
	go copy1(c2, c1)
}

// copy1 performs a 1-way copy from src to dst. It closes both
// connections once the copy has completed or an error is encountered.
func copy1(dst, src io.ReadWriteCloser) {
	// Note: io.Copy won't return an EOF error when it encounters EOF.
	if _, err := io.Copy(dst, src); err != nil {
		log.Printf("Error copying connection:\n    %v\n", err)
	}
	dst.Close()
	src.Close()
}
