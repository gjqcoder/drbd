//
// Copyright (c) Zhou Peng <p@ctriple.cn>
//
package fs

import (
	"log"
	"os/exec"
	"strings"
)

// Mount mount device to mount point specified by path, if mount point does
// not exists, it will create it.
func Mount(device, path string) error {
	// Make sure mount point exists
	if out, err := exec.Command("mkdir", "-p", path).CombinedOutput(); err != nil {
		log.Println("mkdir -p", path, string(out))
		return err
	}

	if out, err := exec.Command("mount", device, path).CombinedOutput(); err != nil {
		log.Println("mount", device, path, string(out))
		return err
	}

	return nil
}

// Unmount unmount mount point specified by path
func Unmount(path string) error {
	// Mount pont must be directory
	if out, err := exec.Command("test", "-d", path).CombinedOutput(); err != nil {
		log.Println("test -d", path, string(out))
		return err
	}

	// Mount point isn't mounted
	if out, err := exec.Command("findmnt", "-f", path).CombinedOutput(); err != nil {
		log.Println("findmnt -f", path, string(out))
	}

	if out, err := exec.Command("umount", path).CombinedOutput(); err != nil {
		log.Println("umount", path, string(out))
		return err
	}

	return nil
}

// Format format device to filesystem type specified by fsType, if device is
// already has desired filesystem type, it will do nothing and return
// successfully.
func Format(device, fsType string) error {
	needFormat := true

	// Check if block device already has desired filesystem type
	if out, err := exec.Command("blkid", "-o", "udev", device).CombinedOutput(); err == nil {
		const FSKEY = "ID_FS_TYPE"

		fields := strings.Fields(string(out))
		for _, pair := range fields {
			p := strings.Split(pair, "=")
			if len(p) < 2 {
				continue
			}
			if FSKEY == p[0] && fsType == p[1] {
				// Already has desired fielsystem type
				needFormat = false
			}
		}
	}

	if needFormat {
		if out, err := exec.Command("mkfs", "-t", fsType, device).CombinedOutput(); err != nil {
			log.Println("mkfs -t", fsType, device, string(out))
			return err
		}
	}

	return nil
}
