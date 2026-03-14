# garden

A CLI tool for tracking your seed inventory, planning when to plant things, and figuring out frost dates for your location.

Data is stored in a SQLite database at `~/.garden/garden.db`.

## Install

```
go build -o garden.exe .
```

Or just use `go run .` during development.

## First-time setup

Tell the app where you live so it can calculate frost dates:

```
garden locale set --zip 80203
# or by state:
garden locale set --state CO
```

## Commands

### seeds — your seed stash

```
garden seeds list
garden seeds add --name Tomato --variety "Cherokee Purple" --qty 2 --unit packets
garden seeds remove 3
garden seeds link <seed-id> <spec-id>   # connect a seed to a plant spec
```

### plants — the plant library

A built-in reference of plant specs with growing info (days to maturity, sun, spacing, frost timing, etc.).

```
garden plants list
garden plants list --sun full
garden plants search tomato
garden plants show 12
```

### schedule — plan what to plant and when

```
garden schedule list
garden schedule add --plant Tomato --type indoor_start --date 2026-03-01
garden schedule suggest --plant Tomato        # calculates dates from your frost dates
garden schedule done 4                        # mark entry #4 as planted
garden schedule remove 4
```

Planting types: `indoor_start`, `transplant`, `direct_sow`

The `suggest` command looks up your frost dates and tells you the optimal window to start seeds indoors or direct sow outside.

### locale — your location

```
garden locale show
garden locale set --zip 80203
garden locale set --state CO
```

### serve — web UI

```
garden serve
garden serve --port 9090
```

Opens a browser UI at `http://localhost:8080`.

## Database location

Default: `~/.garden/garden.db`

Override with the `--db` flag on any command:

```
garden --db ./mygarden.db seeds list
```

---

## Deploying to Portainer — Hey Tony 👋

There are two ways to deploy. **Option A is easier** if the server has internet access.

### Option A — Deploy straight from GitHub (recommended)

1. Open Portainer → **Stacks** → **Add stack**
2. Give it a name, e.g. `garden-app`
3. Select **Git repository**
4. Set the repository URL:
   ```
   https://github.com/alexramsey92/garden-app
   ```
5. Compose path:
   ```
   docker-compose.yml
   ```
6. Leave everything else as default and click **Deploy the stack**

Portainer will pull the code, build the image, and start the container.

To update the app later: go to the stack in Portainer and click **Pull and redeploy**.

---

### Option B — Pre-built image from a registry

If you'd rather not build on the server, build the image on your own machine first:

```bash
docker build -t alexramsey92/garden-app:latest .
docker push alexramsey92/garden-app:latest
```

Then edit `docker-compose.yml` — replace `build: .` with:

```yaml
image: alexramsey92/garden-app:latest
```

Deploy the same way as Option A.

---

### Accessing the app

Once deployed, the app runs on **port 8080**:

```
http://<server-ip>:8080
```

To use a different port, edit `docker-compose.yml` before deploying — change the left side of the port mapping:

```yaml
ports:
  - "9000:8080"   # host port : container port
```

---

### Your data is safe

The database is stored in a Docker named volume (`garden-data`). It survives container restarts, redeployments, and app updates.

To back it up:

```bash
docker run --rm -v garden-data:/data -v $(pwd):/backup alpine \
  tar czf /backup/garden-backup.tar.gz /data
```

---

### Quick reference

| | |
|---|---|
| Default port | `8080` |
| Database path (inside container) | `/data/garden.db` |
| Docker volume | `garden-data` |
| Runs as | Non-root user (`garden`) |
