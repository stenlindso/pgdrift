# pgdrift

Detects schema drift between PostgreSQL databases and outputs structured reports.

## Installation

```bash
go install github.com/yourusername/pgdrift@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/pgdrift.git && cd pgdrift && go build ./...
```

## Usage

Compare two PostgreSQL databases and generate a drift report:

```bash
pgdrift compare \
  --source "postgres://user:pass@localhost:5432/db_production" \
  --target "postgres://user:pass@localhost:5432/db_staging"
```

Output as JSON:

```bash
pgdrift compare \
  --source "postgres://user:pass@source-host/mydb" \
  --target "postgres://user:pass@target-host/mydb" \
  --format json
```

### Example Output

```
[MISSING TABLE]  public.audit_logs        (found in source, missing in target)
[COLUMN DRIFT]   public.users.email       (source: varchar(255), target: text)
[MISSING INDEX]  public.orders.idx_status (found in source, missing in target)

Summary: 3 drift(s) detected across 2 table(s)
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--source` | Source database DSN | required |
| `--target` | Target database DSN | required |
| `--format` | Output format (`text`, `json`) | `text` |
| `--schema` | Schema to inspect | `public` |

## License

MIT