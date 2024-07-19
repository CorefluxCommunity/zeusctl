# Zeus CLI

- Simple CLI to facilitate Vault operations, like auth and fetching secrets using context configs

## Usage

### Auth

```
zeusctl vault auth --host HOST:PORT --ca-path CA_CERT_PATH --user USER --password PASSWORD
```

### Load and Export Context

```
eval $(zeusctl vault context CONTEXT_NAME)
```

