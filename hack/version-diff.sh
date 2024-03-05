#!/bin/bash

# According to:
# https://github.com/cli/cli/blob/235fdcdcc4c3df76adb4a84ecb0ba1ad470fbf5b/pkg/cmd/release/list/http.go#L32
# the GH cli will list releases in descending order so if we list from newest to oldest
while [ "$1" != "" ]; do
  case $1 in
  --version-latest)
    shift
    version_latest=$1
    ;;
  --version-previous)
    shift
    version_previous=$1
    ;;
    esac
  shift
done

if [[ -z $version_latest  || -z $version_previous ]]; then
  read -r version_latest version_previous <<< "$(${GITHUB_CLI_BIN} release list -L 2 | awk '{print $1}' | xargs)"
fi

DIFF="$(${SEMVER_CLI_BIN} diff "${version_latest}" "${version_previous}")"
PR_TYPE=""
case $DIFF in
  patch)
    PR_TYPE=fix
  ;;
  minor)
    PR_TYPE=feat
  ;;
  major)
    PR_TYPE=feat!
  ;;
  *)
    echo "Bump not necessary"
esac
echo ${PR_TYPE}
