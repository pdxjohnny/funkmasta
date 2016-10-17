package run

import (
	"crypto/md5"
	"fmt"
	"io"
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
	s.tempDir, err = ioutil.TempDir("", tempDirPrefix)
	if err != nil {
		return fmt.Errorf("Creating TempDir: %s", err.Error())
	}

	// Create the EnvSetup file
	err = ioutil.WriteFile(
		filepath.Join(s.tempDir, s.envSetupFileName),
		[]byte(s.EnvSetup),
		0700,
	)
	if err != nil {
		return fmt.Errorf("Writing EnvFile: %s", err.Error())
	}

	// Create the Payload file
	err = ioutil.WriteFile(
		filepath.Join(s.tempDir, s.payloadFileName),
		[]byte(s.Payload),
		0700,
	)
	if err != nil {
		return fmt.Errorf("Writing Payload: %s", err.Error())
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
	// Make the env empty
	cmd.Env = make([]string, 10)
	// Capture stdin out and err
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("Getting EnvSetup stdin: %s", err.Error())
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Getting EnvSetup stdout: %s", err.Error())
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("Getting EnvSetup stderr: %s", err.Error())
	}
	outputReader := io.MultiReader(stdout, stderr)

	// Start the shell
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("Starting EnvSetup: %s", err.Error())
	}

	// Source the env setup file and capture the env after
	delim := rand.Int63()
	_, err = stdin.Write([]byte(fmt.Sprintf(
		"source %s\necho %d\nexport RUN_ENVSETUP_TEST_WORKING=true\nenv\nexit 0\n",
		filepath.Join(s.tempDir, s.envSetupFileName),
		delim,
	)))
	if err != nil {
		return fmt.Errorf("Running EnvSetup: %s", err.Error())
	}

	outputBytes, err := ioutil.ReadAll(outputReader)
	if err != nil {
		return fmt.Errorf("Could not get EnvSetup output: %s", err.Error())
	}
	output := string(outputBytes)

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("EnvSetup Failed: %s:\n%s", err.Error(), output)
	}

	outputArray := strings.Split(string(output), fmt.Sprintf("%d", delim))
	if len(outputArray) != 2 {
		return fmt.Errorf("Could not get environment for %s", s.Name)
	}

	cmd.Env = strings.Split(outputArray[1], "\n")

	found := false
	for _, i := range cmd.Env {
		j := strings.SplitN(i, "=", 2)
		if j[0] == "RUN_ENVSETUP_TEST_WORKING" {
			found = true
		}
	}

	if !found {
		return fmt.Errorf(
			"RUN_ENVSETUP_TEST_WORKING was not found in environment\n%v\n",
			cmd.Env,
		)
	}

	return nil
}

func (s *Service) RunPayload(r *getfunky.Request) error {
	// Run Payload withe Env as r.Env send Stdout and Stderr to r.Output

	return nil
}
