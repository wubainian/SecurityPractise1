#!/bin/bash

output_dir=output
dst_exe=security_practise1.bin

rm -rf ${output_dir}
mkdir -p ${output_dir}

go build -i -o ${output_dir}/${dst_exe} main.go

