import OpenFeature

Task {
    let provider = CustomProvider()
    // configure a provider, wait for it to complete its initialization tasks
    await OpenFeatureAPI.shared.setProviderAndWait(provider: provider)

    // get a bool flag value
    let client = OpenFeatureAPI.shared.getClient()
    let flagValue = client.getBooleanValue(key: "of_swift_v2_enabled", defaultValue: false)
}