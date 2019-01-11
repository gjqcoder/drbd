//
// Copyright (c) Zhou Peng <p@ctriple.cn>
//
package defs

const (
	Vendor = "ctriple.cn"
	Driver = "drbd"
)

const (
	// ctriple.cn drbd driver identity
	DrbdDriver = Vendor + "/" + Driver

	// Choose drbd resource port within this range (firewall accepted)
	DrbdPortMin = 7000
	DrbdPortMax = 7100

	DrbdReplicaMin = 2
	DrbdReplicaMax = 4

	// Lvm volume group from which drbd backing disk alloc
	DrbdDiskVG = "centos"
)

type SyncJob string

const (
	SyncJob_New = "SYNCJOB_NEW"
	SyncJob_Del = "SYNCJOB_DEL"
)

const (
	SyncJob_EnvJob     = "SYNCJOB_JOB"
	SyncJob_EnvResName = "SYNCJOB_RESOURCE_NAME"
	SyncJob_EnvResSize = "SYNCJOB_RESOURCE_SIZE"
	SyncJob_EnvResHost = "SYNCJOB_RESOURCE_HOST"
	SyncJob_EnvResIP   = "SYNCJOB_RESOURCE_IP"
)
