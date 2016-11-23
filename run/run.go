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

	"github.com/pdxjohnny/funkmasta/funkmasta"
)

const (
	tempDirPrefix = "funkmasta_run"
)

type Service struct {
	*funkmasta.Service
	tempDir          string
	envSetupFileName string
	payloadFileName  string
	payloadEnv       []string
}

func NewService(gs *funkmasta.Service) *Service {
	s := &Service{
		Service:          gs,
		tempDir:          "",
		envSetupFileName: fmt.Sprintf("%x", md5.Sum([]byte(gs.EnvSetup))),
		payloadFileName:  fmt.Sprintf("%x", md5.Sum([]byte(gs.Payload))),
	}
	return s
}

func envArrayToMap(a []string) map[string]string {
	m := make(map[string]string, 10)
	for _, i := range a {
		j := strings.SplitN(i, "=", 2)
		if len(j) == 2 {
			m[j[0]] = j[1]
		}
	}
	return m
}

func envMapToArray(m map[string]string) []string {
	a := make([]string, 10)
	for k, v := range m {
		a = append(a, k+"="+v)
	}
	return a
}

func (s *Service) RunValidate() error {
	if len(s.EnvSetup) < 1 {
		return fmt.Errorf("EnvSetup is empty")
	}

	if len(s.Payload) < 1 {
		return fmt.Errorf("Payload is empty")
	}

	return nil
}

func (s *Service) RunSetup() error {
	// Make sure we have everything we need
	err := s.RunValidate()
	if err != nil {
		return err
	}

	// Create a temporary directory to run this in
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

	s.payloadEnv = strings.Split(outputArray[1], "\n")

	found := false
	for _, i := range s.payloadEnv {
		j := strings.SplitN(i, "=", 2)
		if j[0] == "RUN_ENVSETUP_TEST_WORKING" {
			found = true
		}
	}

	if !found {
		return fmt.Errorf(
			"RUN_ENVSETUP_TEST_WORKING was not found in environment\n%v\n",
			s.payloadEnv,
		)
	}

	return nil
}

func (s *Service) RunPayload(r *funkmasta.Request) error {
	// Run Payload withe Env as r.Env send Stdout and Stderr to r.Output
	if s.payloadEnv == nil {
		return fmt.Errorf("payloadEnv is nil, has RunEnvSetup been called yet?")
	}

	// Run Payload with Env as r.Env
	// The env of bash after this becomes the env of Payload
	cmd := exec.Command(filepath.Join(s.tempDir, s.payloadFileName))
	// Make the env the sent in env
	env := envArrayToMap(r.Env)
	// Overwite with the EnvSetup env
	payloadEnv := envArrayToMap(s.payloadEnv)
	for k, v := range payloadEnv {
		env[k] = v
	}
	cmd.Env = envMapToArray(env)
	// Capture stdout and err
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("Getting Payload stdin: %s", err.Error())
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Getting Payload stdout: %s", err.Error())
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("Getting Payload stderr: %s", err.Error())
	}
	outputReader := io.MultiReader(stdout, stderr)

	// Start the shell
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("Starting Payload: %s", err.Error())
	}

	// Send the input to the command
	go func() {
		io.Copy(stdin, r.Body)
	}()
	// Send the output to the client
	io.Copy(r.Output, outputReader)

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("Payload Failed: %s", err.Error())
	}

	return nil
}
