#!/bin/bash

ROOT=$(dirname "${BASH_SOURCE}")/..

YEAR=`date +%Y`

cd $ROOT

for mdfile in $(find `seq 2015 $YEAR` -name '*.md') # or whatever other pattern...
do
	cat ./scripts/header.txt $mdfile > $mdfile.new
	mv $mdfile.new $mdfile
done

cd - > /dev/null
