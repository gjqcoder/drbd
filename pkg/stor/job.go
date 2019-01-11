//
// Copyright (c) Zhou Peng <p@ctriple.cn>
//
package stor

import (
	"os"

	"github.com/golang/glog"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	SyncJobNamespace      string
	SyncJobPodImage       string
	SyncJobServiceAccount string
)

const (
	SyncJobGenerateName  = "sync-"
	SyncJobContainerName = "sync"
)

func init() {
	SyncJobNamespace = os.Getenv("MY_POD_NAMESPACE")
	if SyncJobNamespace == "" {
		glog.Fatalln("env MY_POD_NAMESPACE not set!")
	}

	SyncJobPodImage = os.Getenv("MY_POD_IMAGE")
	if SyncJobPodImage == "" {
		glog.Fatalln("env MY_POD_IMAGE not set!")
	}

	SyncJobServiceAccount = os.Getenv("MY_POD_SERVICEACCOUNT")
	if SyncJobServiceAccount == "" {
		glog.Fatalln("env MY_POD_SERVICEACCOUNT not set!")
	}
}

// syncJob returns a batch job instance, this job runs ctriple/stor:latest image
// [/sync] command on the specified kubernetes nodes. This job has 2 works to
// do:
//
// 1. install ctriple.cn/drbd flexvolume to the host nodes
// 2. initialize drbd nodes
//
// Job's control parameters are passed in as environment variables by the job
// creater.
//
// Because ctriple/stor:latest image [/sync] command needs to do some privileged
// job, it needs to mount many host pathes and be privileged.
func syncJob() *batchv1.Job {
	privileged := true

	c := v1.Container{
		Name:            SyncJobContainerName,
		Image:           SyncJobPodImage,
		Command:         []string{"/sync"},
		SecurityContext: &v1.SecurityContext{Privileged: &privileged},
		VolumeMounts: []v1.VolumeMount{
			v1.VolumeMount{Name: "host-bin", MountPath: "/bin", ReadOnly: true},
			v1.VolumeMount{Name: "host-sbin", MountPath: "/sbin", ReadOnly: true},
			v1.VolumeMount{Name: "host-usr-bin", MountPath: "/usr/bin", ReadOnly: true},
			v1.VolumeMount{Name: "host-root", MountPath: "/root", ReadOnly: true},
			v1.VolumeMount{Name: "host-lib", MountPath: "/lib", ReadOnly: true},
			v1.VolumeMount{Name: "host-lib64", MountPath: "/lib64", ReadOnly: true},
			v1.VolumeMount{Name: "host-dev", MountPath: "/dev", ReadOnly: false},
			v1.VolumeMount{Name: "host-etc", MountPath: "/etc", ReadOnly: false},
			v1.VolumeMount{Name: "host-flex-driver-dir", MountPath: "/flexmnt", ReadOnly: false},
		},
	}

	vols := []v1.Volume{
		v1.Volume{
			Name: "host-bin",
			VolumeSource: v1.VolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: "/bin",
				},
			},
		},
		v1.Volume{
			Name: "host-sbin",
			VolumeSource: v1.VolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: "/sbin",
				},
			},
		},
		v1.Volume{
			Name: "host-usr-bin",
			VolumeSource: v1.VolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: "/usr/bin",
				},
			},
		},
		v1.Volume{
			Name: "host-root",
			VolumeSource: v1.VolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: "/root",
				},
			},
		},
		v1.Volume{
			Name: "host-lib",
			VolumeSource: v1.VolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: "/lib",
				},
			},
		},
		v1.Volume{
			Name: "host-lib64",
			VolumeSource: v1.VolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: "/lib64",
				},
			},
		},
		v1.Volume{
			Name: "host-dev",
			VolumeSource: v1.VolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: "/dev",
				},
			},
		},
		v1.Volume{
			Name: "host-etc",
			VolumeSource: v1.VolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: "/etc",
				},
			},
		},
		v1.Volume{
			Name: "host-flex-driver-dir",
			VolumeSource: v1.VolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: "/usr/libexec/kubernetes/kubelet-plugins/volume/exec",
				},
			},
		},
	}

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: SyncJobGenerateName,
			Namespace:    SyncJobNamespace,
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: SyncJobGenerateName,
				},
				Spec: v1.PodSpec{
					HostPID:            true,
					HostIPC:            true,
					HostNetwork:        true,
					ServiceAccountName: SyncJobServiceAccount,
					Containers:         []v1.Container{c},
					Volumes:            vols,
					RestartPolicy:      v1.RestartPolicyNever,
				},
			},
		},
	}

	return job
}
