#!/bin/bash
blake=true
command=0
decompose=$((1024*1024))
while getopts b:s:c:d: flag
do
  case "${flag}" in
    b) blake=${OPTARG};;
    s) size=${OPTARG};;
    c) command=${OPTARG};;
    d) decompose=${OPTARG};;
  esac
done
if [ ! -d "experiment_logs_${size}" ]
then
  mkdir "experiment_logs_${size}"
fi
if [ $size == "docker_image" ]
then
    for e in experiments/e{1..3} experiments/e{7..9} experiments/e{13..15}
    do
	    bazel run --override_repository=com_github_buildbarn_bb_storage=/home/debian/bb-storage //cmd/blake3zcc_hasher:blake3zcc_hasher -- -a localhost:8980 -d /home/debian/blake-client/$e -b=${blake} -D -s ${decompose} 2> experiment_logs_${size}/output_${e: 12}_${blake}_$((${decompose}/1024)).log
    done
elif [ $size == "docker_fs" ]
then
    for e in experiments/e{4..6} experiments/e{10..12} experiments/e{16..18}
    do
	    bazel run --override_repository=com_github_buildbarn_bb_storage=/home/debian/bb-storage //cmd/blake3zcc_hasher:blake3zcc_hasher -- -a localhost:8980 -d /home/debian/blake-client/$e -b=${blake} -D -s ${decompose} 2> experiment_logs_${size}/output_${e: 12}_${blake}_$((${decompose}/1024)).log
    done
else
    for e in experiments/e{1..2}
    do
	    bazel run --override_repository=com_github_buildbarn_bb_storage=/home/debian/bb-storage //cmd/blake3zcc_hasher:blake3zcc_hasher -- -a localhost:8980 -d /home/debian/blake-client/$e -b=${blake} -D -s ${decompose} 2> experiment_logs_${size}/output_${e: 12}_${blake}_$((${decompose}/1024)).log
    done
fi
