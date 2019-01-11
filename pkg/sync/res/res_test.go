//
// Copyright (c) Zhou Peng <p@ctriple.cn>
//
package res

import (
	"io/ioutil"
	"path"
	"testing"
)

var (
	resName = "res-testing"
	disk    = "/dev/lvm/res-testing"

	hosts = []string{
		"node1.example.com",
		"node2.example.com",
		"node3.example.com",
	}

	ips = []string{
		"172.25.33.11",
		"172.25.33.12",
		"172.25.33.13",
	}
)

func init() {
	resOutDir = "testdata"
}

func TestNew(t *testing.T) {
	if err := New(resName, disk, hosts, ips); err != nil {
		t.Fatal(err)
	}
}

// This pseudo testcase is used to display the generated resource file before it
// is deleted.
func Test_show_res(t *testing.T) {
	resFile := path.Join(resOutDir, resName+".res")
	data, err := ioutil.ReadFile(resFile)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}

func TestDel(t *testing.T) {
	if err := Del(resName); err != nil {
		t.Fatal(err)
	}
}
