# Telemetry Data Contract — Phase 0 draft

Status: **proposed**; validate with the Factorio 2.1.24 runtime. Parent Issue #17. This is a contract specification, **not an implemented feature**.

## Pipeline

Factorio Lua mod -> `script-output/factory-telemetry/events.jsonl` -> local Go Collector -> SQLite. One object per JSONL line; use `helpers.table_to_json` for serialization where appropriate. Collector must tolerate incomplete trailing lines and must not advance its checkpoint past them.

## Minimal event envelope

```json
{
  "schema_version": 1,
  "world_id": "provisional-save-identity",
  "run_id": "provisional-session-id",
  "sequence": 42,
  "tick": 123456,
  "surface_index": 1,
  "event_type": "research_finished",
  "payload": {
    "force": "player",
    "technology": "electronics",
    "level": 1
  }
}
```

All IDs and numbers are **examples only**. `world_id` / `run_id` generation and host/server routing remain unverified. Schema must be versioned; unknown fields should be tolerated within a version. The Collector validates type and required fields before storage. Include `mod_version`, `factorio_version`, and mod-set fingerprint in a session-start record when implementable.

## Event families (initial scope)

| Event | Essential fields | Notes |
| --- | --- | --- |
| `session_started` | run/world IDs, observed tick, mod/version fingerprint | Not equivalent to unique save identity; revalidate rollback |
| `research_finished` | force, technology name, level, tick | Multi-force and repeated research must be supported |
| `entity_built` | entity id, prototype, position, direction, force, surface | Player, robot, script-raised variants |
| `entity_removed` | entity id, surface, reason | Record before invalidation when necessary |
| `snapshot_begin` | snapshot id, declared scope, tick | Establish scope only; do not delete any DB entities yet |
| `snapshot_entity` | snapshot id, entity identity and properties | Bounded export or paginated batches |
| `snapshot_end` | snapshot id, tick, item count, checksum (optional) | Mark complete only after all parts received |

**Do not assume** all possible runtime entity changes trigger the same event; Phase 1 must produce an event coverage matrix. Changes to recipes/configuration/force merges and runtime scripted manipulation need separate review.

## Initial SQLite sketch

```sql
CREATE TABLE runs (
  run_id TEXT PRIMARY KEY,
  world_id TEXT NOT NULL,
  started_at TEXT,
  first_tick INTEGER NOT NULL,
  last_tick INTEGER,
  reconciliation_required INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE events (
  run_id TEXT NOT NULL,
  sequence INTEGER NOT NULL,
  tick INTEGER NOT NULL,
  event_type TEXT NOT NULL,
  payload_json TEXT NOT NULL,
  PRIMARY KEY (run_id, sequence)
);

CREATE TABLE research (
  run_id TEXT NOT NULL,
  force_name TEXT NOT NULL,
  technology_name TEXT NOT NULL,
  level INTEGER,
  researched INTEGER NOT NULL,
  observed_tick INTEGER NOT NULL,
  PRIMARY KEY (run_id, force_name, technology_name)
);

CREATE TABLE sync_checkpoints (
  source_id TEXT PRIMARY KEY,
  byte_offset INTEGER NOT NULL,
  last_run_id TEXT,
  last_sequence INTEGER
);
```

This is an MVP sketch, not final DDL. Entity and snapshot tables belong to the Phase 1/2 detailed design. The Collector should apply parsed events and move checkpoints **in a single transaction**. File rotation and append/rewrite recovery need file identity in the checkpoint, not only offset.

## Recovery invariants

1. **No phantom deletes:** missing entity in a partial snapshot says nothing about the world.
2. **No cross-save contamination:** never merge events or current state solely by entity number/coordinates.
3. **No replay inflation:** re-reading an event is idempotent by `(run_id, sequence)`.
4. **No unverified freshness:** show source run, tick, capture scope and sync status in every AI query.
5. **Safe rollback:** after any load, copied-save replacement, or run identity ambiguity, mark state untrusted until a complete reconciling snapshot.
6. **Unobserved is unknown**, not zero production or no resources.

## Mandatory acceptance tests

- Research completion appears once in SQLite even after Collector restarts/replays file.
- A full snapshot plus subsequent build/remove events matches the game for a limited test region.
- Partial JSONL line, truncated export, and interrupted ingestion do not corrupt DB state.
- Save rollback, load of a different map and copied-save/AI-server resume invalidate previous-current projections correctly.
- Collecting does not visibly impair UPS for a small production test case (establish measurement thresholds during Phase 1).
- Telemetry has no write effects on gameplay beyond its own saved bookkeeping state and external output.
