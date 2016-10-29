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

type TestDBData struct {
	Key string `json:"key"`
}

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

	d.mem["test"] = TestDBData{
		Key: "value",
	}

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

func TestDBUpdate(t *testing.T) {
	d := NewDB("")

	d.Update("test", TestDBData{
		Key: "value",
	})

	v, ok := d.mem["test"]
	if !ok {
		t.Fatalf("d.mem[\"test\"] was not present")
	}

	switch td := v.(type) {
	case TestDBData:
		if td.Key != "value" {
			t.Fatalf("d.mem[\"test\"].Key was %q should hve been \"value\"", td.Key)
		}
		break
	default:
		t.Fatalf("d.mem[\"test\"] was not of type TestDBData")
		break
	}
}

func TestDBGet(t *testing.T) {
	d := NewDB("")

	d.mem["test"] = TestDBData{
		Key: "value",
	}

	v := d.Get("test")

	if v == nil {
		t.Fatalf("d.mem[\"test\"] was nil")
	}

	switch td := v.(type) {
	case TestDBData:
		if td.Key != "value" {
			t.Fatalf("d.mem[\"test\"].Key was %q should hve been \"value\"", td.Key)
		}
		break
	default:
		t.Fatalf("d.mem[\"test\"] was not of type TestDBData")
		break
	}
}
