#!/bin/sh

export ENV_STR2INT_MAP_DER_NAME=./sample.d/str2int.asn1.der.dat

input=./sample.d/input.maps.asn1.der.dat
output=./sample.d/output.maps.asn1.der.dat

genmap(){
	echo creating str-to-int map...
	python3 ./sample.d/genmap.py |
		dd \
			if=/dev/stdin \
			of="${ENV_STR2INT_MAP_DER_NAME}" \
			status=none
}

originalMaps(){
	python3 ./sample.d/originalMaps.py |
		dd \
			if=/dev/stdin \
			of="${input}" \
			status=none
}

test -f "${ENV_STR2INT_MAP_DER_NAME}" || genmap

test -f "${input}" || originalMaps

echo converting the original map to normalized map...
cat "${input}" |
	./asn1kvpairs-key2int |
	dd \
		if=/dev/stdin \
		of="${output}" \
		bs=1048576 \
		status=none

echo printing the original map...
cat "${input}" |
	python3 ./sample.d/original2print.py |
	jq -c

echo printing the normalized map...
cat "${output}" |
	python3 ./sample.d/mapd2print.py |
	jq -c

ls -l "${input}" "${output}"
