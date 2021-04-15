#!/bin/bash
clean=""
while getopts s:c flag
do
  case "${flag}" in
    s) size=${OPTARG};;
    c) clean=TRUE;;
  esac
done
rm -rf experiments
if [ $clean = "TRUE" ]
then
  exit 0
fi
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
    docker save ubuntu:base > in0
    cd ../e2
    docker save ubuntu:wget > in0
    cd ../e3
    docker save ubuntu:git > in0
    cd ../e4
    id=$(docker create ubuntu:base)
    docker export $id > in0
    docker container rm $id
    cd ../e5
    id=$(docker create ubuntu:wget)
    docker export $id > in0
    docker container rm $id
    cd ../e6
    id=$(docker create ubuntu:git)
    docker export $id > in0
    docker container rm $id
    cd ../e7
    docker save alpine:base > in0
    cd ../e8
    docker save alpine:wget > in0
    cd ../e9
    docker save alpine:git > in0
    cd ../e10
    id=$(docker create alpine:base)
    docker export $id > in0
    docker container rm $id
    cd ../e11
    id=$(docker create alpine:wget)
    docker export $id > in0
    docker container rm $id
    cd ../e12
    id=$(docker create alpine:git)
    docker export $id > in0
    docker container rm $id
    mkdir ../e13
    cd ../e13
    docker save debian:base > in0
    mkdir ../e14
    cd ../e14
    docker save debian:wget > in0
    mkdir ../e15
    cd ../e15
    docker save debian:git > in0
    mkdir ../e16
    cd ../e16
    id=$(docker create debian:base)
    docker export $id > in0
    docker container rm $id
    mkdir ../e17
    cd ../e17
    id=$(docker create debian:wget)
    docker export $id > in0
    docker container rm $id
    mkdir ../e18
    cd ../e18
    id=$(docker create debian:git)
    docker export $id > in0
    docker container rm $id
    ;;
qcow2)
    cp ~/arch-boxes/output/Arch-Linux-x86_64-basic-20210415.0.qcow2 ./in0
    cd ../e2
    cp ~/arch-boxes/output/Arch-Linux-x86_64-cloudimg-20210415.0.qcow2 ./in0
    ;;
esac
