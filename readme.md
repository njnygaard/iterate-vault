# vault-iterator

## Testing

```bash
cd vault-iterator
cat << 'EOF' > .auth.yaml
token: [Vault Token]
vault_addr: https://vault.build.splicemachine-dev.io/
EOF
go test
```

## CLI Usage

```bash
cd cli
cat << 'EOF' > .auth.yaml
token: [Vault Token]
vault_addr: https://vault.build.splicemachine-dev.io/
EOF
go build
./cli g secret/deployments/k8s/default/services/cloudmanager
```
