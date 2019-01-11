//
// Copyright (c) Zhou Peng <p@ctriple.cn>
//
package main

import (
	"os"
	"path"
	"strings"

	"github.com/ctriple/drbd/pkg/defs"
	"github.com/ctriple/drbd/pkg/drbdadm"
	"github.com/ctriple/drbd/pkg/sync"
	"github.com/ctriple/drbd/pkg/sync/lvm"
	"github.com/ctriple/drbd/pkg/sync/res"
	"github.com/golang/glog"
)

func init() {
	if err := sync.InstallDriver(); err != nil {
		glog.Fatal(err)
	}
}

func main() {
	var (
		job     = os.Getenv(defs.SyncJob_EnvJob)
		resName = os.Getenv(defs.SyncJob_EnvResName)
		resSize = os.Getenv(defs.SyncJob_EnvResSize)
		host    = os.Getenv(defs.SyncJob_EnvResHost)
		ip      = os.Getenv(defs.SyncJob_EnvResIP)

		hosts = strings.Split(host, ",")
		ips   = strings.Split(ip, ",")
	)

	switch defs.SyncJob(job) {
	case defs.SyncJob_New:
		if drbdadm.ShResource(resName) {
			glog.Fatalln(resName, "already exist!")
		}
		if err := doNew(resName, resSize, hosts, ips); err != nil {
			glog.Fatalln(err)
		}

	case defs.SyncJob_Del:
		if !drbdadm.ShResource(resName) {
			glog.Fatalln(resName, "does not exist!")
		}
		if err := doDel(resName); err != nil {
			glog.Fatalln(err)
		}

	default:
		glog.Fatalln("env:", defs.SyncJob_EnvJob, "must be", defs.SyncJob_New, "or", defs.SyncJob_Del)
	}
}

func doNew(resName, resSize string, hosts, ips []string) error {
	if err := lvm.Create(defs.DrbdDiskVG, resName, resSize); err != nil {
		return err
	}
	// lvm allocated disk pattern: /dev/{vg}/{name}
	disk := path.Join("/dev", defs.DrbdDiskVG, resName)
	if err := res.New(resName, disk, hosts, ips); err != nil {
		lvm.Remove(disk)
		return err
	}
	if err := drbdadm.CreateMD(resName); err != nil {
		lvm.Remove(disk)
		res.Del(resName)
		return err
	}
	if err := drbdadm.Up(resName); err != nil {
		return err
	}

	return nil
}

func doDel(resName string) error {
	disk, err := drbdadm.ShLlDev(resName)
	if err != nil {
		return err
	}
	if err := drbdadm.Down(resName); err != nil {
		return err
	}
	if err := lvm.Remove(disk); err != nil {
		drbdadm.Up(resName)
		return err
	}
	if err := res.Del(resName); err != nil {
		return err
	}

	return nil
}
