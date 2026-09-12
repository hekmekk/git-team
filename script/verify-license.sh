#!/usr/bin/env bash

error_lines=()

# TODO: verify content
# full content match or just some heuristic like below
# wget -q https://www.gnu.org/licenses/gpl-3.0.txt -O -
if [ ! -f COPYING ]; then
  error_lines+=( "error: missing COPYING file" )
fi

if [ "$(readlink LICENSE)" != "COPYING" ]; then
  error_lines+=( "error: missing LICENSE symlink" )
fi

copyright_notice_pattern="Copyright (C) $(date +'%Y') Rea Sand"
while IFS= read -r error_line; do
  error_lines+=( "${error_line}" )
done < <( find main.go src/ hookscript-tests/ acceptance-tests/ -type f -exec grep -zvq "${copyright_notice_pattern}" {} \; -exec printf "error: bad/missing copyright notice in file '%s'\n" {} \; )

license_notice_pattern="This file is part of.*This program is free software.*GNU General Public License.*WITHOUT ANY WARRANTY.*If not, see <https://www.gnu.org/licenses/>"
while IFS= read -r error_line; do
  error_lines+=( "${error_line}" )
done < <( find main.go src/ hookscript-tests/ acceptance-tests/ -type f -exec grep -zvq "${license_notice_pattern}" {} \; -exec printf "error: bad/missing license notice in file '%s'\n" {} \; )

if [ ${#error_lines[@]} -gt 0 ]; then
  printf '%s\n' "${error_lines[@]}"
  exit 1
else
  printf "OK\n"
  exit 0
fi
