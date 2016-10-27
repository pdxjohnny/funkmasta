package db

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

const (
	TestDBTempFilePrefix  = "TestDBTempFilePrefix"
	TestDBTempFileContent = `{"test":{"key":"value"}}`
)

func TestDBSave(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", TestDBTempFilePrefix)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	d := NewDB(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	d.mem["test"] = make(map[string]string, 1)
	(d.mem["test"].(map[string]string))["key"] = "value"

	d.Save()
	if err != nil {
		t.Fatal(err)
	}

	j, err := ioutil.ReadAll(tmpfile)
	if err != nil {
		t.Fatal(err)
	}
	b := string(j[:len(j)-1])

	err = tmpfile.Close()
	if err != nil {
		t.Fatal(err)
	}

	if b != TestDBTempFileContent {
		t.Fatalf("Expected: %v, got: %v", []byte(TestDBTempFileContent), []byte(b))
	}
}

func TestDBLoad(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", TestDBTempFilePrefix)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	d := NewDB(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	_, err = tmpfile.Write([]byte(TestDBTempFileContent))
	if err != nil {
		t.Fatal(err)
	}

	err = tmpfile.Close()
	if err != nil {
		t.Fatal(err)
	}

	d.Load()
	if err != nil {
		t.Fatal(err)
	}

	j, err := json.Marshal(d.mem)
	if err != nil {
		t.Fatal(err)
	}
	b := string(j)

	if b != TestDBTempFileContent {
		t.Fatalf("Expected: %v, got: %v", []byte(TestDBTempFileContent), []byte(b))
	}
}
