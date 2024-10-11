## Patches

### 001-remote-and-local-path-options.patch

Added 2 parameters for proxy operation mode:
- `remotepathonly`;
- `localpathalias`;

Example:
```yaml
proxy:
  remoteurl: "..."
  remotepathonly: "sys/deckhouse-oss"
  localpathalias: "system/deckhouse"
  username: "..."
  password: "..."
  ttl: 72h
```
Allows you to specify the allowed path to the registry, as well as replace the path for accessing the caching (local) registry

```bash
# not 'docker pull localhost:5001/sys/deckhouse-oss/install:latest'
docker pull localhost:5001/system/deckhouse/install:latest
```