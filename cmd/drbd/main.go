//
// Copyright (c) Zhou Peng <p@ctriple.cn>
//
package main

import (
	"log"
	"os"

	"github.com/ctriple/drbd/pkg/flex"
)

func main() {
	logfile, err := os.OpenFile("/tmp/drbd.log", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer logfile.Close()

	log.SetOutput(logfile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) < 2 {
		log.Fatal(`Usage: drbd action [...]
Action list:
  init
  attach          Not Supported
  detach          Not Supported
  waitforattach   Not Supported
  isattached      Not Supported
  mountdevice     Not Supported
  unmountdeivce   Not Supported
  mount
  umount

More FlexVolume Specification:

https://github.com/kubernetes/community/contributors/devel/flexvolume.md
`)

	}

	action, args := os.Args[1], os.Args[1:]
	exitCode := flex.DriverRun(action, args)

	os.Exit(int(exitCode))
}
