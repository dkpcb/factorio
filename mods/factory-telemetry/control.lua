-- Factorio World Model: observation-only prototype.
-- Output: <active Factorio server data>/script-output/factory-telemetry/events.jsonl
-- No gameplay mutation or automated construction.
local OUTPUT = "factory-telemetry/events.jsonl"
local MAX_QUEUE = 256

local function state()
  if not storage.factory_telemetry then
    storage.factory_telemetry = {sequence = 0, queue = {}}
  end
  return storage.factory_telemetry
end

local function emit(kind, tick, payload, surface)
  local s = state()
  s.sequence = s.sequence + 1
  s.queue[#s.queue + 1] = {
    schema_version = 1,
    -- Phase 1 deliberately does not promise globally unique world/run identity.
    -- Phase 2 will establish save lineage and reconcile snapshot/event history.
    world_id = "unverified",
    run_id = "unverified",
    sequence = s.sequence,
    tick = tick,
    surface_index = surface,
    event_type = kind,
    payload = payload
  }
  if #s.queue >= MAX_QUEUE then
    -- Do not silently drop on overflow. Flush immediately (bounded to 256 records).
    local lines = {}
    for _, event in ipairs(s.queue) do
      lines[#lines + 1] = helpers.table_to_json(event)
    end
    helpers.write_file(OUTPUT, table.concat(lines, "\n") .. "\n", true, 0)
    s.queue = {}
  end
end

local function flush()
  local s = state()
  if #s.queue == 0 then return end
  local lines = {}
  for _, event in ipairs(s.queue) do
    lines[#lines + 1] = helpers.table_to_json(event)
  end
  helpers.write_file(OUTPUT, table.concat(lines, "\n") .. "\n", true, 0)
  s.queue = {}
end

local function encode_entity(e)
  if not (e and e.valid) then return nil end
  return {
    unit_number = e.unit_number,
    name = e.name,
    type = e.type,
    x = e.position.x,
    y = e.position.y,
    direction = e.direction,
    force = e.force.name
  }
end

script.on_init(function()
  state()
  emit("mod_initialized", game.tick, {mod_version = "0.1.0"})
  flush()
end)

script.on_configuration_changed(function()
  state()
  emit("mod_configuration_changed", game.tick, {mod_version = "0.1.0"})
  flush()
end)

script.on_nth_tick(60, function()
  flush()
end)

script.on_event(defines.events.on_research_finished, function(event)
  local tech = event.research
  emit("research_finished", event.tick, {
    technology = tech.name,
    level = tech.level,
    force = tech.force.name,
    by_script = event.by_script
  })
end)

local function built(event)
  local entity = event.entity or event.created_entity
  local payload = encode_entity(entity)
  if payload then emit("entity_built", event.tick, payload, entity.surface.index) end
end

local function removed(event)
  local entity = event.entity
  local payload = encode_entity(entity)
  if payload then emit("entity_removed", event.tick, payload, entity.surface.index) end
end

script.on_event({
  defines.events.on_built_entity,
  defines.events.on_robot_built_entity,
  defines.events.script_raised_built,
  defines.events.script_raised_revive
}, built)

script.on_event({
  defines.events.on_player_mined_entity,
  defines.events.on_robot_mined_entity,
  defines.events.script_raised_destroy,
  defines.events.on_entity_died
}, removed)

-- On-demand limited snapshot; intentionally NOT a complete map scan.
-- In multiplayer run as an admin on the connected host.
commands.add_command("telemetry-snapshot",
  "Export player-force research and nearby entities (radius 32) without changing gameplay.",
  function(command)
    local player = command.player_index and game.get_player(command.player_index)
    if not player or not player.valid or not player.admin then
      return
    end

    local sid = tostring(game.tick) .. "-" .. tostring(command.player_index)
    local researched = {}
    for name, tech in pairs(player.force.technologies) do
      if tech.researched then researched[#researched + 1] = name end
    end
    table.sort(researched)
    emit("snapshot_research", game.tick, {
      snapshot_id = sid,
      force = player.force.name,
      technologies = researched,
      complete = true
    })

    local pos = player.position
    local area = {{pos.x - 32, pos.y - 32}, {pos.x + 32, pos.y + 32}}
    emit("snapshot_begin", game.tick, {
      snapshot_id = sid,
      scope = {kind = "area", x1 = area[1][1], y1 = area[1][2],
               x2 = area[2][1], y2 = area[2][2]},
      complete = false
    }, player.surface.index)

    local count = 0
    for _, entity in pairs(player.surface.find_entities_filtered{area = area, force = player.force}) do
      if count >= 300 then break end
      local payload = encode_entity(entity)
      if payload then
        payload.snapshot_id = sid
        emit("snapshot_entity", game.tick, payload, entity.surface.index)
        count = count + 1
      end
    end
    -- Snapshots capped at 300 entities cannot prove absence without a full scan.
    emit("snapshot_end", game.tick, {
      snapshot_id = sid, count = count, complete = false, limit = 300
    }, player.surface.index)
    flush()
    player.print("Factory Telemetry: exported limited snapshot (" .. count .. " entities).")
  end)
