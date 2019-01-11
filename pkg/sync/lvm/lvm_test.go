//
// Copyright (c) Zhou Peng <p@ctriple.cn>
//
package lvm

import (
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/ctriple/drbd/pkg/defs"
)

var (
	vg   = "centos"
	name = "test-sshlvm-disk"
	size = "500M"
	disk = path.Join("/dev", vg, name)
)

var skip = false

func init() {
	if _, err := exec.Command("vgdisplay", defs.DrbdDiskVG).CombinedOutput(); err != nil {
		skip = true
	}
}

func TestCreate(t *testing.T) {
	if skip {
		t.Skipf("lvm prerequisite does not meet!")
	}

	err := Create(vg, name, size)
	if err != nil {
		t.Fatal(err)
	}
}

// This pseudo testcase just show disk info to confirm creating succeed
func Test_show_disk_info(t *testing.T) {
	if skip {
		t.Skipf("lvm prerequisite does not meet!")
	}

	cmd := []string{"fdisk", "-l", disk}
	out, err := sshexec(strings.Join(cmd, " "))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(out)
}

func TestRemove(t *testing.T) {
	if skip {
		t.Skipf("lvm prerequisite does not meet!")
	}

	err := Remove(disk)
	if err != nil {
		t.Fatal(err)
	}
}
