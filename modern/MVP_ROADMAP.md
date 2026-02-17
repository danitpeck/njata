# NJATA Skills System: Minimal Viable Design

**MVP Progress: ~70% Complete** | Last Updated: Feb 17, 2026

**✅ DONE**: Spell system (8 spells), classless design, starter kits, character creation, Study command (MVP), save/load  
**⏳ IN PROGRESS**: Combat maneuvers (Slash), magical item placement, simple combat resolution  
**📋 NEXT**: Implement combat maneuver system (Option A)

---

## Philosophy

Build small, add meaningful. Like the original njata.c approach:
- Start with only what's needed (6-8 core spells)
- Create clean extension points for custom additions
- Grow organically based on actual gameplay needs, not pre-designed systems
- Each addition should have purpose within the game world

---

## Part 1: Minimal Viable Product (MVP) - Week 1

**PRIORITY: Discovery-Based Progression**

The absolute minimum to have a working spell system for gameplay testing—with discovery-based learning as the centerpiece.

### What Goes In

**Core Mechanic: Study Ability**

Study is NJATA-specific progression:
```
All players can learn spells by studying magical items (wands, scrolls)
- Discover items around the world
- Make proficiency check to learn the spell inside
- Items are consumed (create scarcity and value)
- Study skill improves with practice
- Different from auto-learning—feels like achievement
```

This changes from "generic eight spells" to **"eight spells + discovery-driven learning"**

**Supporting 8-Spell Roster** (Colista-aligned, anyone can learn):
```
1001: Arcane Bolt      (core offensive, foundation magic)
1002: Leviathan's Fire (Immak fire, deep sea magic)
1003: Mend             (life magic, healing)
1004: Shadow Veil      (darkness + illusion, debuffs and control)
1005: Ephemeral Step   (fairy magic, escape and stealth)
1006: Path Shift       (the Courier's secret, teleportation)
1007: Winter's Whisper (mountain/water magic, control through cold)
1008: Knowing          (Bronze Artisan's insight, utility)
```

**Core Systems**:
- ✅ Spell definitions (8 spells, JSON format)
- ✅ Spell loading (skills.Load())
- ✅ Proficiency tracking (30-100%, increases with casting)
- ✅ Mana cost & cooldown validation
- ✅ Basic damage formula (dice + attributes)
- ✅ spellbook command (view learned spells with proficiency)
- ✅ cast command (execute spell, improve proficiency)
- ✅ Study command (MVP version with keyword mapping)
- ✅ Player save/load (persist skills, proficiency, learned status)
- ✅ Classless design (starter kits replace classes)
- ✅ Character creation (race + starter kit + age + sex)
- ⏳ 8 magical items with embedded spell IDs (need to create & place)
- ⏳ Combat maneuver system (Slash defined in kits, needs implementation)
- ⏳ Trainer NPCs (for teaching maneuvers)

### Progression System: Equipment + Skills (No Levels)

**Core Philosophy**: Players grow stronger through **equipment improvement** + **skill mastery**, not experience points or levels.

**How It Works**:

**Equipment Progression** (Self-reinforcing loop):
- Weak monsters drop weak gear (low damage, low armor)
- Player defeats weak monsters, loots their gear
- Better gear enables fighting stronger monsters
- Stronger monsters drop better gear
- Cycle repeats: gear → monsters → gear
- Natural difficulty scaling without artificial gates

**Skill Progression** (Proficiency-based):
- Weapons/spells have **proficiency** (0-100%)
- Each use slightly increases proficiency
- Higher proficiency = higher damage multiplier/effectiveness
- Example: Longsword at 10% skill vs 90% skill makes huge difference
- **Scaling formula**: `base_damage * (1.0 + proficiency * 0.5)` (scales to 150% at 100%)

**Combat-Focused Player Example**:
1. Start: basic sword, Slash proficiency 5%
2. Fight goblins → loot dagger, short sword, leather armor
3. Equip better gear, improve Slash proficiency to 15%
4. Fight stronger goblins → loot hand axe, chain mail
5. Discover trainer in Alklyu → learn Power Attack maneuver
6. Use Power Attack, proficiency rises to 30%
7. Hunt harder areas → find greatsword
8. Repeat cycle indefinitely—always stronger monsters, better gear, improved skills

