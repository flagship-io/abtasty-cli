import { OpenFeature } from "@openfeature/server-sdk";

// Register your feature flag provider
await OpenFeature.setProviderAndWait(new YourProviderOfChoice());

// create a new client
const client = OpenFeature.getClient();

// Evaluate your feature flag
const v2Enabled = await client.getBooleanValue("of_ts_v2_enabled", false);

if (v2Enabled) {
  console.log("v2 is enabled");
}
