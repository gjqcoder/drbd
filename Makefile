#
# Copyright (c) Zhou Peng <p@ctriple.cn>
#

all: test build image clean

test:
	env MY_POD_NAMESPACE=ctriple-drbd MY_POD_IMAGE=ctriple/drbd:latest MY_POD_SERVICEACCOUNT=drbd go test github.com/ctriple/drbd/...

build:
	go build github.com/ctriple/drbd/cmd/drbd
	go build github.com/ctriple/drbd/cmd/stor
	go build github.com/ctriple/drbd/cmd/sync

image:
	docker build -t ctriple/drbd:latest .
	docker push     ctriple/drbd:latest

clean:
	@git clean -xdf .
