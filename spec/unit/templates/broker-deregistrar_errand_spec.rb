RSpec.describe 'broker deregistration errand template', template: true do
  let(:output) {
    compiled_template('broker-deregistrar', 'errand.sh', manifest_properties)
  }

  describe 'authentication methods' do
    context 'when user credentials are configured' do
      let(:manifest_properties) do
        {
          'cf' => {
            'admin_username' => 'some-username',
            'admin_password' => 'some-password'
          }
        }
      end

      it 'authenticates using admin username and password' do
        expect(output).to include 'CF_AUTH_COMMAND="$CF_ADMIN_USERNAME $CF_ADMIN_PASSWORD"'
      end
    end

    context 'when oauth client is configured' do
      let(:manifest_properties) do
        {
          'cf' => {
            'admin_client' => 'some-client',
            'admin_client_secret' => 'some-secret'
          }
        }
      end

      it 'authenticates with cf api using admin client and secret' do
        expect(output).to include 'CF_AUTH_COMMAND="--client-credentials $CF_ADMIN_CLIENT $CF_ADMIN_CLIENT_SECRET"'
      end
    end

    context 'when both oauth client and admin user are configured' do
      let(:manifest_properties) do
        {
          'cf' => {
            'admin_username' => 'some-username',
            'admin_password' => 'some-password',
            'admin_client' => 'some-client',
            'admin_client_secret' => 'some-secret'
          }
        }
      end

      it 'authenticates with cf api using admin client and secret' do
        expect(output).to include 'CF_AUTH_COMMAND="--client-credentials $CF_ADMIN_CLIENT $CF_ADMIN_CLIENT_SECRET"'
      end
    end
  end

  context 'when ssl validation is skipped' do
    let(:manifest_properties) { { 'cf' => { 'skip_ssl_validation' => true } } }

    it 'skips ssl validation' do
      expect(output).to include 'cf api --skip-ssl-validation'
    end
  end

  context 'when ssl validation is respected' do
    let(:manifest_properties) { { 'cf' => { 'skip_ssl_validation' => false } } }

    it 'do not skip ssl validation' do
      expect(output).not_to include '--skip-ssl-validation'
      expect(output).to include 'cf api'
    end
  end

  context 'when I set my broker service name' do
		let(:manifest_properties) { { 'broker' => { 'service' => { 'name' => 'broker_service_name' } } } }

    it 'purges the correct service offering' do
      expect(output).to include "cf purge-service-offering -f 'broker_service_name'"
    end
  end

  context 'when deleting the service broker' do
    let(:manifest_properties) { { 'broker' => { 'name' => 'changed_broker_name' } } }

    it 'deletes the correct broker' do
      expect(output).to include "cf delete-service-broker -f 'changed_broker_name'"
    end
  end
end

