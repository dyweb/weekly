#!/bin/bash

ROOT=$(dirname ${BASH_SOURCE})/..
ROOT=$(realpath ${ROOT})

WORKDIR=${ROOT}/.jekyll
POSTDIR=${WORKDIR}/_posts

echo 'root = '${ROOT}
echo 'workdir ='${WORKDIR}

cd ${WORKDIR}
# bundle install
mkdir -p ${POSTDIR}
seq 2015 `date +%Y` | xargs -I% mv ${ROOT}/% ${POSTDIR}

bundle exec jekyll build --source ${WORKDIR} --destination ${ROOT}/docs
ls ${POSTDIR} | xargs -I% mv ${POSTDIR}/% ${ROOT}
cd - > /dev/null
