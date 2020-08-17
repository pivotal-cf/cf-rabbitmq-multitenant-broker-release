#!/bin/bash

set -euxo pipefail

bosh interpolate --var deployment-name=cf-rabbitmq-multitenant-broker-release-ci \
	--var-errs --ops-file=git-bosh-release/manifests/add-cf-rabbitmq.yml \
	--ops-file=git-bosh-release/manifests/change-vcap-password.yml \
	--ops-file=git-bosh-release/manifests/add-go-syslogd.yml \
	--ops-file=git-bosh-release/manifests/add-syslog.yml \
	--ops-file=git-cf-rabbitmq-pipelines/manifests/ops-files/add-embedded-tests.yml \
	--vars-file=git-cf-rabbitmq-pipelines/manifests/vars-files/cf-rabbitmq-vars.yml \
	--vars-file=git-cf-rabbitmq-pipelines/manifests/vars-files/cf-rabbitmq-multitenant-broker-vars.yml \
	--vars-file=git-cf-rabbitmq-pipelines/manifests/vars-files/smith-cf-deployment-vars.yml \
	--var-file stemcell-version=./stemcell-resource/version git-bosh-release/manifests/cf-rabbitmq-broker-template.yml > manifest/manifest.yml
