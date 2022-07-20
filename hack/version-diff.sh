#!/bin/bash

# According to:
# https://github.com/cli/cli/blob/235fdcdcc4c3df76adb4a84ecb0ba1ad470fbf5b/pkg/cmd/release/list/http.go#L32
# the GH cli will list releases in descending order so if we list from newest to oldest

read -r vLatest vPrev <<< $(${GITHUB_CLI_BIN} release list -L 2 | awk '{print $1}' | xargs)
DIFF=$(${SEMVER_CLI_BIN} diff $vLatest $vPrev)
PR_TYPE=""
case $DIFF in
  patch)
    PR_TYPE=fix
  ;;
  minor)
    PR_TYPE=feat
  ;;
  major)
    PR_TYPE=!feat
  ;;
  *)
    echo "Bump not necessary"
esac
echo ${PR_TYPE}
