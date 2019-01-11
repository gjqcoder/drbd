//
// Copyright (c) Zhou Peng <p@ctriple.cn>
//
package stor

import (
	"sort"
	"testing"
)

func TestBySizeNr(t *testing.T) {
	hosts := []string{
		"node1.example.com",
		"node2.example.com",
		"node3.example.com",
		"node4.example.com",
	}
	ips := []string{
		"172.25.33.11",
		"172.25.33.12",
		"172.25.33.13",
		"172.25.33.14",
	}
	loadSize := map[string]int{
		"node1.example.com": 5,
		"node2.example.com": 5,
		"node3.example.com": 10,
		"node4.example.com": 10,
	}
	loadNr := map[string]int{
		"node1.example.com": 2,
		"node2.example.com": 1,
		"node3.example.com": 2,
		"node4.example.com": 1,
	}

	t.Log(hosts)
	t.Log(ips)

	sort.Stable(bySizeNr{
		host: hosts,
		ip:   ips,
		size: loadSize,
		nr:   loadNr,
	})

	t.Log(hosts)
	t.Log(ips)
}
