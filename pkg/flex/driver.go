//
// Copyright (c) Zhou Peng <p@ctriple.cn>
//
package flex

import (
	"encoding/json"
	"fmt"

	"github.com/ctriple/drbd/pkg/drbdadm"
	"github.com/ctriple/drbd/pkg/flex/fs"
)

func DriverRun(action string, args []string) exitCode {
	switch action {
	case "init":
		return doInit()
	case "attach":
		return doAttach(args)
	case "detach":
		return doDetach(args)
	case "waitforattach":
		return doWaitForAttach(args)
	case "isattached":
		return doIsAttached(args)
	case "mountdevice":
		return doMountDevice(args)
	case "unmountdeivce":
		return doUnmountDevice(args)
	case "mount":
		return doMount(args)
	case "unmount":
		return doUnmount(args)
	default:
		stdoutJson(errorNotSupported)
		return ExitFailure
	}
}

// Init:
//
// Initializes the driver. Called during Kubelet & Controller manager
// initialization. On success, the function returns a capabilities map showing
// whether each FlexVolume capabilities is supported by the driver. Current
//
// capabilities:
//
//   `attach` -- a boolean filed indicating whether the driver requires attach
//               and detach operations. This field is *required*, although for
//               backward-compatibility the default value is set to `true`, i.e.
//               requires attach and detach. See specificition for the capabilities
//               map format.
//
// <driver executable> init
//
func doInit() exitCode {
	echo := initCallEcho{
		callEcho: callEcho{
			Status:  StatusSuccess,
			Message: "ctriple.cn/drbd driver initialization ok",
		},
		Capabilities: initCallCapabilities{
			Attach: false, // NOTE: we don't support controller manager attach/detach
		},
	}
	stdoutJson(echo)

	return ExitSuccess
}

// Attach:
//
// Attach the volume specified by the given spec on the given node. On success,
// returns the device path where the device is attached on the node. Called from
// Controller Manager.
//
// This call-out does not pass "secrets" specified in Flexvolume spec. If your
// driver requires secrets, do not implement this call-out and instead use
// "mount" call-out and implement attach and mount in that call-out.
//
// <driver executable> attach <json options> <node name>
//
func doAttach(args []string) exitCode {
	stdoutJson(errorNotSupported)
	return ExitFailure
}

// Detach:
//
// Detach the volume form the node. Called from Controller Manager.
//
// <driver executable> detach <mount device> <node name>
//
func doDetach(args []string) exitCode {
	stdoutJson(errorNotSupported)
	return ExitFailure
}

// Wait for attach:
//
// Wait for the volume to be attached on the remote node. On success, the path
// to the device is returned. Called from Controller Manager. The timeout should
// be 10m (based on
// https://git.k8s.io/kubernetes/pkg/kubelet/volumemanager/volume_manager.go#L88)
//
// <driver executable> waitforattach <mount device> <json options>
//
func doWaitForAttach(args []string) exitCode {
	stdoutJson(errorNotSupported)
	return ExitFailure
}

// Volume is Attached:
//
// Check the volume is attached on the node. Called from Controller Manager.
//
// <driver executable> isattached <json options> <node name>
//
func doIsAttached(args []string) exitCode {
	stdoutJson(errorNotSupported)
	return ExitFailure
}

// Mount device:
//
// Mount device mounts the device to a global path with inidvidual pods can then
// bind mount. Called only from Kubelet.
//
// Thsi call-out does not pass "secrets" specified in Flexvolume spec. If your
// driver requires secrets, do not implement this call-out and instead use
// "mount" call-out and implement attach and mount in that call-out.
//
// <driver executable> mountdevice <mount dir> <mount device> <json options>
//
func doMountDevice(args []string) exitCode {
	stdoutJson(errorNotSupported)
	return ExitFailure
}

// Unmount device:
//
// Unmounts the global mount for the device. This is called once all bind mounts
// have been unmounted. Called only from Kubelet.
//
// In addition to the user-specified options and default JSON options, the
// following options capturing information about the pod are passed through and
// generated automatically.
//
//   kubernetes.io/pod.name
//   kubernetes.io/pod.namespace
//   kubernetes.io/pod.uid
//   kubernetes.io/serviceAccount.name
//
// <driver executable> unmountdeivce <mount device>
//
func doUnmountDevice(args []string) exitCode {
	stdoutJson(errorNotSupported)
	return ExitFailure
}

