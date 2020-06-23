package databasetest

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"testing"
)

// Container type holds the information
// of each docker container spun up for testing.
// All tests are done inside a docker container.Container
type Container struct {
	ID   string
	Host string
}

// StartContainer spins up a postgresql container to use as testenv
func StartContainer(t *testing.T) *Container {
	t.Helper()

	cmd := exec.Command("docker", "run", "-P", "-d", "postgres:11.1-alpine")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("Could not start test container: %v", err)
	}

	id := out.String()[:12]
	t.Log("Test db continer id is", id)

	cmd = exec.Command("docker", "inspect", id)
	out.Reset()
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("Could not inspect container %s: %v", id, err)
	}

	// Creating a data structure to hold info from the output
	// of the docker inspect command above
	var doc []struct {
		NetworkSettings struct {
			Port struct {
				TCP5432 []struct {
					HostIP   string `json:"HostIp"`
					HostPort string `json:"HostPort"`
				} `json:"5432/tcp"`
			} `json:"Ports"`
		} `json:"NetworkSettings"`
	}
	if err := json.Unmarshal(out.Bytes(), &doc); err != nil {
		t.Fatalf("Could not decode json %s: %v", id, err)
	}

	network := doc[0].NetworkSettings.Port.TCP5432[0]

	c := Container{
		ID:   id,
		Host: network.HostIP + ":" + network.HostPort,
	}

	t.Log("Test DB host is", c.Host)

	return &c
}

// StopContainer stops and removes the specified container.
func StopContainer(t *testing.T, c *Container) {
	t.Helper()

	if err := exec.Command("docker", "stop", c.ID).Run(); err != nil {
		t.Fatalf("Could not stop container %s: %v", c.ID, err)
	}
	t.Log("Stopped:", c.ID)

	if err := exec.Command("docker", "rm", c.ID, "-v").Run(); err != nil {
		t.Fatalf("Could not remove container %s: %v", c.ID, err)
	}
	t.Log("Removed:", c.ID)
}

// DumpContainerLogs runs "docker logs" and sends all logs to the test logger
func DumpContainerLogs(t *testing.T, c *Container) {
	t.Helper()

	out, err := exec.Command("docker", "logs", c.ID).CombinedOutput()
	if err != nil {
		t.Fatalf("Could not log container %s: %v", c.ID, err)
	}
	t.Logf("Logs for %s\n%s:", c.ID, out)
}
