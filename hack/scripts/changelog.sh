#!/usr/bin/env sh
# Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
# Use of this source code is covered by the following dual licenses:
# (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
# (2) Affero General Public License 3.0 (AGPL 3.0)
# license that can be found in the LICENSE file.

set -eu

echo '# Changelog'
echo

tag=
git tag -l 'v*' | sort -rV | while read last; do
  if [ "$tag" != "" ]; then
    echo "## $(git for-each-ref --format='%(refname:strip=2) (%(creatordate:short))' refs/tags/${tag})"
    echo
    git_log='git --no-pager log --no-merges --invert-grep --grep=^\(build\|ci\|docs\|test\|chore\):'
	  $git_log --format=' * [%h](https://github.com/easysoft/quickon_cli/commit/%H) %s' $last..$tag
	  echo
	  echo "### Contributors"
	  echo
	  $git_log --format=' * %an'  $last..$tag | sort -u
	  echo
  fi
  tag=$last
done
