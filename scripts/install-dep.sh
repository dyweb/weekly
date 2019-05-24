#!/bin/bash

ROOT=$(dirname ${BASH_SOURCE})/..
ROOT=$(realpath ${ROOT})

WORKDIR=${ROOT}/.jekyll
POSTDIR=${WORKDIR}/_posts

echo 'root = '${ROOT}
echo 'workdir ='${WORKDIR}

cd ${WORKDIR}
bundle install
cd - > /dev/null
