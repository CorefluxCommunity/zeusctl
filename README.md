# Coreflux CLI

- Simple CLI with Coreflux Cloud basic utilities, like auth and fetching Vault secrets using context configs

## Usage

### Auth

```
cfctl vault auth --host HOST:PORT --ca-path CA_CERT_PATH --user USER --password PASSWORD
```

### Load and Export Context

```
eval $(cfctl vault context CONTEXT_NAME)
```

