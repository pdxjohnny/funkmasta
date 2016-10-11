package run

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pdxjohnny/getfunky/getfunky"
)

const (
	tempDirPrefix = "getfunky_run"
)

type Service struct {
	*getfunky.Service
	tempDir          string
	envSetupFileName string
	payloadFileName  string
}

func NewService(gs *getfunky.Service) *Service {
	s := &Service{
		Service:          gs,
		tempDir:          "",
		envSetupFileName: fmt.Sprintf("%x", md5.Sum([]byte(gs.EnvSetup))),
		payloadFileName:  fmt.Sprintf("%x", md5.Sum([]byte(gs.Payload))),
	}
	return s
}

func (s *Service) RunSetup() error {
	// Create a temporary directory to run this in
	var err error
	s.tempDir, err = ioutil.TempDir(s.envSetupFileName+s.payloadFileName, tempDirPrefix)
	if err != nil {
		return err
	}

	// Create the EnvSetup file
	err = ioutil.WriteFile(
		filepath.Join(s.tempDir, s.envSetupFileName),
		[]byte(s.EnvSetup),
		0700,
	)
	if err != nil {
		return err
	}

	// Create the Payload file
	err = ioutil.WriteFile(
		filepath.Join(s.tempDir, s.payloadFileName),
		[]byte(s.Payload),
		0700,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) RunTeardown() error {
	return os.RemoveAll(s.tempDir)
}

func (s *Service) RunEnvSetup() error {
	// Run EnvSetup with Env as r.Env
	// The env of bash after this becomes the env of Payload
	cmd := exec.Command("bash", "-e")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}

	delim := rand.Int63()
	_, err = stdin.Write([]byte(fmt.Sprintf("source %s\necho %ld\nenv\nexit 0\n", s.envSetupFileName, delim)))
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	outputArray := strings.Split(string(output), fmt.Sprintf("%ld", delim))
	if len(outputArray) != 2 {
		return fmt.Errorf("Could not get environment for %s", s.Name)
	}

	fmt.Println("source", outputArray[0])
	fmt.Println("env", outputArray[1])

	return nil
}

func (s *Service) RunPayload(r *getfunky.Request) error {
	// Run Payload withe Env as r.Env send Stdout and Stderr to r.Output

	return nil
}
