#!/bin/bash

set -e
set -x

python "$(dirname $0)"/updateLicense.py \
	$(go list -json ./... \
	| jq -r '.Dir + "/" + (.GoFiles | .[])' \
	| grep -v /gen/internal/tests/ \
	)
