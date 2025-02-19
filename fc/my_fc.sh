#!/bin/sh

cp *.fo targets/
pushd targets && ./fc_all.sh && popd
