//
// Copyright (c) Zhou Peng <p@ctriple.cn>
//
package lvm

import (
	"fmt"
	"os/exec"
	"strings"
)

func Create(vg, name string, sizeMb string) error {
	sshcmd := []string{"lvcreate", "--name", name, "--size", sizeMb, vg}

	if _, err := sshexec(strings.Join(sshcmd, " ")); err != nil {
		return err
	}

	return nil
}

func Remove(diskPath string) error {
	// No disk to remove
	lscmd := []string{"stat", diskPath}
	if _, err := sshexec(strings.Join(lscmd, " ")); err != nil {
		return nil
	}

	rmcmd := []string{"lvremove", "-f", diskPath}
	if _, err := sshexec(strings.Join(rmcmd, " ")); err != nil {
		return err
	}

	return nil
}

// FIXME: This is a technical compromise
//
// Since lvm utility needs to access host's /sys and /proc, which makes it not
// possible to get lvm utility work right inside a container. So we ssh to host
// to workaround this limitation.
func sshexec(cmd string) (string, error) {
	out, err := exec.Command("ssh", "-T", "root@localhost", cmd).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ssh:%s error:%s", string(out), err)
	}
	return string(out), nil
}
