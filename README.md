# Cloud Foundry RabbitMQ Multi-tenant Broker

This repository contains the release for a multi-tenant RabbitMQ service broker
for Cloud Foundry. It's deployable by BOSH in the usual way.

## Dependencies

- [bundler](http://bundler.io/)

## Install

Clone the repository and run `./scripts/update-release` to update submodules and install dependencies.

## Deploying

To deploy the release into BOSH you will need a deployment manifest. You can generate a deployment manifest using the following command:
```sh
boshgo interpolate \
  --vars-file=manifests/lite-vars-file.yml \
  --var=director-uuid=$(bosh status --uuid) \
  manifests/cf-rabbitmq-broker-template.yml > manifests/cf-rabbitmq-broker.yml
```

Once you have a [BOSH Lite up and running locally](https://github.com/cloudfoundry/bosh-lite), run `scripts/deploy-to-bosh-lite`.

## Testing

To run all the tests do `bundle exec rake spec`.

Use `rspec` to run a specific test:
`bundle exec rspec spec/integration/broker_registrar_spec.rb`

### Unit Tests

To run only unit tests locally, run: `bundle exec rake spec:unit`. Unit tests do not require the release to be deployed.

### Integration Tests

Integration tests require this release to be deployed into a BOSH director (see [Deploying section above](#deploying)).

To run integration tests, run: `bundle exec rake spec:integration`.
