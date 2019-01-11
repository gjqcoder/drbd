# Getting Started

This file describe design and implementation details, you need some backgrounds
to get deep into this project.

## Kubernetes

- External Storage, FlexVolume
- Pod, Job, Deployment
- StorageClass, PersistentVolume, PersistentVolumeClaim

[kubernetes](https://kubernetes.io)

## DRBD

DRBD (distributed replication block device) is working at the linux kernel
level. it means that its work is transparent and high performance, but also
means that lack of userspace awareness such as filesystem.

[drbd](https://docs.linbit.com/docs/users-guide-9.0/)

## Building

This project was written in Go programming, so you should have a working Go
development environment to build.

[Go](https://golang.org)

## Deploying

You should have a working kubernetes cluster, with some prerequisite already
meet. See [README](README.md).

# Design

## drbd

FlexVolume driver, it is called by kubelet to do PersistentVolume house keeping
like mount/umount disk to Pod.

`Run out of kubernetes cluster as a host process.`

## stor

Kubernetes external storage provisioner, it anwsers PersistentVolumeClaim
requests and provision PersistentVolume.

`Run as a Deployment inside cluster.`

## sync

Working as a temporary job to do node specific work such disk allocation and
free.

`Run as a temporary Job inside cluster.`

# Implementation

## lvm

We use lvm utility to manage disk dynamic provision on all kubernetes nodes, you
need to get all your kubernetes nodes with volume group **centos** created in
advance. Because lvm utility needs to access host /proc and /sys, but our sync
job runs inside container, we use ssh from container to host to workaround this
limitation.

## drbdadm

As the same to lvm, we have built many wrapper functions to use drbdadm from Go
program, such as.

- primary/secondary
- up/down
- sh-dev/sh-ll-dev
- sh-resources

## sync job

The sync programm implements all node operations like create/delete drbd
resource, the work parameters are passed in by container environments at each
launch of a job.

## image ctriple/drbd:latest

For easy deploy and management, we package all executables into one docker image
ctriple/drbd:latest, and container start command is specified explicitly at
launch time.
