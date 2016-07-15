#!/bin/sh -xe

dnf install rpm-build -y

curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| GH=mh-cbon/go-bin-rpm sh -xe

curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| GH=mh-cbon/changelog sh -xe

cd /docker
TAG=$1
NAME=$2

if [[ -z ${TAG} ]]; then TAG="0.0.0"; fi

VERBOSE=* go-bin-rpm generate -a 386 --version ${TAG} -b pkg-build/386/ -o ${NAME}-386.rpm
VERBOSE=* go-bin-rpm generate -a amd64 --version ${TAG} -b pkg-build/amd64/ -o ${NAME}-amd64.rpm
