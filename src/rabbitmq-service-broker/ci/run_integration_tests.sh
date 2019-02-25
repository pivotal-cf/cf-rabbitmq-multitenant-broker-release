#!/usr/bin/env bash

# Runner for the integration tests in a container
# - on a development mac, run 'make integration_tests'
# - in Concourse, run this script

set -eu

cd "$(dirname "${BASH_SOURCE[0]}")/.."

echo "About to start rabbitmq"
service rabbitmq-server start
rabbitmq-plugins enable rabbitmq_management

echo "About to run service broker integration tests"
GOFLAGS='-mod=vendor' ginkgo -v -r integrationtests
