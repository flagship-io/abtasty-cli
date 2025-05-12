public async Task Example()
{
    // Register your feature flag provider
    await Api.Instance.SetProviderAsync(new InMemoryProvider());

    // Create a new client
    FeatureClient client = Api.Instance.GetClient();

    // Evaluate your feature flag
    bool v2Enabled = await client.GetBooleanValueAsync("of_dotnet_v2_enabled", false);

    if ( v2Enabled )
    {
        //Do some work
    }
}

