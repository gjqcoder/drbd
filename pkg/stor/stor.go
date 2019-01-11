//
// Copyright (c) Zhou Peng <p@ctriple.cn>
//
package stor

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ctriple/drbd/pkg/defs"
	"github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/controller"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/kubelet/apis"
)

const (
	pvCreatedBy = "kubernetes.io/createdby"
)

var _ controller.Provisioner = &flexProvisioner{}

type flexProvisioner struct {
	client   kubernetes.Interface
	identity types.UID
}

func NewFlexProvisioner(client kubernetes.Interface) *flexProvisioner {
	var identity types.UID

	flexProvisioner := &flexProvisioner{
		client:   client,
		identity: identity,
	}

	return flexProvisioner
}

func (p *flexProvisioner) Provision(options controller.VolumeOptions) (*v1.PersistentVolume, error) {
	capacity := options.PVC.Spec.Resources.Requests[v1.ResourceStorage]
	requestedBytes := capacity.Value()

	replicas := defs.DrbdReplicaMin
	fstype := "ext4"

	for k, v := range options.Parameters {
		switch strings.ToLower(k) {
		case "replicas":
			r64, err := strconv.ParseInt(v, 10, 64)
			if err == nil {
				r := int(r64)
				switch {
				case r < defs.DrbdReplicaMin:
					replicas = defs.DrbdReplicaMin
				case r > defs.DrbdReplicaMax:
					replicas = defs.DrbdReplicaMax
				default:
					replicas = r
				}
			}
		case "fstype":
			fstype = v
		}
	}

	resName := fmt.Sprintf("%s-%s", options.PVC.ObjectMeta.Namespace, options.PVC.ObjectMeta.Name)
	resSize := fmt.Sprintf("%dM", (requestedBytes/1024/1024 + 1))

	// -- Use our host choosen algorithm
	hosts, ips, err := p.candidates()
	if err != nil {
		return nil, err
	}
	if len(hosts) < replicas {
		return nil, fmt.Errorf("candidates:%v less than exptected replicas:%d", hosts, replicas)
	}
	hosts = hosts[:replicas]
	ips = ips[:replicas]

	// -- Run sync job on each choosen host

	jobClient := p.client.BatchV1().Jobs(SyncJobNamespace)
	syncJob := syncJob()

	jobEnvs := []v1.EnvVar{
		{Name: defs.SyncJob_EnvJob, Value: defs.SyncJob_New},
		{Name: defs.SyncJob_EnvResName, Value: resName},
		{Name: defs.SyncJob_EnvResSize, Value: resSize},
		{Name: defs.SyncJob_EnvResHost, Value: strings.Join(hosts, ",")},
		{Name: defs.SyncJob_EnvResIP, Value: strings.Join(ips, ",")},
	}
	syncJob.Spec.Template.Spec.Containers[0].Env = jobEnvs

	complete := []string{}
	failed := []string{}

	for _, h := range hosts {
		// Run job on this host
		syncJob.Spec.Template.Spec.NodeSelector = map[string]string{apis.LabelHostname: h}
		newJob, err := jobClient.Create(syncJob)
		if err != nil {
			failed = append(failed, h)
			continue
		}

	ThisJobExit:
		for {
			getJob, err := jobClient.Get(newJob.Name, metav1.GetOptions{IncludeUninitialized: true})
			if err != nil {
				failed = append(failed, h)
				break
			}
			for _, c := range getJob.Status.Conditions {
				switch c.Type {
				case batchv1.JobComplete:
					complete = append(complete, h)
					break ThisJobExit
				case batchv1.JobFailed:
					failed = append(failed, h)
					break ThisJobExit
				}
			}
		}
	}

	// -- Partially completion, should clean up already completed host
	if len(complete) < len(hosts) {
		jobEnvs := []v1.EnvVar{
			{Name: defs.SyncJob_EnvJob, Value: defs.SyncJob_Del},
			{Name: defs.SyncJob_EnvResName, Value: resName},
			{Name: defs.SyncJob_EnvResSize, Value: "not-used"},
			{Name: defs.SyncJob_EnvResHost, Value: "not-used"},
			{Name: defs.SyncJob_EnvResIP, Value: "not-used"},
		}
		syncJob.Spec.Template.Spec.Containers[0].Env = jobEnvs

		for _, h := range complete {
			syncJob.Spec.Template.Spec.NodeSelector = map[string]string{apis.LabelHostname: h}
			jobClient.Create(syncJob)
		}

		return nil, fmt.Errorf("Sync job complete:%v failed:%v", complete, failed)
	}

	// -- All sync job completed successfully, pv provision ok. Note that
	// this pv is available only on the drbd nodes.

	pv := &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:   resName,
			Labels: map[string]string{},
			Annotations: map[string]string{
				pvCreatedBy: defs.DrbdDriver,
			},
		},
		Spec: v1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: options.PersistentVolumeReclaimPolicy,
			AccessModes:                   options.PVC.Spec.AccessModes,
			Capacity: v1.ResourceList{
				v1.ResourceStorage: options.PVC.Spec.Resources.Requests[v1.ResourceStorage],
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{
				FlexVolume: &v1.FlexPersistentVolumeSource{
					Driver: defs.DrbdDriver,
					FSType: fstype,
					Options: map[string]string{
						"resource": resName,
					},
				},
			},
			NodeAffinity: &v1.VolumeNodeAffinity{
				Required: &v1.NodeSelector{
					NodeSelectorTerms: []v1.NodeSelectorTerm{
						{
							MatchExpressions: []v1.NodeSelectorRequirement{
								{
									Key:      apis.LabelHostname,
									Operator: v1.NodeSelectorOpIn,
									Values:   hosts,
								},
							},
						},
					},
				},
			},
		},
	}

	return pv, nil
}

