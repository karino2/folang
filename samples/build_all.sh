sed 's/ .*$//' filelist.txt | xargs -L1 ./myfc.sh
./build_sample_md filelist.txt
