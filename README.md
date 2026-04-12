# envoy-cli

> A CLI tool for managing and syncing `.env` files across local and remote environments with secret redaction support.

---

## Installation

```bash
go install github.com/yourusername/envoy-cli@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/envoy-cli/releases).

---

## Usage

```bash
# Push your local .env to a remote environment
envoy push --env .env --remote staging

# Pull remote environment variables to a local file
envoy pull --remote production --out .env.production

# Sync environments with secrets automatically redacted in logs
envoy sync --from staging --to production --redact

# List all tracked environments
envoy list
```

### Flags

| Flag | Description |
|------|-------------|
| `--env` | Path to the local `.env` file |
| `--remote` | Target remote environment name |
| `--redact` | Mask secret values in output logs |
| `--out` | Output file path for pulled variables |

---

## Configuration

`envoy-cli` looks for a `envoy.yaml` config file in the current directory or `~/.config/envoy/envoy.yaml`.

```yaml
remote:
  url: https://your-env-server.example.com
  token: $ENVOY_TOKEN
environments:
  - staging
  - production
```

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

---

## License

[MIT](LICENSE)