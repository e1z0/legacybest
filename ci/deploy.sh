#!/bin/bash

VERSION=$(cat VERSION)
BUILD=`date "+%F-%T"`
COMMIT=$(git rev-parse HEAD)
BASENEIMAS=legacybest-${VERSION}-${COMMIT}
SSH_HOST="$DEPLOY_HOST"
SSH_USER="$DEPLOY_LOGIN"
SSH_PORT="$DEPLOY_PORT"
SSH_PASS="$DEPLOY_PASS"
SSH_DIR="$DEPLOY_PATH"

echo "Deploying releases..."
# linux
sshpass -p ${SSH_PASS} ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ${SSH_USER}@${SSH_HOST} -p ${SSH_PORT} "mkdir -p ${SSH_DIR}/linux/${BUILD}"
sshpass -p ${SSH_PASS} scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -P ${SSH_PORT} builds/${BASENEIMAS}-linux-x64.zip ${SSH_USER}@${SSH_HOST}:${SSH_DIR}/linux/${BUILD}
sshpass -p ${SSH_PASS} ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ${SSH_USER}@${SSH_HOST} -p ${SSH_PORT} "ln -sf linux/${BUILD}/${BASENEIMAS}-linux-x64.zip ${SSH_DIR}/latest-linux-x64.zip"
# windows
sshpass -p ${SSH_PASS} ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ${SSH_USER}@${SSH_HOST} -p ${SSH_PORT} "mkdir -p ${SSH_DIR}/windows/${BUILD}"
sshpass -p ${SSH_PASS} scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -P ${SSH_PORT} builds/${BASENEIMAS}-win-x64.zip ${SSH_USER}@${SSH_HOST}:${SSH_DIR}/windows/${BUILD}
sshpass -p ${SSH_PASS} scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -P ${SSH_PORT} builds/${BASENEIMAS}-win-x86.zip ${SSH_USER}@${SSH_HOST}:${SSH_DIR}/windows/${BUILD}
sshpass -p ${SSH_PASS} ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ${SSH_USER}@${SSH_HOST} -p ${SSH_PORT} "ln -sf windows/${BUILD}/${BASENEIMAS}-win-x64.zip ${SSH_DIR}/latest-windows-x64.zip"
sshpass -p ${SSH_PASS} ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ${SSH_USER}@${SSH_HOST} -p ${SSH_PORT} "ln -sf windows/${BUILD}/${BASENEIMAS}-win-x86.zip ${SSH_DIR}/latest-windows-x86.zip"
# mac
sshpass -p ${SSH_PASS} ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ${SSH_USER}@${SSH_HOST} -p ${SSH_PORT} "mkdir -p ${SSH_DIR}/mac/${BUILD}"
sshpass -p ${SSH_PASS} scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -P ${SSH_PORT} builds/${BASENEIMAS}-mac-x64.zip ${SSH_USER}@${SSH_HOST}:${SSH_DIR}/mac/${BUILD}
sshpass -p ${SSH_PASS} ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ${SSH_USER}@${SSH_HOST} -p ${SSH_PORT} "ln -sf mac/${BUILD}/${BASENEIMAS}-mac-x64.zip ${SSH_DIR}/latest-mac-x64.zip"

