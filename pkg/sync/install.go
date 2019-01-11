//
// Copyright (c) Zhou Peng <p@ctriple.cn>
//
package sync

import (
	"log"
	"os/exec"
	"path"

	"github.com/ctriple/drbd/pkg/defs"
)

// InstallDriver installs ctriple.cn/drbd kubernetes flexvolume driver to host path
func InstallDriver() error {
	const (
		flexMnt   = "/flexmnt" // Mounted from hostPath '/usr/libexec/kubernetes/kubelet-plugins/volume/exec'
		newDriver = "/drbd"    // Newer flexvolume driver shipped with this container

		vendorDir = defs.Vendor + "~" + defs.Driver
		driverExe = defs.Driver
	)

	var (
		driverDir = path.Join(flexMnt, vendorDir)
		tmpDriver = path.Join(driverDir, "."+driverExe)
		dstDriver = path.Join(driverDir, driverExe)
	)

	if out, err := exec.Command("mkdir", "-p", driverDir).CombinedOutput(); err != nil {
		log.Println("mkdir -p", driverDir, string(out))
		return err
	}

	// NOTE: copy is not an atomic operation, so we first copy as a temp
	// file, and rename by atomic.
	if out, err := exec.Command("cp", "-p", newDriver, tmpDriver).CombinedOutput(); err != nil {
		log.Println("cp -p", newDriver, tmpDriver, string(out))
		return err
	}
	if out, err := exec.Command("mv", "-f", tmpDriver, dstDriver).CombinedOutput(); err != nil {
		log.Println("mv -f", tmpDriver, dstDriver, string(out))
		return err
	}

	return nil
}
