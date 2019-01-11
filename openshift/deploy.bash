#!/usr/bin/env bash

#
# Copyright (c) Zhou Peng <p@ctriple.cn>
#

# This scripts will deploy drbd powered kubernetes dynamic storage solution into
# your cluster. It will first create namespaces and serviceaccount as needed,
# and grant priorities to the serviceaccount. then create external provisioner
# and storageclass.

set -o errexit
set -o nounset
set -o pipefail

oc create -f 1-ns.yaml
oc create -f 2-sa.yaml

oc adm policy add-scc-to-user          privileged    -z drbd -n ctriple-drbd
oc adm policy add-cluster-role-to-user cluster-admin -z drbd -n ctriple-drbd

oc create -f 3-dc.yaml
oc create -f 4-sc.yaml
