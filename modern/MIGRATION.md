# Migration Checklist

## Milestone 1: Vertical Slice (Playable)

- [x] TCP server accepts connections
- [x] Basic login with name prompt
- [x] Command loop
- [x] Core commands: look, say, who, help, quit
- [x] In-memory single-room world
- [x] Parser tests
- [x] One gameplay rule test (say broadcast)

## Milestone 2: Parser + Room System

- [x] Command table parity with [legacy/src/interp.c](../legacy/src/interp.c)
- [x] Room data structures and exits (from area files)
- [x] Room look formatting (from [legacy/src/act_info.c](../legacy/src/act_info.c))
- [x] Area metadata display (name, author)

## Milestone 3: Persistence

- [ ] Area loading/parsing (from [legacy/src/db.c](../legacy/src/db.c))
- [ ] Player save/load (from [legacy/src/save.c](../legacy/src/save.c))
- [ ] Reset logic (from [legacy/src/reset.c](../legacy/src/reset.c))

## Milestone 4: Combat + Skills

- [ ] Combat loop (from [legacy/src/fight.c](../legacy/src/fight.c))
- [ ] Skills/spells (from [legacy/src/skills.c](../legacy/src/skills.c))
- [ ] Equipment/objects handling (from [legacy/src/handler.c](../legacy/src/handler.c))

## Milestone 5: Timers + Updates

- [ ] Update loop and pulses (from [legacy/src/update.c](../legacy/src/update.c))
- [ ] Weather/time system (from [legacy/src/weather.c](../legacy/src/weather.c))
- [ ] Mob/obj progs (from [legacy/src/mud_prog.c](../legacy/src/mud_prog.c))

## Data Formats (Not Yet Parsed)

- Area files: [legacy/area/*.are](../legacy/area) (SMAUG area file format)
- Player files: [legacy/player/](../legacy/player) (classic pfile format)
- Boards: [legacy/boards/](../legacy/boards) (board text files)
