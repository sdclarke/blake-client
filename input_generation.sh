#!/bin/bash

while getopts s: flag
do
  case "${flag}" in
    s) size=${OPTARG};;
  esac
done
case $size in 
  small)
    rm -rf experiments
    mkdir -p experiments/e{1..12}
    cd experiments/e1
    dd if=/dev/zero of=in0 bs=512 count=2
    dd if=/dev/zero of=in1 bs=512 count=3
    cd ../e2
    dd if=/dev/zero of=in0 bs=512 count=2
    dd if=/dev/zero of=in1 bs=512 count=15
    cd ../e3
    dd if=/dev/zero of=in0 bs=512 count=2
    dd if=/dev/zero of=in1 bs=512 count=16
    cd ../e4
    dd if=/dev/zero of=in0 bs=512 count=3
    dd if=/dev/zero of=in1 bs=512 count=15
    cd ../e5
    dd if=/dev/zero of=in0 bs=512 count=3
    dd if=/dev/zero of=in1 bs=512 count=16
    cd ../e6
    dd if=/dev/zero of=in0 bs=512 count=15
    dd if=/dev/zero of=in1 bs=512 count=16
    cd ../e7
    dd if=/dev/zero of=in0 bs=512 count=16
    dd if=/dev/zero of=in1 bs=512 count=15
    cd ../e8
    dd if=/dev/zero of=in0 bs=512 count=16
    dd if=/dev/zero of=in1 bs=512 count=3
    cd ../e9
    dd if=/dev/zero of=in0 bs=512 count=16
    dd if=/dev/zero of=in1 bs=512 count=2
    cd ../e10
    dd if=/dev/zero of=in0 bs=512 count=15
    dd if=/dev/zero of=in1 bs=512 count=3
    cd ../e11
    dd if=/dev/zero of=in0 bs=512 count=15
    dd if=/dev/zero of=in1 bs=512 count=2
    cd ../e12
    dd if=/dev/zero of=in0 bs=512 count=3
    dd if=/dev/zero of=in1 bs=512 count=2
    ;;
esac
