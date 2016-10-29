package db

import (
	"encoding/json"
	"os"
)

const (
	// Only user read write
	DBFilePerm = 0600
)

type DB struct {
	File string
	mem  map[string]interface{}
}

// NewDB create a new DB that will save to file
func NewDB(file string) *DB {
	return &DB{
		File: file,
		mem:  make(map[string]interface{}, 100),
	}
}

func (d *DB) Save() error {
	f, err := os.OpenFile(d.File, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, DBFilePerm)
	if err != nil {
		return err
	}

	dec := json.NewEncoder(f)
	err = dec.Encode(&d.mem)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) Load() error {
	f, err := os.OpenFile(d.File, os.O_RDONLY, DBFilePerm)
	if err != nil {
		return err
	}

	dec := json.NewDecoder(f)
	err = dec.Decode(&d.mem)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) Update(key string, value interface{}) {
	d.mem[key] = value
	d.Save()
}

func (d *DB) Get(key string) interface{} {
	v, ok := d.mem[key]
	if !ok {
		return nil
	}

	return v
}
