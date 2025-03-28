public void example() {

    // flags defined in memory
    Map<String, Flag<?>> myFlags = new HashMap<>();
    myFlags.put("v2_enabled", Flag.builder()
            .variant("on", true)
            .variant("off", false)
            .defaultVariant("on")
            .build());

    // configure a provider
    OpenFeatureAPI api = OpenFeatureAPI.getInstance();
    api.setProviderAndWait(new InMemoryProvider(myFlags));

    // create a client
    Client client = api.getClient();

    // get a bool flag value
    boolean flagValue = client.getBooleanValue("of_java_v2_enabled", false);
}