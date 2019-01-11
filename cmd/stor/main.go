//
// Copyright (c) Zhou Peng <p@ctriple.cn>
//
package main

import (
	"github.com/ctriple/drbd/pkg/defs"
	"github.com/ctriple/drbd/pkg/stor"
	"github.com/golang/glog"
	"github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/controller"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		glog.Fatalf("Failed to create config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Fatalf("Failed to create client: %v", err)
	}

	serverVersion, err := clientset.Discovery().ServerVersion()
	if err != nil {
		glog.Fatalf("Error getting server version: %v", err)
	}

	flexProvisioner := stor.NewFlexProvisioner(clientset)

	pc := controller.NewProvisionController(clientset, defs.DrbdDriver, flexProvisioner, serverVersion.GitVersion)

	pc.Run(wait.NeverStop)
}
