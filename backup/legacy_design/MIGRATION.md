# Migration Checklist

## Milestone 1: Vertical Slice (Playable) ✅

- [x] TCP server accepts connections
- [x] Basic login with name prompt
- [x] Command loop
- [x] Core commands: look, say, who, help, quit, stats, score
- [x] In-memory single-room world
- [x] Parser tests
- [x] One gameplay rule test (say broadcast)

## Milestone 2: Parser + Room System ✅

- [x] Command table parity with [legacy/src/interp.c](../legacy/src/interp.c)
- [x] Room data structures and exits (from area files)
- [x] Room look formatting (from [legacy/src/act_info.c](../legacy/src/act_info.c))
- [x] Area metadata display (name, author)
- [x] Movement commands (n/s/e/w/u/d + diagonals)
- [x] Autoexits command + toggle

## Milestone 3: Persistence ✅ (Phase 1 Complete)

### Phase 1: Core Persistence + Command System
- [x] Area loading/parsing - JSON converter (35 areas, 2178+ rooms)
- [x] Player model expansion (30+ fields: class, race, sex, level, attributes[7], vitals, combat stats)
- [x] Player save/load with extended schema
- [x] Character display (stats, score commands with legacy formatting)
- [x] Legacy player data loading (Zoie imported with full stats)
- [x] Command abbreviations/prefix matching (l→look, sa→say, st→stats, etc.)
- [x] Registry-based command dispatch with ordered prefix matching
- [x] Astat command stub for area statistics display
- [x] UX polish: Remove area author from standard look output (astat only)

### Phase 2: Skills System (In Progress)
- [x] Load skill definitions from JSON
- [x] Skill package with proficiency tracking
- [x] Player skill storage and persistence
- [x] Spellbook command to list learned skills
- [x] Cast command with mana/cooldown validation
- [ ] Skill damage/healing calculations
- [ ] Skill progression and proficiency leveling

### Phase 3: Character Creation (In Progress)
- [x] Race/class JSON loaders with menu formatting
- [x] Character creation flow (race → class selection)
- [x] Apply race attribute modifiers (STR/INT/WIS/DEX/CON/LCK/CHA)
- [x] Apply class modifiers (HP, Mana, hitroll)
- [x] Display final calculated stats
- [x] Smart login: new player detection + creation flow vs existing player load
- [x] Sex selection during creation (Male/Female/Neutral with confirmation)
- [ ] Attribute assignment/customization
- [ ] Starting skill proficiencies

### Phase 4: Advanced Persistence (In Progress)
- [ ] Reset logic (from [legacy/src/reset.c](../legacy/src/reset.c))
	- [x] Initial spawn from resets (no respawn yet)
- [x] Reset scheduler (per-area timers with a default interval)
  - [x] Per-area reset_minutes field in area JSON
  - [x] RespawnTick implementation with per-area timing
  - [x] Mob/object respawn (clears and re-instantiates from resets)
- [ ] Object/item persistence (inventory system)

## Milestone 4: Combat + Skills (Not Started)

- [ ] Combat loop (from [legacy/src/fight.c](../legacy/src/fight.c))
- [ ] Skills/spells (from [legacy/src/skills.c](../legacy/src/skills.c))
- [ ] Equipment/objects handling (from [legacy/src/handler.c](../legacy/src/handler.c))

## Milestone 5: Timers + Updates (Not Started)

- [ ] Update loop and pulses (from [legacy/src/update.c](../legacy/src/update.c))
- [ ] Weather/time system (from [legacy/src/weather.c](../legacy/src/weather.c))
- [ ] Mob/obj progs (from [legacy/src/mud_prog.c](../legacy/src/mud_prog.c))

## Data Format Conversions ✅

- [x] Area files: [legacy/area/*.are](../legacy/area) → [areas/*.json](../areas) (35 files, 2178+ rooms, mobiles, objects)
- [x] Race definitions: [legacy/races/*.race](../legacy/races) → [races/*.json](../races) (26 races)
- [x] Class definitions: [legacy/classes/*.class](../legacy/classes) → [classes/*.json](../classes) (14 classes)
- [ ] Player files: [legacy/player/](../legacy/player) (classic pfile format) → modern JSON
- [ ] Boards: [legacy/boards/](../legacy/boards) (board text files) → structured format
- [ ] Object definitions: [legacy/src/obj/](../legacy/src/obj) (if present)
- [ ] Mobile definitions: [legacy/src/mob/](../legacy/src/mob) (if present)

## Runtime Config

- [x] Config file for startup settings (config.json)
	- [x] Start room override (start_room_vnum)
	- [x] Default respawn interval (respawn_default_minutes)