// Mount:
//
// Mount the volume at the mount dir. This call-out defaults to bind mount for
// drivers which implement attach & mount-device call-outs. Called only from
// Kubelet.
//
// <driver executable> mount <mount dir> <json options>
//
func doMount(args []string) exitCode {
	if len(args) < 3 {
		stdoutJson(errorArgs)
		return ExitFailure
	}

	mountDir, rawOpts := args[1], args[2]
	opts, err := parseOptions(rawOpts)
	if err != nil {
		stdoutJson(errorArgs)
		return ExitFailure
	}

	// First: promote drbd resource as primary role
	if err := drbdadm.Primary(opts.ResName); err != nil {
		echo := callEcho{
			Status:  StatusFailure,
			Message: fmt.Sprintf("%s", err),
		}
		stdoutJson(echo)
		return ExitFailure
	}

	device, err := drbdadm.ShDev(opts.ResName)
	if err != nil {
		derr := drbdadm.Secondary(opts.ResName)

		echo := callEcho{
			Status:  StatusFailure,
			Message: fmt.Sprintf("%s, %s", err, derr),
		}
		stdoutJson(echo)
		return ExitFailure
	}

	// Second: format and mount device
	if err := fs.Format(device, opts.FsType); err != nil {
		derr := drbdadm.Secondary(opts.ResName)

		echo := callEcho{
			Status:  StatusFailure,
			Message: fmt.Sprintf("%s, %s", err, derr),
		}
		stdoutJson(echo)
		return ExitFailure
	}
	if err := fs.Mount(device, mountDir); err != nil {
		derr := drbdadm.Secondary(opts.ResName)

		echo := callEcho{
			Status:  StatusFailure,
			Message: fmt.Sprintf("%s, %s", err, derr),
		}
		stdoutJson(echo)
		return ExitFailure
	}

	echo := callEcho{
		Status:  StatusSuccess,
		Message: fmt.Sprintf("resource: %s, device: %s, mount: %s, fstype: %s", opts.ResName, device, mountDir, opts.FsType),
	}
	stdoutJson(echo)

	return ExitSuccess
}

// Unmount:
//
// Unmount the volume. This call-out defaults to bind mount for drivers which
// implement attach & mount-device call-outs. Called only from Kubelet.
//
// <driver executable> unmount <mount dir>
//
func doUnmount(args []string) exitCode {
	if len(args) < 2 {
		stdoutJson(errorArgs)
		return ExitFailure
	}

	mountDir := args[1]

	resName, err := drbdadm.ResByMntDir(mountDir)
	if err != nil {
		echo := callEcho{
			Status:  StatusFailure,
			Message: fmt.Sprintf("%s", err),
		}
		stdoutJson(echo)
		return ExitFailure
	}

	// First: unmount directory
	if err := fs.Unmount(mountDir); err != nil {
		echo := callEcho{
			Status:  StatusFailure,
			Message: fmt.Sprintf("%s", err),
		}
		stdoutJson(echo)
		return ExitFailure
	}

	// Second: demote drbd resource as role secondary
	if err := drbdadm.Secondary(resName); err != nil {
		echo := callEcho{
			Status:  StatusFailure,
			Message: fmt.Sprintf("%s", err),
		}
		stdoutJson(echo)
		return ExitFailure
	}

	echo := callEcho{
		Status:  StatusSuccess,
		Message: fmt.Sprintf("resource: %s, mount: %s", resName, mountDir),
	}
	stdoutJson(echo)

	return ExitSuccess
}

// stdoutJson will do json marshal val, and print the json string to standard
// output. if val marshaled failed, it will print out an empty json object
// string.
func stdoutJson(val interface{}) {
	out, err := json.Marshal(val)
	if err != nil {
		fmt.Print("{}")
	}

	fmt.Print(string(out))
}
