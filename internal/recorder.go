package internal

import (
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
}

func NewRecorder(filename string) *Recorder {
	r := &Recorder{
		filename: filename,
	}
	r.reset()
	return r
}

func (r *Recorder) Save(record Record) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	err := r.enc.Encode(record)
	if err != nil {
		if err := r.reset(); err != nil {
			return err
		}
		return r.enc.Encode(record)
	}
	return nil
}

func (r *Recorder) reset() error {
	f, err := os.Create(r.filename)
	if err != nil {
		return err
	}
	r.enc = gob.NewEncoder(f)
	return nil
}

// Record can be encoded into gob
type Record struct {
	Path         string
	Method       string
	StatusCode   int
	ResponseTime time.Duration
}
