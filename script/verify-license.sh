#!/usr/bin/env bash

error_lines=()

license_file="$(pwd)/LICENSE"

if [ ! -f ${license_file} ]; then
  error_lines+=( "error: missing ${license_file}" )
else
  tmp_license=$(mktemp)

  wget -q https://joinup.ec.europa.eu/sites/default/files/custom-page/attachment/2020-03/EUPL-1.2%20EN.txt -O ${tmp_license}

  if ! /usr/bin/diff --brief ${license_file} ${tmp_license} >/dev/null; then
    error_lines+=( "error: license content is not correct" )
  fi

  rm ${tmp_license}
fi

while IFS= read -r error_line; do
  error_lines+=( "${error_line}" )
done < <( find main.go src/ hookscript-tests/ acceptance-tests/ -type f -exec grep -zvq "Copyright (C) $(date +'%Y') Rea Sand" {} \; -exec printf "error: bad/missing copyright notice in file '%s'\n" {} \; )

while IFS= read -r error_line; do
  error_lines+=( "${error_line}" )
done < <( find main.go src/ hookscript-tests/ acceptance-tests/ -type f -exec grep -zvq "Licensed under the EUPL" {} \; -exec printf "error: bad/missing license notice in file '%s'\n" {} \; )

if [ ${#error_lines[@]} -gt 0 ]; then
  printf '%s\n' "${error_lines[@]}"
  exit 1
else
  printf "OK\n"
  exit 0
fi
