# Kubernetes Dynamic Storage Solution Powered By DRBD

This project implements a simple storage solution for kubernetes. it is based on
DRBD (distribute replication block device) open source software, which is
natively data redundancy and high performance.

## Prerequisite

We asumed that you have all your kubernetes nodes installed lvm-utils and at
least DRBD 9.0 (userspace utils and kernel module). For more about DRBD, please
refer to [LINBIT](https://linbit.com/docs/users-guide-9.0).

## DRBD Network RAID

Userspace Apps filesystem operations get replicated to another disk on the
network node, the whole process is transparent and high performance. Apps can be
migrated to that secondary node quickly and seamlessly with data already have
been placed right there.

```text
      DRBD Primary                          DRBD Secondary
+-----------------------+               +----------------------+
| Userspace             |               | Userspace            |
|                       |               |                      |
| open/read/write/close |               |                      |
+--------------|--------+               +----------------------+
| Kernel       |        |    tcp/ip     | Kernel               |
|              +----------->---->---->--------------+          |
|              |        |               |           |          |
+--------------v--------+               +-----------v----------+
| Disk: /dev/drbd0      |               | Disk: /dev/drbd0     |
+-----------------------+               +----------------------+
```

## Kubernetes meets DRBD

Kubernetes is good at stateless Pod migration, but this is not true as to
satefull Pod. So here DRBD comes to play.

- Pod mysql runs on k8s-node-1

```text
        k8s-node-1                                 k8s-node-2
+--------------------------+               +---------------------------+
| Kubelet: drbd flexvolume |               | Kubelete: drbd flexvolume |
|--------------------------|               |---------------------------|
| Pod: mysql               |               |                           |
|        |                 |               |                           |
+--------v-----------------+    tcp/ip     +---------------------------+
| Disk: /dev/drbd0 -->-->--|--->--->--->---| Disk: /dev/drbd0          |
+--------------------------+               +---------------------------+
        DRBD Primary                               DRBD Secondary
```

- Pod mysql got killed and migrated to k8s-node-2

```text
        k8s-node-1                                 k8s-node-2
+--------------------------+               +---------------------------+
| Kubelet: drbd flexvolume |               | Kubelete: drbd flexvolume |
|--------------------------|               |---------------------------|
|                          |               | Pod: mysql                |
|                          |               |        |                  |
+--------------------------+    tcp/ip     +--------v------------------+
| Disk: /dev/drbd0 --<--<--|---<---<---<---| Disk: /dev/drbd0          |
+--------------------------+               +---------------------------+
        DRBD Secondary                             DRBD Primary
```

## Limitations

### Pros

- Simple
- No **extra** storage cluster
- Data natively redundancy and high performance
- Easy to use (StorageClass, PersistentVolume, PersistentVolumeClaim)

### Cons

- DRBD>=9.0 can be up to 16 nodes, but only at most one primary node simultaneously
- Pod must run as single replica instance (ReadWritOnce)
- PV is available only on resource placed nodes
