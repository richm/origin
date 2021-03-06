#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

OS_ROOT=$(dirname "${BASH_SOURCE}")/../..
source "${OS_ROOT}/hack/util.sh"
source "${OS_ROOT}/hack/cmd_util.sh"
os::log::install_errexit

# Cleanup cluster resources created by this test
(
  set +e
  oc delete all,templates --all
  exit 0
) &>/dev/null


# This test validates the debug command
os::cmd::expect_success 'oc create -f test/integration/fixtures/test-deployment-config.yaml'
os::cmd::expect_success_and_text "oc debug dc/test-deployment-config -o yaml" '\- /bin/sh'
os::cmd::expect_success_and_text "oc debug dc/test-deployment-config --keep-annotations -o yaml" 'annotations:'
os::cmd::expect_success_and_text "oc debug dc/test-deployment-config --as-root -o yaml" 'runAsUser: 0'
os::cmd::expect_success_and_text "oc debug dc/test-deployment-config --keep-liveness --keep-readiness -o yaml" ''
os::cmd::expect_success_and_text "oc debug dc/test-deployment-config -o yaml -- /bin/env" '\- /bin/env'
os::cmd::expect_success_and_not_text "oc debug dc/test-deployment-config -o yaml -- /bin/env" 'stdin'
os::cmd::expect_success_and_not_text "oc debug dc/test-deployment-config -o yaml -- /bin/env" 'tty'
# Does not require a real resource on the server
os::cmd::expect_success_and_text "oc debug -f examples/hello-openshift/hello-pod.json --keep-liveness --keep-readiness -o yaml" ''
os::cmd::expect_success_and_text "oc debug -f examples/hello-openshift/hello-pod.json -o yaml -- /bin/env" '\- /bin/env'
os::cmd::expect_success_and_not_text "oc debug -f examples/hello-openshift/hello-pod.json -o yaml -- /bin/env" 'stdin'
os::cmd::expect_success_and_not_text "oc debug -f examples/hello-openshift/hello-pod.json -o yaml -- /bin/env" 'tty'
echo "debug: ok"
