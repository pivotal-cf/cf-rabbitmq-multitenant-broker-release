# Cloud Foundry RabbitMQ Multi-tenant Broker

This repository contains the release for a multi-tenant RabbitMQ service broker
for Cloud Foundry. It's deployable by BOSH in the usual way.

## Dependencies

In order to test this release locally you will need:

- [bundler](http://bundler.io/)
- [BOSH CLI v2](https://bosh.io/docs/cli-v2.html#install)
- [BOSH Lite](https://bosh.io/docs/bosh-lite)

## Install

Clone the repository and run `./scripts/update-release` to update submodules and install dependencies.

## Deploying

To deploy the release into BOSH you will need a deployment manifest. You can generate a deployment manifest using the following command:
```sh
alias boshgo=bosh # This is just to make pcf-rabbitmq tile team's life simpler
boshgo interpolate \
  --vars-file=manifests/lite-vars-file.yml \
  --var=director-uuid=$(bosh status --uuid) \
  manifests/cf-rabbitmq-broker-template.yml > manifests/cf-rabbitmq-broker.yml
```

Once you have a [BOSH Lite up and running locally](https://github.com/cloudfoundry/bosh-lite), run `./scripts/deploy-to-bosh-lite`.

## Testing

To run all the tests do `bundle exec rake spec`.

Use `rspec` to run a specific test:
`bundle exec rspec spec/integration/broker_registrar_spec.rb`

### Unit Tests

To run only unit tests locally, run: `./scripts/unit-tests`. Unit tests do not require the release to be deployed.

## Troubleshooting

### An error occurred while installing capybara-webkit (macOS)
```bash
An error occurred while installing capybara-webkit (1.11.1), and Bundler cannot continue.
Make sure that `gem install capybara-webkit -v '1.11.1'` succeeds before bundling.
```

Some of the tests in this repository use `prof`, which depends on `capybara`.
The error occurs when Xcode is not installed, and `capybara` needs Xcode to get installed. More details [here](https://github.com/thoughtbot/capybara-webkit/issues/813)

#### To solve the problem:
- Go to the App Store and install Xcode
- run `sudo xcode-select --switch /Applications/Xcode.app/Contents/Developer`
- run `sudo xcodebuild -license`
