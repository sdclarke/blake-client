#!/bin/bash

while getopts s: flag
do
  case "${flag}" in
    s) size=${OPTARG};;
  esac
done
rm -rf experiments
mkdir -p experiments/e{1..12}
cd experiments/e1
case $size in 
  small)
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
  small_multi)
    dd if=/dev/zero of=in0 bs=$((1*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((2*8192+1)) count=1
    cd ../e2
    dd if=/dev/zero of=in0 bs=$((1*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((4*8192+1)) count=1
    cd ../e3
    dd if=/dev/zero of=in0 bs=$((1*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((5*8192-1)) count=1
    cd ../e4
    dd if=/dev/zero of=in0 bs=$((2*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((4*8192+1)) count=1
    cd ../e5
    dd if=/dev/zero of=in0 bs=$((2*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((5*8192-1)) count=1
    cd ../e6
    dd if=/dev/zero of=in0 bs=$((4*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((5*8192-1)) count=1
    cd ../e7
    dd if=/dev/zero of=in0 bs=$((5*8192-1)) count=1
    dd if=/dev/zero of=in1 bs=$((4*8192+1)) count=1
    cd ../e8
    dd if=/dev/zero of=in0 bs=$((5*8192-1)) count=1
    dd if=/dev/zero of=in1 bs=$((2*8192+1)) count=1
    cd ../e9
    dd if=/dev/zero of=in0 bs=$((5*8192-1)) count=1
    dd if=/dev/zero of=in1 bs=$((1*8192+1)) count=1
    cd ../e10
    dd if=/dev/zero of=in0 bs=$((4*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((2*8192+1)) count=1
    cd ../e11
    dd if=/dev/zero of=in0 bs=$((4*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((1*8192+1)) count=1
    cd ../e12
    dd if=/dev/zero of=in0 bs=$((2*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((1*8192+1)) count=1
    ;;
  medium)
    dd if=/dev/zero of=in0 bs=$((5*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((20*8192+1)) count=1
    cd ../e2
    dd if=/dev/zero of=in0 bs=$((5*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((35*8192+1)) count=1
    cd ../e3
    dd if=/dev/zero of=in0 bs=$((5*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((55*8192+1)) count=1
    cd ../e4
    dd if=/dev/zero of=in0 bs=$((20*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((35*8192+1)) count=1
    cd ../e5
    dd if=/dev/zero of=in0 bs=$((20*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((55*8192+1)) count=1
    cd ../e6
    dd if=/dev/zero of=in0 bs=$((35*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((55*8192+1)) count=1
    cd ../e7
    dd if=/dev/zero of=in0 bs=$((55*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((35*8192+1)) count=1
    cd ../e8
    dd if=/dev/zero of=in0 bs=$((55*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((20*8192+1)) count=1
    cd ../e9
    dd if=/dev/zero of=in0 bs=$((55*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((5*8192+1)) count=1
    cd ../e10
    dd if=/dev/zero of=in0 bs=$((35*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((20*8192+1)) count=1
    cd ../e11
    dd if=/dev/zero of=in0 bs=$((35*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((5*8192+1)) count=1
    cd ../e12
    dd if=/dev/zero of=in0 bs=$((20*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((5*8192+1)) count=1
    ;;
  medium_high)
    dd if=/dev/zero of=in0 bs=$((65*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((107*8192+1)) count=1
    cd ../e2
    dd if=/dev/zero of=in0 bs=$((65*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((149*8192+1)) count=1
    cd ../e3
    dd if=/dev/zero of=in0 bs=$((65*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((191*8192+1)) count=1
    cd ../e4
    dd if=/dev/zero of=in0 bs=$((107*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((149*8192+1)) count=1
    cd ../e5
    dd if=/dev/zero of=in0 bs=$((107*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((191*8192+1)) count=1
    cd ../e6
    dd if=/dev/zero of=in0 bs=$((149*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((191*8192+1)) count=1
    cd ../e7
    dd if=/dev/zero of=in0 bs=$((191*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((149*8192+1)) count=1
    cd ../e8
    dd if=/dev/zero of=in0 bs=$((191*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((107*8192+1)) count=1
    cd ../e9
    dd if=/dev/zero of=in0 bs=$((191*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((65*8192+1)) count=1
    cd ../e10
    dd if=/dev/zero of=in0 bs=$((149*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((107*8192+1)) count=1
    cd ../e11
    dd if=/dev/zero of=in0 bs=$((149*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((65*8192+1)) count=1
    cd ../e12
    dd if=/dev/zero of=in0 bs=$((107*8192+1)) count=1
    dd if=/dev/zero of=in1 bs=$((65*8192+1)) count=1
    ;;
  large)
    dd if=/dev/zero of=in0 bs=$((2*1024*1024+1)) count=1
    dd if=/dev/zero of=in1 bs=$((18*1024*1024+1)) count=1
    cd ../e2
    dd if=/dev/zero of=in0 bs=$((2*1024*1024+1)) count=1
    dd if=/dev/zero of=in1 bs=$((34*1024*1024+1)) count=1
    cd ../e3
    dd if=/dev/zero of=in0 bs=$((2*1024*1024+1)) count=1
    dd if=/dev/zero of=in1 bs=$((50*1024*1024)) count=1
    cd ../e4
    dd if=/dev/zero of=in0 bs=$((18*1024*1024+1)) count=1
    dd if=/dev/zero of=in1 bs=$((34*1024*1024))  count=1
    cd ../e5
    dd if=/dev/zero of=in0 bs=$((18*1024*1024+1)) count=1
    dd if=/dev/zero of=in1 bs=$((50*1024*1024))  count=1
    cd ../e6
    dd if=/dev/zero of=in0 bs=$((34*1024*1024+1)) count=1
    dd if=/dev/zero of=in1 bs=$((50*1024*1024))  count=1
    cd ../e7
    dd if=/dev/zero of=in0 bs=$((50*1024*1024+1)) count=1
    dd if=/dev/zero of=in1 bs=$((34*1024*1024))  count=1
    cd ../e8
    dd if=/dev/zero of=in0 bs=$((50*1024*1024+1)) count=1
    dd if=/dev/zero of=in1 bs=$((18*1024*1024))  count=1
    cd ../e9
    dd if=/dev/zero of=in0 bs=$((50*1024*1024+1)) count=1
    dd if=/dev/zero of=in1 bs=$((2*1024*1024)) count=1
    cd ../e10
    dd if=/dev/zero of=in0 bs=$((34*1024*1024+1)) count=1
    dd if=/dev/zero of=in1 bs=$((18*1024*1024))  count=1
    cd ../e11
    dd if=/dev/zero of=in0 bs=$((34*1024*1024+1)) count=1
    dd if=/dev/zero of=in1 bs=$((2*1024*1024)) count=1
    cd ../e12
    dd if=/dev/zero of=in0 bs=$((18*1024*1024+1)) count=1
    dd if=/dev/zero of=in1 bs=$((2*1024*1024)) count=1
    ;;
  huge)
    dd if=/dev/zero of=in0 bs=$((265*1024*1024+1)) count=1
    dd if=/dev/zero of=in1 bs=$((275*1024*1024+1)) count=1
    cd ../e2
    dd if=/dev/zero of=in0 bs=$((265*1024*1024+1)) count=1
    dd if=/dev/zero of=in1 bs=$((530*1024*1024+1)) count=1
    cd ../e3
    dd if=/dev/zero of=in0 bs=$((265*1024*1024+1)) count=1
    dd if=/dev/zero of=in1 bs=$((540*1024*1024)) count=1
    cd ../e4
    dd if=/dev/zero of=in0 bs=$((275*1024*1024+1)) count=1
    dd if=/dev/zero of=in1 bs=$((530*1024*1024))  count=1
    cd ../e5
    dd if=/dev/zero of=in0 bs=$((275*1024*1024+1)) count=1
    dd if=/dev/zero of=in1 bs=$((540*1024*1024))  count=1
    cd ../e6
    dd if=/dev/zero of=in0 bs=$((530*1024*1024+1)) count=1
    dd if=/dev/zero of=in1 bs=$((540*1024*1024))  count=1
    cd ../e7
    dd if=/dev/zero of=in0 bs=$((540*1024*1024+1)) count=1
    dd if=/dev/zero of=in1 bs=$((530*1024*1024))  count=1
    cd ../e8
    dd if=/dev/zero of=in0 bs=$((540*1024*1024+1)) count=1
    dd if=/dev/zero of=in1 bs=$((275*1024*1024))  count=1
    cd ../e9
    dd if=/dev/zero of=in0 bs=$((540*1024*1024+1)) count=1
    dd if=/dev/zero of=in1 bs=$((265*1024*1024)) count=1
    cd ../e10
    dd if=/dev/zero of=in0 bs=$((530*1024*1024+1)) count=1
    dd if=/dev/zero of=in1 bs=$((275*1024*1024))  count=1
    cd ../e11
    dd if=/dev/zero of=in0 bs=$((530*1024*1024+1)) count=1
    dd if=/dev/zero of=in1 bs=$((265*1024*1024)) count=1
    cd ../e12
    dd if=/dev/zero of=in0 bs=$((275*1024*1024+1)) count=1
    dd if=/dev/zero of=in1 bs=$((265*1024*1024)) count=1
    ;;
  huge_single)
    dd if=/dev/zero of=in0 bs=$((1024*1024)) count=200
    cd ../e2
    dd if=/dev/zero of=in0 bs=$((1024*1024)) count=250
    cd ../e3
    dd if=/dev/zero of=in0 bs=$((1024*1024)) count=300
    cd ../e4
    dd if=/dev/zero of=in0 bs=$((1024*1024)) count=350
    cd ../e5
    dd if=/dev/zero of=in0 bs=$((1024*1024)) count=400
    cd ../e6
    dd if=/dev/zero of=in0 bs=$((1024*1024)) count=450
    cd ../e7
    dd if=/dev/zero of=in0 bs=$((1024*1024)) count=500
    cd ../e8
    dd if=/dev/zero of=in0 bs=$((1024*1024)) count=550
    cd ../e9
    dd if=/dev/zero of=in0 bs=$((1024*1024)) count=600
    cd ../e10
    dd if=/dev/zero of=in0 bs=$((1024*1024)) count=650
    cd ../e11
    dd if=/dev/zero of=in0 bs=$((1024*1024)) count=700
    cd ../e12
    dd if=/dev/zero of=in0 bs=$((1024*1024)) count=750
    ;;
  docker)
    docker save ubuntu:0 > in0
    cd ../e2
    docker save ubuntu:1 > in0
    ;;
esac
