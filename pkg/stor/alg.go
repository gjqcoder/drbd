//
// Copyright (c) Zhou Peng <p@ctriple.cn>
//
package stor

import (
	"sort"

	"github.com/ctriple/drbd/pkg/defs"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// candidates returns all available nodes, best-fit nodes will be at first
// position.
//
// Candidate nodes order are sorted by two factors:
//   - node's total allocated resource size
//   - node's resource number already on
func (p *flexProvisioner) candidates() (hosts, ips []string, err error) {
	loadNr, loadSize, err := p.load()
	if err != nil {
		return
	}

	hosts, ips, err = p.nodes()
	if err != nil {
		return
	}

	sort.Stable(bySizeNr{
		host: hosts,
		ip:   ips,
		size: loadSize,
		nr:   loadNr,
	})

	return
}

func (p *flexProvisioner) nodes() (hosts, ips []string, err error) {
	nodes, err := p.client.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return
	}

	for _, node := range nodes.Items {
		var host string
		var ip string

		for _, addr := range node.Status.Addresses {
			switch addr.Type {
			case v1.NodeInternalIP:
				ip = addr.Address
			case v1.NodeHostName:
				host = addr.Address
			}
		}
		hosts = append(hosts, host)
		ips = append(ips, ip)
	}

	return
}

func (p *flexProvisioner) load() (loadNr, loadSize map[string]int, err error) {
	pvs, err := p.client.CoreV1().PersistentVolumes().List(metav1.ListOptions{})
	if err != nil {
		return
	}

	loadNr = make(map[string]int)
	loadSize = make(map[string]int)

	for _, pv := range pvs.Items {
		if provisioner := pv.Annotations[pvCreatedBy]; provisioner != defs.DrbdDriver {
			continue
		}

		hosts := pv.Spec.NodeAffinity.Required.NodeSelectorTerms[0].MatchExpressions[0].Values
		capacity := pv.Spec.Capacity[v1.ResourceStorage]
		size := int(capacity.Value())
		for _, h := range hosts {
			loadNr[h] = loadNr[h] + 1
			loadSize[h] = loadSize[h] + size
		}
	}

	return
}

type bySizeNr struct {
	host []string
	ip   []string
	size map[string]int
	nr   map[string]int
}

func (by bySizeNr) Len() int {
	return len(by.host)
}

func (by bySizeNr) Less(i, j int) bool {
	hi, hj := by.host[i], by.host[j]
	switch {

	// First: size
	case by.size[hi] < by.size[hj]:
		return true
	case by.size[hi] > by.size[hj]:
		return false

	// Second: number
	case by.nr[hi] < by.nr[hj]:
		return true
	case by.nr[hi] > by.nr[hj]:
		return false

	// Third: name
	default:
		return hi < hj
	}
}

func (by bySizeNr) Swap(i, j int) {
	by.host[i], by.host[j] = by.host[j], by.host[i]
	by.ip[i], by.ip[j] = by.ip[j], by.ip[i]
}
