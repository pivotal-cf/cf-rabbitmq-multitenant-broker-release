#!/bin/bash

set -euxo pipefail

ENVIRONMENT_LOCK_FILE="${ENVIRONMENT_LOCK_FILE:='environment-lock/metadata'}"
ENVIRONMENT="$(jq -r '.name' "$ENVIRONMENT_LOCK_FILE")"
API_URL="$(jq -r '.cf.api_url' "$ENVIRONMENT_LOCK_FILE")"
DOMAIN=${API_URL//api./}

set +x
  eval "$(smith -l "$ENVIRONMENT_LOCK_FILE" bosh)"
  eval "$(smith -l "$ENVIRONMENT_LOCK_FILE" om)"
  cf_guid="$(om -k curl -s -p /api/v0/staged/products | jq -r '.[] | select(.type == "cf") | .guid')"
set -x

bosh interpolate \
	--var deployment-name=cf-rabbitmq-multitenant-broker-release-ci \
	--var rabbitmq-release-name="${RABBITMQ_RELEASE_NAME:=cf-rabbitmq}" \
	--var-errs \
	--ops-file=git-bosh-release/manifests/add-cf-rabbitmq.yml \
	--ops-file=git-bosh-release/manifests/change-vcap-password.yml \
	--ops-file=git-bosh-release/manifests/add-go-syslogd.yml \
	--ops-file=git-bosh-release/manifests/add-syslog.yml \
	--ops-file=cf-rabbitmq-pipelines/manifests/ops-files/add-embedded-tests.yml \
	--vars-file=cf-rabbitmq-pipelines/manifests/vars-files/cf-rabbitmq-vars.yml \
	--vars-file=cf-rabbitmq-pipelines/manifests/vars-files/cf-rabbitmq-multitenant-broker-vars.yml \
	--vars-file=cf-rabbitmq-pipelines/manifests/vars-files/smith-cf-deployment-vars.yml \
	--var cf-admin-password="$(om -k curl -s -p "/api/v0/deployed/products/$cf_guid/credentials/.uaa.admin_credentials" | jq -r .credential.value.password)" \
	--var nats-client-cert="((/opsmgr/${BOSH_DEPLOYMENT}/nats_client_cert.cert_pem))" \
	--var nats-client-key="((/opsmgr/${BOSH_DEPLOYMENT}/nats_client_cert.private_key_pem))" \
	--var system-domain="$DOMAIN" \
	--var cf-deployment-name="${BOSH_DEPLOYMENT}" \
	--var-file stemcell-version=./stemcell-resource/version git-bosh-release/manifests/cf-rabbitmq-broker-template.yml > manifest/manifest.yml
