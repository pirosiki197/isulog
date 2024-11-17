package internal

import (
	"bufio"
	"encoding/gob"
	"os"
	"sync"
	"time"
)

// TODO: add buffering?
type Recorder struct {
	filename string
	enc      *gob.Encoder
	mu       sync.Mutex
	records  []Record
}

func NewRecorder(filename string) *Recorder {
	r := &Recorder{
		filename: filename,
		records:  make([]Record, 1024),
	}
	r.reset()
	return r
}

func (r *Recorder) Save(record Record) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if err := r.enc.Encode(record); err != nil {
		if err := r.reset(); err != nil {
			return err
		}
	}
	return nil
}

func (r *Recorder) reset() error {
	f, err := os.Create(r.filename)
	if err != nil {
		return err
	}
	r.enc = gob.NewEncoder(bufio.NewWriter(f))
	return nil
}

// Record can be encoded into gob
type Record struct {
	Path         string
	Method       string
	StatusCode   int
	ResponseTime time.Duration
}
