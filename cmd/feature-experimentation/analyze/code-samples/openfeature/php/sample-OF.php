use OpenFeature\OpenFeatureAPI;
use OpenFeature\Providers\Flagd\FlagdProvider;

function example()
{
$api = OpenFeatureAPI::getInstance();

// configure a provider
$api->setProvider(new FlagdProvider());

// create a client
$client = $api->getClient();

// get a bool flag value
$client->getBooleanValue('of_php_v2_enabled', false);
}