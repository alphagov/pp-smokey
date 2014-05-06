#!/bin/bash -e

# Download the JSON results file of pp-smokey from Jenkins. Exit with a code
# which can be interpreted as a Sensu check: http://sensuapp.org/docs/0.12/checks
#
# Example: check_feature.sh smoke_test_name
#
# To run locally, do this first:
# ```
# export JENKINS_URL=https://deploy.preview.performance.service.gov.uk
# ```
# Exit codes:
# 0 = OK
# 1 = WARNING
# 2 = CRITICAL
# 3 = UNKNOWN



JENKINS_URL=${JENKINS_URL-http://jenkins:8080}

THIS_SCRIPT=$(readlink -f "$0")
THIS_DIR=$(dirname "${THIS_SCRIPT}")
CUCUMBER_DIR="${THIS_DIR}/.."

TEST_FEATURE_NAME="$1"

cd ${CUCUMBER_DIR}

function parse_arguments {
    if [ "${TEST_FEATURE_NAME}" = "" ]; then
        echo "Usage: $0 <feature_name>"
        exit 1
    fi
}

function download_cucumber_results_file {
    RESULTS_JSON=$(mktemp --suffix=json)
    curl  -o ${RESULTS_JSON} ${JENKINS_URL}/job/pp-smokey/lastBuild/artifact/results.json
}

function parse_results_json_file {
    set +e
    sensu/parse_result_file.py ${TEST_FEATURE_NAME} ${RESULTS_JSON}
    EXIT_CODE=$?
    set -e
}

parse_arguments
download_cucumber_results_file
parse_results_json_file

echo "Exiting with code: ${EXIT_CODE}"
exit ${EXIT_CODE}
