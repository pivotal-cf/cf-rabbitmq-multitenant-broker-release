#!/bin/bash

set -euxo pipefail

ENVIRONMENT="$(jq -r '.name' environment-lock/metadata)"
API_URL="$(jq -r '.cf.api_url' environment-lock/metadata)"
DOMAIN=${API_URL//api./}

bosh interpolate \
	--var deployment-name=cf-rabbitmq-multitenant-broker-release-ci \
	--var rabbitmq-release-name="${RABBITMQ_RELEASE_NAME:=cf-rabbitmq}" \
	--var-errs \
	--ops-file=git-bosh-release/manifests/add-cf-rabbitmq.yml \
	--ops-file=git-bosh-release/manifests/change-vcap-password.yml \
	--ops-file=git-bosh-release/manifests/add-go-syslogd.yml \
	--ops-file=git-bosh-release/manifests/add-syslog.yml \
	--ops-file=git-bosh-release/manifests/add-bosh-dns-aliases.yml \
	--ops-file=cf-rabbitmq-pipelines/manifests/ops-files/add-embedded-tests.yml \
	--vars-file=cf-rabbitmq-pipelines/manifests/vars-files/cf-rabbitmq-vars.yml \
	--vars-file=cf-rabbitmq-pipelines/manifests/vars-files/cf-rabbitmq-multitenant-broker-vars.yml \
	--vars-file=cf-rabbitmq-pipelines/manifests/vars-files/smith-cf-deployment-vars.yml \
	--var cf-admin-password="((/bosh-${ENVIRONMENT}/cf/cf_admin_password))" \
	--var nats-client-cert="((/bosh-${ENVIRONMENT}/cf/nats_client_cert.certificate))" \
	--var nats-client-key="((/bosh-${ENVIRONMENT}/cf/nats_client_cert.private_key))" \
	--var system-domain="$DOMAIN" \
	--var-file stemcell-version=./stemcell-resource/version git-bosh-release/manifests/cf-rabbitmq-broker-template.yml > manifest/manifest.yml