**Magic-Focused Player Example**:
1. Start: Arcane Bolt (from starter kit or early Study), Study skill 5%
2. Find Wand of Leviathan's Fire → Study DC 50
3. Succeed → Learn Leviathan's Fire at 30% proficiency
4. Cast spells in combat → improve Arcane Bolt to 40%, Leviathan's to 50%
5. Hunt deeper areas for rare scrolls (Knowing, Winter's Whisper)
6. Study better items → expand repertoire
7. Powerful spells with high proficiency = can tackle harder content
8. Exploration cycle: hunt → study → grow → hunt harder areas

**Hybrid Player Example**:
1. Start with sword + basic spell
2. Learn maneuvers from trainers, study magical items when found
3. Equip STR/DEX gear for combat, keep INT-focused backup for spells
4. Switch tactics based on enemy types (physical vs. magic-vulnerable)
5. Natural specialization emerges from gear choices and time investment

**Why This Works**:
- No artificial plateaus (no "soft caps")
- Equipment drops are tangible rewards
- Skill feels earned through use
- Both matter: bad gear limits early progress, low skill makes late-game difficult
- Encourages player choice: focus combat, focus magic, or mix both
- Natural specialization through gear and time investment (INT gear vs STR gear)
- Scales naturally from newbie to veteran

**Systems to Skip (for now)**:
- ❌ Components system
- ❌ Multiple save types (just Reflex)
- ❌ Complex damage scaling
- ❌ Effect stacking
- ❌ Cooldown categories
- ❌ Area selection complexity
- ❌ NPC teaching complexity (Study and trainers are simple MVP versions)
- ❌ Item crafting (can add later)

### MVP Data Structures

**Spell Definition** (Minimal):
```json
{
  "id": 1002,
  "name": "Leviathan's Fire",
  "description": "Fire magic born from the depths of Immak's naiad cove.",
  "mana_cost": 30,
  "cooldown_seconds": 5,
  "level_required": 1,
  
  "targeting": {
    "mode": "hostile_area",
    "range": 30,
    "radius": 3
  },
  
  "effects": {
    "damage": "4d8 + I",
    "damage_type": "fire",
    "save_type": "reflex",
    "save_dc": 14
  },
  
  "messages": {
    "cast": "You hurl Leviathan's fire at $target.",
    "hit": "$actor's flame scorches you!",
    "miss": "$actor's fire sails past you!",
    "save": "$target resists the flames!"
  }
}
```

**Magical Item** (For Study):
```json
{
  "vnum": 5002,
  "short_descr": "Wand of Leviathan's Fire",
  "long_descr": "A glowing wand crackles with deep, ancient fire.",
  "item_type": "WAND",
  "values": [
    0,    // value[0] charges (cosmetic)
    0,    // value[1] unused
    0,    // value[2] unused
    1002  // value[3] SPELL ID for Leviathan's Fire
  ]
}
```

### MVP Go Code (Minimal)

```go
// internal/spells/spell.go
type Spell struct {
    ID            int
    Name          string
    ManaCost      int
    CooldownSecs  int
    Damage        string  // "4d8 + I"
    DamageType    string  // "fire"
    SaveDC        int
    Range         int
    AreaRadius    int
    Messages      map[string]string
}

// internal/commands/cast.go - Very simple
func CastSpell(p *Player, args []string) {
    spell := GetSpell(args[0])
    
    // Simple checks
    if p.Mana < spell.ManaCost { error }
    if !IsReadyToCast(p, spell.ID) { error }
    if !p.KnowsSpell(spell.ID) { error }
    
    // Roll damage
    damage := RollDice(spell.Damage)
    
    // Apply to targets
    for _, target := range GetTargets(p, spell) {
        if target.MakeSave(spell.SaveDC) {
            damage = damage / 2
        }
        target.TakeDamage(damage)
    }
    
    p.Mana -= spell.ManaCost
    SetCooldown(p, spell.ID, spell.CooldownSecs)
}

// internal/commands/study.go - Available to all players
func StudyItem(player *Player, itemName string) {
    // Find item in room
    item := GetItemInRoom(player, itemName)
    if item == nil { error }
    
    // Validate it contains a spell
    spellID := item.Values[3]
    if spellID == 0 { error }
    
    // Make Study skill check
    studyProficiency := player.Skills[STUDY_SKILL_ID]
    dc := 55 - int(float64(studyProficiency) * 0.8)
    roll := Rand(1, 100)
    
    if roll < dc {
        // Failed
        Send(player, "You cannot glean any knowledge from it.")
        player.Skills[STUDY_SKILL_ID] += 1  // +1% on failure
        item.Remove()
        return
    }
    
    // Success
    if player.KnowsSpell(spellID) {
        Send(player, "You already know that spell!")
        return
    }
    
    player.LearnSpell(spellID, 30)  // 30% proficiency
    player.Skills[STUDY_SKILL_ID] += 2  // +2% on success
    item.Remove()
    Send(player, 
        fmt.Sprintf("You have learned the art of %s!", GetSpellName(spellID)))
}
```

### MVP Testing

```go
// Spell casting tests
✓ Can cast spell with enough mana
✓ Spell fails with insufficient mana
✓ Cooldown prevents quick re-cast
✓ Damage dealt to target

// Study tests (available to all players)
✓ Can study magical item in room
✓ Item must contain spell (value[3])
✓ Study success check works (DC vs proficiency)
✓ On success: Learn spell at 30%, item consumed
✓ On failure: No spell learned, item consumed
✓ Can't learn same spell twice
✓ Study proficiency increases with use

// Starter kit tests
✓ Character creation shows three starter kits
✓ Scholar's Kit starts with Arcane Bolt (30%) + Study (10%)
✓ Warrior's Kit starts with Slash maneuver (10%)
✓ Wanderer's Kit starts with Arcane Bolt (20%) + Study (5%) + Slash (5%)
✓ Final stats display shows race and sex (no class/kit after creation)

// Combat maneuver tests (TODO)
☐ Can use 'slash <target>' command to attack
☐ Slash proficiency improves with use
☐ Damage scales with proficiency (like spells)
☐ Simple mob combat works (goblins fight back)
```

### Deliverable: MVP

**Current Status: ~70% Complete**

**What Ships (MVP):**
- ✅ 8 core spells (Colista-themed, anyone can learn)
- ✅ Starter kit selection (Scholar/Warrior/Wanderer)
- ✅ Study command (learn spells from items - MVP keyword version)
- ✅ Cast command with proficiency improvement
- ✅ Spellbook command (view learned spells)
- ✅ Classless design (no more class restrictions)
- ✅ Player persistence (skills, proficiency, learned spells)
- ⏳ Combat maneuver system (Slash + basic combat)
- ⏳ 8 magical items placed in world
- ⏳ Simple combat testing (fight goblins, improve Slash proficiency)

### Combat Maneuver MVP Specification

**Goal**: Warrior players can use Slash maneuver to fight, improving proficiency through use (parallel to Scholar's spell casting).

**Minimum Implementation**:
1. **Slash Maneuver** (skill ID 9002):
   - Command: `slash <target>` or `maneuver slash <target>`
   - Base damage: 1d6 + STR/2 + proficiency_bonus
   - Proficiency scaling: same as spells (base_damage * (1.0 + proficiency * 0.5))
   - Improves 1% per use
   - No mana cost, no cooldown (basic attack)

2. **Simple Combat Resolution**:
   - Target selection (mob in room)
   - Damage calculation
   - Mob HP tracking
   - Mob death / loot (optional for MVP)
   - Combat messages (you slash, mob takes damage)

3. **Test Scenario**:
   - Create goblin mob (HP: 20)
   - Warrior uses `slash goblin`
   - Goblin takes damage
   - Slash proficiency increases from 10% → 11%
   - Repeat until goblin dies

**Deferred Features**:
- Mob counter-attacks (AI)
- Multiple maneuvers (Power Attack, Defensive Stance, etc.)
- Cooldowns on maneuvers
- Equipment bonuses to damage
- Combat rounds/initiative
- Area attacks
- Advanced targeting

**Implementation Checklist**:
```
[ ] Define Slash as a "maneuver" type skill in skills.json or separate system
[ ] Create cmdSlash or cmdManeuver command in player_commands.go
[ ] Add damage calculation (1d6 + STR/2, scale with proficiency)
[ ] Update proficiency on each use (+1%)
[ ] Create simple mob HP system (extend game.Mobile)
[ ] Add combat resolution (deal damage, check if mob dies)
[ ] Add combat messages (you slash, mob takes X damage)
[ ] Test: Warrior can kill goblin using slash command
[ ] Test: Slash proficiency improves from 10% → 11% → 12% with use
```

**Deferred to Post-MVP:**
- Equipment system (gear drops, wear slots, stat focus)
- Trainer NPCs (can add advanced maneuvers later)
- Study proficiency checks (DC system, success/failure)
- Item consumption on Study
- Complex combat (targeting, mob AI, tactics)

**Why Classless Works** (✅ Validated through implementation): 
- **Player agency**: Discover your own path, not locked into pre-made archetypes
- **Natural specialization**: Starter kits provide direction without permanent restrictions
- **Exploration incentive**: Both combat and magic paths reward exploration (trainers, spell items)
- **Lower cognitive load**: No "what if I picked the wrong class?" anxiety
- **Flexible progression**: Can pivot playstyle as gear allows
- **Dev's brilliant insight**: Classes were artificial constraints on player creativity 💯

**Progression Philosophy**:
- **No levels** - skill-based advancement only
- Two primary progression paths (not mutually exclusive):
  - Magic: Study items → learn spells → improve spell proficiency via casting
  - Combat: Seek trainers → learn maneuvers → improve proficiency via fighting
- Different incentive structures create natural specialization:
  - Magic-focused: Must explore to find spell items, rewards curiosity
  - Combat-focused: Can fight immediately, optionally explore for advanced trainers
  - Hybrid: Balanced approach, gear flexibility required (INT and STR stats)

**Deferred Content**:
- Additional starter kits (merchant, performer, etc.)
- Spirit-communing mechanics tied to Colistani animism (Fox/Owl/Deer/Snake spirits)
- Performance-based inspiration system (bard-style buffing)

**Time estimate**: 3-4 days for one developer (includes 8 item definitions)

---

## Part 2: After MVP - Growth Points

Once MVP is working, extend based on actual need, not pre-planning.

### Growth Point 1: More Spells (Week 2-3)

**When**: After playing with 8 spells, identify gaps

**Examples**:
- "We need poison damage" → Add Poison spell
- "Healing is OP, need diminishing returns" → Add cooldown scaling
- "Mages need utility" → Add Dispel Magic
- "Rangers need escape" → Add Smoke Bomb

**Philosophy**: Each new spell solves a specific gameplay problem observed during play, not from a pre-made list.

### Growth Point 2: Components System

**When**: After spells feel too cheap/valuable

**What**: Add cost beyond mana
```
"This rare spell needs rare components"
→ Gate spell: 50 gold + rare reagent
```

**Trigger**: Only if game balance requires gatekeeping powerful spells

### Growth Point 3: Complex Damage Types

**When**: Combat feels flat (after 100+ test sessions)

**What**: Add resistances, damage scaling
```
"This dragon is immune to fire"
"Cold-based spell does 2x to fire creatures"
```

**Trigger**: When you notice certain spells always dominate

### Growth Point 4: Spell Effects (Buffs/Debuffs)

**When**: After 50+ hours of gameplay

**Systems**:
- Strength buff: +2 STR for 60 seconds
- Blindness: AC penalty for 30 seconds
- Haste: +25% attack speed for 45 seconds

**Trigger**: When combat becomes repetitive (just damage/heal cycling)

### Growth Point 5: Advanced Mechanics

**When**: 200+ hours logged, core gameplay solid

**Examples**:
- Spell interruption (damage breaks casting)
- Channeled spells (cast over time)
- Spell chains (use one spell to power up another)
- Mana shield (convert HP to mana temporarily)

**Trigger**: When experienced players want depth

### Growth Point 6: Progression Paths (Architecture)

**Combat Path**:
- **Combat Maneuvers** - Tactics learned from trainers (Vojvoda Zsa's academy, combat masters in towns)
- Learn basic techniques (Power Attack, Defensive Stance, Guard)
- **Lower barrier to entry**: grab a weapon, learn basic moves, fight monsters
- Maneuver proficiency improves through combat use
- Can optionally explore to find distant trainers for advanced techniques
- Simple, self-directed loop: fight → improve → fight better

**Magic Path**:
- **Spell Discovery** - learn new spells via Study
- Must actively hunt for magical items
- Can't just cast Arcane Bolt forever and plateau—exploration is built-in progression
- Proficiency improves through casting and successful studying
- Engaging, exploration-driven loop: hunt → study → expand repertoire
- Higher cognitive load but more rewarding for curious players

**Design Principle**: Both skill-based. Different incentive structures:
- Combat-focused: fight monsters first, optionally explore for trainers
- Magic-focused: must explore to stay viable
- Hybrid: balance time between both paths, gear flexibility required

### Growth Point 7: Content Expansion (After MVP)

**Additional Starter Kits**:
- Performer's Kit: Simple instrument, basic inspiration ability (buff allies)
- Merchant's Kit: Trading goods, appraisal skill, negotiation focus
- Hunter's Kit: Bow, tracking skill, wilderness survival

**Spirit-Touched Mechanics** (Post-MVP, designed from animism):
Tie to Colistani animal spirits (Fox/Owl/Deer/Snake):
- Fox-touched: hunting/tracking bonuses, spirit favor system
- Owl-touched: wisdom/utility, vision in darkness
- Deer-touched: blessing/warding, protective magic
- Snake-touched: underworld/shadow magic, stealth bonuses

**Design Philosophy**: Each new system needs:
1. Lore grounding (tied to world, not generic archetype)
2. Unique mechanic (not just "different spell list")
3. Actual design document (not template placeholders)

**What We're NOT Doing**: Avoid generic fantasy archetypes. Only add systems that fit NJATA's world (Wislyu, Colista, animal spirits).

---

## Part 3: Design for Growth (Architecture)

These "growth points" work because of design choices made in MVP:

### Extensible Schema

```json
// MVP: Minimal required fields
{
  "id": 1001,
  "name": "Arcane Bolt",
  "mana_cost": 15,
  "cooldown_seconds": 2,
  "damage": "1d6 + I",
  "messages": {...}
}

// Growth: Easy to add fields without breaking old spells
// Add to new spells only: components, resistances, effects, etc.
{
  "id": 2001,
  "name": "Gate Portal",
  "mana_cost": 75,
  "cooldown_seconds": 60,
  "components": [{"type": "gold", "amount": 50}],
  "damage": null,
  "effects": [{"type": "teleport"}],
  "messages": {...}
}
```

### Plugin Architecture

```go
// Core spell system doesn't know about extensions
// But can call hooks during casting

type SpellHook interface {
    OnCast(caster, target, spell)
    OnHit(caster, target, spell, damage)
    OnEffect(caster, target, spell, effect)
}

// New spell systems just register hooks
RegisterHook("resistances", ResistanceHook)
RegisterHook("effects", EffectHook)
RegisterHook("proficiency", ProficiencyHook)
```

### Separation of Concerns

```
MVP (Core, stable):
  ├─ Spell loading
  ├─ Mana/cooldown validation
  ├─ Basic damage rolling
  └─ Message output

Growth 1 (Damage types):
  ├─ Resistance checking
  ├─ Typed damage system
  └─ Modifier application

Growth 2 (Effects):
  ├─ Affect tracking
  ├─ Duration management
  └─ Stat modification

Growth 3 (Complexity):
  ├─ Chains/combos
  ├─ Interruption
  └─ Advanced targeting
```

---

## Part 4: Recommended Week-by-Week

### Week 1: MVP Foundation

**Completed:**
- ✅ Spell definitions (8 spells: Arcane Bolt, Leviathan's Fire, Mend, Shadow Veil, Ephemeral Step, Path Shift, Winter's Whisper, Knowing)
- ✅ Spell loader (skills/skills.json)
- ✅ Cast command (with proficiency improvement)
- ✅ Spellbook command (displays learned spells)
- ✅ Study command (MVP keyword-based version)
- ✅ Damage calculation (dice + attribute modifiers)
- ✅ Proficiency tracking (0-100%, increases on use)
- ✅ Player persistence (JSON save/load with skills)
- ✅ Classless character creation (removed class system)
- ✅ Starter kit selection (Scholar/Warrior/Wanderer)
- ✅ Test: All 8 spells castable and tracked

**In Progress:**
- ⏳ Combat maneuver system:
  - Define Slash maneuver (basic attack)
  - Create maneuver command (like cast, but for combat)
  - Proficiency improvement through use
  - Simple AI/combat for testing
- ⏳ 8 magical items with spell IDs:
  - Define items in JSON with value[3] = spell_id
  - Place in world areas for discovery
  - Update Study command to check for actual items in room
- ⏳ Optional: Trainer NPCs for teaching advanced maneuvers

### Week 2: Polish & Test
- [ ] Run 20+ hours of gameplay
- [ ] Identify balance issues
- [ ] Identify missing features

### Week 3: First Growth (Based on Observations)
- [ ] Add 2-4 spells players wanted
- [ ] Fix balance issues found
- [ ] Add 1 system based on testing (e.g., resistances)

### Weeks 4-8: Organic Evolution
- [ ] Add features as gameplay demands
- [ ] Each addition answering a real need
- [ ] Regular playtesting between additions

This is better than trying to implement 28 spells + full component system + multiple save types + effect stacking all at once.

---

## Part 5: Design Artifacts (Reference)

The full `SKILLS_SYSTEM_DESIGN.md` document exists for reference when you need to add features:

- **Need a new spell effect?** → Check "Spell Effects & Affects System" 
- **Want proficiency scaling?** → Check "Proficiency & Learning System"
- **Need resistances?** → Check "Damage Type Resistances"
- **Want NPC teachers?** → Check "Learning Mechanics"
- **Building components system?** → Check "Spell Components System"

It's a **reference library, not a spec**. Use it to answer design questions as they come up, not as a prescriptive roadmap.

---

## Part 6: Custom Content Hooks

Like njata.c, add custom mechanics for specific races/classes:

### Example: Dragon Breath (Racial Ability, Not Generic Spell)

```go
// Not a spell - unique dragon ability
func (p *Player) DragonBreath(target *Character) {
    if p.Race != RACE_DRAGON {
        return "Only dragons can breathe fire!"
    }
    
    damage := RollDice("6d8 + I")  // Better than spells
    target.TakeDamage(damage)
}
```

### Example: Unicorn's Healing Presence (Racial Aura)

```go
// Passive effect, not a castable spell
func (p *Player) UnicornPresence() {
    if p.Race == RACE_UNICORN {
        p.PassiveHeal = 2  // 2 HP per tick
    }
}
```

### Example: Fairy's Luck Spell (Class/Race Hybrid)

```go
// Custom spell tied to specific race
func CastFairyLuck(p *Player, args []string) {
    if p.Race != RACE_FAIRY {
        SendMessage(p, "Only fairies can use this magic!")
        return
    }
    
    // Unique behavior to fairies only
    for i := 0; i < 3; i++ {
        crit := RollDice("1d20") + p.Stats.Luck
        if crit > 18 {
            p.Gold += 100
        }
    }
}
```

This is how NJATA grows: race-specific abilities, not generic spell lists.

---

## Summary: MVP → Growth → Custom

```
┌─────────────────────────────────────────────────────────┐
│ MVP: 8 spells, simple system, play for a week           │
├─────────────────────────────────────────────────────────┤
│ Observe gameplay, identify needs                        │
├─────────────────────────────────────────────────────────┤
│ Growth 1: Add features based on observations            │
│ (resistances? effects? new spells?)                     │
├─────────────────────────────────────────────────────────┤
│ Repeat growth cycle as game evolves                     │
├─────────────────────────────────────────────────────────┤
│ Custom Content: Race/class unique mechanics             │
│ (not generic spells, but njata-specific flavor)        │
└─────────────────────────────────────────────────────────┘
```

**Start small. Add based on real needs. Keep njata feeling custom and alive.**

