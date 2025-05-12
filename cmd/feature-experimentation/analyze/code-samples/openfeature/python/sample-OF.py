from openfeature import api
from openfeature.provider.in_memory_provider import InMemoryFlag, InMemoryProvider

# flags defined in memory
my_flags = {
  "v2_enabled": InMemoryFlag("on", {"on": True, "off": False})
}

# configure a provider
api.set_provider(InMemoryProvider(my_flags))

# create a client
client = api.get_client()

# get a bool flag value
flag_value = client.get_boolean_value("of_python_v2_enabled", False)
print("Value: " + str(flag_value))