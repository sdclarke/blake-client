#!/bin/bash
blake=true
command=0
while getopts b:s:c: flag
do
  case "${flag}" in
    b) blake=${OPTARG};;
    s) size=${OPTARG};;
    c) command=${OPTARG};;
  esac
done
if [ ! -d "experiment_logs_${size}" ]
then
  mkdir "experiment_logs_${size}"
fi
for e in experiments/*
do
  if [ $size == "huge_single" ] || [ $size == "huge_single_dirty" ] || [ $size == "huge_single_1MiB" ]|| [ $size == "huge_single_1MiB_dirty" ]
  then
    bazel run --override_repository=com_github_buildbarn_bb_storage=/home/debian/bb-storage //cmd/blake3zcc_hasher:blake3zcc_hasher -- -a localhost:8980 -d /home/debian/blake-client/$e -c ${command} -b=${blake} 2> experiment_logs_${size}/output_${e: 12}_${command}_${blake}.log
    continue
  fi
  for c in {0..4}
  do
    if [ $c -lt 2 -o $c -gt 3 ]
    then
      inlist="false"
      for d in "e1" "e4" "e6" "e7"
      do
        if [ $e == "experiments/$d" ]
        then
          inlist="true"
        fi
      done
      if [ $inlist == "false" ]
      then
        continue
      fi
    fi
    bazel run --override_repository=com_github_buildbarn_bb_storage=/home/debian/bb-storage //cmd/blake3zcc_hasher:blake3zcc_hasher -- -a localhost:8980 -d /home/debian/blake-client/$e -c $c -b=${blake} 2> experiment_logs_${size}/output_${e: 12}_${c}_${blake}.log
  done
done
