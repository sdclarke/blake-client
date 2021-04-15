while getopts s:b: flag
do
  case "${flag}" in
    s) size=${OPTARG};;
    b) blake=${OPTARG};;
  esac
done
rm output_${size}_${blake}.txt
for f in experiment_logs_${size}/*${blake}*.log
do
  rm ${f}_parsed.txt
  grep -o "Bytes.*$" $f | tee -a output_${size}_${blake}.txt >> ${f}_parsed.txt
  grep -o "Time.*$" $f | tee -a output_${size}_${blake}.txt >> ${f}_parsed.txt
  grep -o "execution_start_timestamp:{seconds:[0-9]*\s*nanos:[0-9]*}\s*execution_completed_timestamp:{seconds:[0-9]*\s*nanos:[0-9]*}" $f | sed -e 's/execution_start_timestamp:{seconds:\([0-9]*\)\s*nanos:\([0-9]*\)}\s*execution_completed_timestamp:{seconds:\([0-9]*\)\s*nanos:\([0-9]*\)}/Execute: \1s\2 \3s\4/' | tee -a output_${size}_${blake}.txt >> ${f}_parsed.txt
  grep -o "output_upload_start_timestamp:{seconds:[0-9]*\s*nanos:[0-9]*}\s*output_upload_completed_timestamp:{seconds:[0-9]*\s*nanos:[0-9]*}" $f | sed -e 's/output_upload_start_timestamp:{seconds:\([0-9]*\)\s*nanos:\([0-9]*\)}\s*output_upload_completed_timestamp:{seconds:\([0-9]*\)\s*nanos:\([0-9]*\)}/Upload: \1s\2 \3s\4/' | tee -a output_${size}_${blake}.txt >> ${f}_parsed.txt
done
