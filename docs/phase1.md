# Phase 1 — Telemetry vertical slice (experimental)

Parent: [Issue #17](https://github.com/dkpcb/factorio/issues/17). Phase 0 design: `docs/architecture.md`, `docs/telemetry.md`.

## Scope implemented in this branch

- Thin `mods/factory-telemetry` Lua mod: `on_research_finished`, basic player/robot/script built and removed events, 60-tick batched JSONL file output, on-demand **limited** research + area snapshot command.
- `collector/`: dependency-free Go JSONL tail/validation CLI with tests. It **prints valid events to stdout**; **does not write SQLite** yet.
- Does **not** mutate gameplay: no placement, research or entity modification.
- The file is written by the **dedicated Factorio server process**, not automatically by the multiplayer client, to `script-output/factory-telemetry/events.jsonl` in that process's Factorio user data directory.

## Preparation (do NOT test on your main save)

1. Merge Phase 0 PR #18 before merging this branch, or review this PR as a stacked change.
2. On the Mac, inspect `~/factorio-llm-use/scripts/start_server.py`, the actual server `--mod-directory` / saves directories, and output path. These files cannot be read from GitHub and are not audited by this PR.
3. Install the telemetry MOD as a Factorio mod directory or correctly packaged ZIP in the **dedicated server's effective mods path**. Check mod-load and server logs before continuing. Use a throwaway save.
4. Verify that a simple gameplay event writes an NDJSON line to the **server's** script-output directory. If no file appears, verify the process's user-data path before changing Lua code.
5. From `collector/` run:
   ```sh
   go test ./...
   go run . -input "/absolute/path/to/server/script-output/factory-telemetry/events.jsonl"
   ```
   The collector shows one JSON object per complete input line. Use `-once` for a one-shot file read.
6. To test the limited initial snapshot, open the server game as an administrator and run `/telemetry-snapshot`. The snapshot covers only a **64 × 64 tile region** around that player, capped at 300 entities. `snapshot_end.complete=false` prevents treating missing entities as proof of absence.

## Explicit Phase 1 limitations

- `world_id` and `run_id` are deliberately `unverified`. Current JSONL ordering and `sequence` are not guaranteed safe across save reload/rollback. **Do not use this stream to authoritatively populate production DB state.**
- Collector's byte offset is only in memory: restarting it replays records. There is no persisted deduplication, transaction, checkpoint or SQLite integration; all are Phase 2.
- The snapshot is manually triggered and partial; event coverage is incomplete for some script entities/Space Age construction paths, and changes to settings or inventories aren't tracked.
- If the event output is truncated or overwritten, the collector resets offset with a warning; data integrity isn't yet ensured.
- A Factorio server may copy a save. Identify the **active server world** before comparing state with other sessions.
- Lua runtime API usage has been checked against the online 2.1.x docs; 2.1.24 runtime/mod interoperability is **not yet tested**.
- File writes and event extraction can affect UPS. No performance measurement has been performed.

## Acceptance criteria for Phase 1 before closing the phase

- [ ] Mod is loaded from the actual server's mod path, with no Lua errors.
- [ ] Research-complete event appears in JSONL and collector stdout.
- [ ] Player build/removal event appears, with surface and entity details.
- [ ] On-demand limited snapshot emits research + entity records correctly.
- [ ] `go test ./...` passes on the Mac.
- [ ] CPU/UPS overhead is measured on a throwaway save.
- [ ] Relevant unsupported event pathways are documented.

Once these are true, start Phase 2 on save/run identity, reconciliation, and durable SQLite.
