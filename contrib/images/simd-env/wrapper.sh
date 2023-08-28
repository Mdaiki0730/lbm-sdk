#!/usr/bin/env sh
set -euo pipefail
set -x

BINARY=/rollupd/${BINARY:-rollupd}
ID=${ID:-0}
LOG=${LOG:-rollupd.log}

if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'rollupd'"
	exit 1
fi

export SIMDHOME="/data/node${ID}/rollupd"

if [ -d "$(dirname "${SIMDHOME}"/"${LOG}")" ]; then
  "${BINARY}" --home "${SIMDHOME}" "$@" | tee "${SIMDHOME}/${LOG}"
else
  "${BINARY}" --home "${SIMDHOME}" "$@"
fi
