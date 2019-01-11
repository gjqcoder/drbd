//
// Copyright (c) Zhou Peng <p@ctriple.cn>
//
package drbdadm

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// Primary promote this drbd node as primary role
func Primary(resName string) error {
	out, err := exec.Command("drbdadm", "primary", resName, "--force").CombinedOutput()
	if err != nil {
		log.Println("drbdadm primary", resName, string(out))
		return err
	}

	return nil
}

// Secondary demote this drbd node as secondary role
func Secondary(resName string) error {
	out, err := exec.Command("drbdadm", "secondary", resName).CombinedOutput()
	if err != nil {
		log.Println("drbdadm secondary", resName, string(out))
		return err
	}

	return nil
}

// CreateMD create metadata on this newly drbd resource backing physical disk
func CreateMD(resName string) error {
	out, err := exec.Command("drbdadm", "create-md", resName, "--force").CombinedOutput()
	if err != nil {
		log.Println("drbdadm create-md", resName, string(out))
		return err
	}

	return nil
}

// Up makes this resource on the current drbd node start serving
func Up(resName string) error {
	out, err := exec.Command("drbdadm", "up", resName).CombinedOutput()
	if err != nil {
		log.Println("drbdadm up", resName, string(out))
		return err
	}

	return nil
}

// Down makes this resource on the current drbd node stop serving
func Down(resName string) error {
	out, err := exec.Command("drbdadm", "down", resName).CombinedOutput()
	if err != nil {
		log.Println("drbdadm down", resName, string(out))
		return err
	}

	return nil
}

// Adjust makes new resource config take effect, it basically equals down and up
func Adjust(resName string) error {
	out, err := exec.Command("drbdadm", "adjust", resName).CombinedOutput()
	if err != nil {
		log.Println("drbdadm adjust", resName, string(out))
		return err
	}

	return nil
}

// ShResources returns all resource names on this drbd node
func ShResources() ([]string, error) {
	out, err := exec.Command("drbdadm", "sh-resources").CombinedOutput()
	if err != nil {
		log.Println("drbdadm sh-resources", string(out))
		return []string{}, err
	}

	return strings.Fields(string(out)), nil
}

// ShResource returns true if this resource on the current host
func ShResource(resName string) bool {
	resNames, err := ShResources()
	if err != nil {
		return false
	}

	for _, res := range resNames {
		if res == resName {
			return true
		}
	}

	return false
}

// ShDev returns drbd resource virtual device
func ShDev(resName string) (string, error) {
	out, err := exec.Command("drbdadm", "sh-dev", resName).CombinedOutput()
	if err != nil {
		log.Println("drbdadm sh-dev", resName, string(out))
		return "", err
	}

	// NOTE: command output has superfluous whitespace
	device := strings.TrimSpace(string(out))

	return device, nil
}

// ShLlDev returns drbd resource backing physical disk
func ShLlDev(resName string) (string, error) {
	out, err := exec.Command("drbdadm", "sh-ll-dev", resName).CombinedOutput()
	if err != nil {
		log.Println("drbdadm sh-ll-dev", resName, string(out))
		return "", err
	}

	// NOTE: command output has superfluous whitespace
	disk := strings.TrimSpace(string(out))

	return disk, nil
}

// ResByMntDir try to find drbd resource name by mount point
func ResByMntDir(mountDir string) (string, error) {
	device, err := devByMntDir(mountDir)
	if err != nil {
		return "", err
	}

	return resByDev(device)
}

// devByMntDir try to find drbd virtual device by mount point
func devByMntDir(mountDir string) (string, error) {
	// Mount point must be directory
	if out, err := exec.Command("test", "-d", mountDir).CombinedOutput(); err != nil {
		log.Println("test -d", mountDir, string(out))
		return "", err
	}

	out, err := exec.Command("findmnt", "-f", "-n", "--output", "SOURCE", mountDir).CombinedOutput()
	if err != nil {
		log.Println("findmnt -f -n --output SOURCE", mountDir, string(out))
		return "", err
	}

	// NOTE: command output has superfluous whitespace
	device := strings.TrimSpace(string(out))

	return device, nil
}

// resByDev try to find drbd resource name by drbd device
func resByDev(device string) (string, error) {
	resNames, err := ShResources()
	if err != nil {
		return "", err
	}

	for _, resName := range resNames {
		deviceName, err := ShDev(resName)
		if err != nil {
			continue
		}
		if deviceName == device {
			return resName, nil
		}
	}

	return "", fmt.Errorf("device: %s has no drbd resource related.", device)
}