func (p *flexProvisioner) Delete(volume *v1.PersistentVolume) error {
	// -- pv not provisioned by ctriple.cn/drbd
	provisioner := volume.Annotations[pvCreatedBy]
	if provisioner != defs.DrbdDriver {
		return fmt.Errorf("pv not provisioned by: %v", defs.DrbdDriver)
	}

	resName := volume.Name

	jobClient := p.client.BatchV1().Jobs(SyncJobNamespace)
	syncJob := syncJob()
	jobEnvs := []v1.EnvVar{
		{Name: defs.SyncJob_EnvJob, Value: defs.SyncJob_Del},
		{Name: defs.SyncJob_EnvResName, Value: resName},
		{Name: defs.SyncJob_EnvResSize, Value: "not-used"},
		{Name: defs.SyncJob_EnvResHost, Value: "not-used"},
		{Name: defs.SyncJob_EnvResIP, Value: "not-used"},
	}
	syncJob.Spec.Template.Spec.Containers[0].Env = jobEnvs
	hosts := volume.Spec.NodeAffinity.Required.NodeSelectorTerms[0].MatchExpressions[0].Values

	complete := []string{}
	failed := []string{}

	for _, h := range hosts {
		// Run job on this host
		syncJob.Spec.Template.Spec.NodeSelector = map[string]string{apis.LabelHostname: h}
		newJob, err := jobClient.Create(syncJob)
		if err != nil {
			failed = append(failed, h)
			continue
		}

	ThisJobExit:
		for {
			getJob, err := jobClient.Get(newJob.Name, metav1.GetOptions{IncludeUninitialized: true})
			if err != nil {
				failed = append(failed, h)
				break
			}
			for _, c := range getJob.Status.Conditions {
				switch c.Type {
				case batchv1.JobComplete:
					complete = append(complete, h)
					break ThisJobExit
				case batchv1.JobFailed:
					failed = append(failed, h)
					break ThisJobExit
				}
			}
		}
	}

	// FIXME: If ctriple.cn/drbd provisioned pv was deleted partially
	// completed, i have no idea how to recover from this situation, you
	// should never use this pv anymore and try to clean up by hand.
	if len(complete) < len(hosts) {
		return fmt.Errorf("Sync job complete:%v failed:%v", complete, failed)
	}

	return nil
}
