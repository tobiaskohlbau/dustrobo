package wav

import (
	"io"
	"log"

	"github.com/pkg/errors"
)

type Decoder struct {
	chunkID   []byte
	chunkSize uint32
	riffType  []byte
	src       io.ReadCloser
}

func NewDecoder(r io.ReadCloser) (*Decoder, error) {
	buf := make([]byte, 4)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, errors.Wrap(err, "failed to read wav header")
	}

	log.Printf("%x", buf)

	// var chunkSize uint32
	// err := binary.Read(bytes.NewReader(buf[4:8]), binary.LittleEndian, &chunkSize)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "failed to read chunkSize")
	// }

	// if !bytes.Equal(buf[0:4], []byte("RIFF")) {
	// 	return nil, errors.Errorf("invalid chunk id %x", buf[0:4])
	// }

	// if !bytes.Equal(buf[8:12], []byte("WAVE")) {
	// 	return nil, errors.New("invalid riff type")
	// }

	// size := int64(buf[4]) | int64(buf[5])<<8 | int64(buf[6])<<16 | int64(buf[7])<<24
	// var size int64
	// if err := binary.Read(bytes.NewReader(buf[4:8]), binary.LittleEndian, &size); err != nil {
	// 	return nil, errors.Wrap(err, "failed to read fmt header size")
	// }
	// log.Println(size)

	// for {
	// 	buf := make([]byte, 8)
	// 	if _, err := io.ReadFull(r, buff); err != nil {
	// 		return nil, errors.Wrap(err, "failed to read full chunk")
	// 	}

	// }

	return &Decoder{
		src: r,
		// chunkID:   buf[0:3],
		// chunkSize: chunkSize,
		// riffType:  buf[8:11],
	}, nil
}

func (d *Decoder) Close() error {
	return d.src.Close()
}

func (d *Decoder) Read(buf []byte) (int, error) {
	return d.src.Read(buf)
}
