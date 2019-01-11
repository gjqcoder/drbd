#
# ctriple/drbd:latest
#
FROM alpine:3.8
MAINTAINER Zhou Peng <p@ctriple.cn>

ADD drbd /drbd
ADD sync /sync
ADD stor /stor
