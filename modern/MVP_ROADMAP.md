# NJATA Skills System: Minimal Viable Design

## Philosophy

Build small, add meaningful. Like the original njata.c approach:
- Start with only what's needed (6-8 core spells)
- Create clean extension points for custom additions
- Grow organically based on actual gameplay needs, not pre-designed systems
- Each addition should have purpose within the game world

---

## Part 1: Minimal Viable Product (MVP) - Week 1

## Part 1: Minimal Viable Product (MVP) - Week 1

**PRIORITY: Scholar's Study System**

The absolute minimum to have a working spell system for gameplay testing—with Scholar's unique discovery mechanics as the centerpiece.

### What Goes In

**One Spotlight: Scholar Class + Study Ability**

Scholar is NJATA-specific custom content:
```
Scholars learn spells by studying magical items (wands, scrolls)
- Discover items around the world
- Make proficiency check to learn the spell inside
- Items are consumed (create scarcity and value)
- Study skill improves with practice
- Different from auto-learning—feels like achievement
```

This changes from "generic eight spells" to **"eight spells + how Scholar uniquely learns them"**

**Supporting 8-Spell Roster** (Colista-aligned, available to both Scholar and Warrior):
```
1001: Arcane Bolt      (core offensive, foundation magic)
1002: Leviathan's Fire (Immak fire, Scholars hunt in libraries)
1003: Mend             (life magic, both classes can learn)
1004: Shadow Veil      (darkness + illusion, debuffs and control)
1005: Ephemeral Step   (fairy magic, escape and stealth)
1006: Path Shift       (the Courier's secret, teleportation)
1007: Winter's Whisper (mountain/water magic, control through cold)
1008: Knowing          (Bronze Artisan's insight, utility)
```

**Core Systems**:
- ✅ Spell definitions (JSON, minimal)
- ✅ Spell loading
- ✅ Proficiency tracking (30-100%, learned through Study or auto-learn)
- ✅ Mana cost & cooldown validation
- ✅ Basic damage formula (no complex scaling yet)
- ✅ spellbook command (view learned spells)
- ✅ cast command (execute spell)
- ✅ **NEW: Study command** (Scholar-specific, central feature)
- ✅ Study skill proficiency tracking
- ✅ Spell learning from studied items
- ✅ Simple persistence (player file)
- ✅ 8 magical items with embedded spell IDs

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

**Warrior Example Loop**:
1. Start: basic sword, Slash proficiency 5%
2. Fight goblins → loot dagger, short sword, leather armor
3. Equip better gear, improve Slash proficiency to 15%
4. Fight stronger goblins → loot hand axe, chain mail
5. Discover trainer in Alklyu → learn Power Attack maneuver
6. Use Power Attack, proficiency rises to 30%
7. Hunt harder areas → find greatsword
8. Repeat cycle indefinitely—always stronger monsters, better gear, improved skills

**Scholar Example Loop**:
1. Start: Arcane Bolt (learned via training), Study skill 5%
2. Find Wand of Leviathan's Fire → Study DC 50
3. Succeed → Learn Leviathan's Fire at 30% proficiency
4. Cast spells in combat → improve Arcane Bolt to 40%, Leviathan's to 50%
5. Hunt deeper areas for rare scrolls (Knowing, Winter's Whisper)
6. Study better items → expand repertoire
7. Powerful spells with high proficiency = can tackle harder content
8. Exploration cycle: hunt → study → grow → hunt harder areas

**Why This Works**:
- No artificial plateaus (no "soft caps")
- Equipment drops are tangible rewards
- Skill feels earned through use
- Both matter: bad gear limits early progress, low skill makes late-game difficult
- Encourages both combat (Warrior) and exploration (Scholar)
- Scales naturally from newbie to veteran

**Systems to Skip (for now)**:
- ❌ Components system
- ❌ Multiple save types (just Reflex)
- ❌ Complex damage scaling
- ❌ Effect stacking
- ❌ Cooldown categories
- ❌ Area selection complexity
- ❌ NPC teaching (except Scholar's unique study mechanic)
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

// internal/commands/study.go - NEW for Scholar
func StudyItem(scholar *Player, itemName string) {
    // Find item in room
    item := GetItemInRoom(scholar, itemName)
    if item == nil { error }
    
    // Validate it contains a spell
    spellID := item.Values[3]
    if spellID == 0 { error }
    
    // Make Study skill check
    studyProficiency := scholar.Skills[STUDY_SKILL_ID]
    dc := 55 - int(float64(studyProficiency) * 0.8)
    roll := Rand(1, 100)
    
    if roll < dc {
        // Failed
        Send(scholar, "You cannot glean any knowledge from it.")
        scholar.Skills[STUDY_SKILL_ID] += 1  // +1% on failure
        item.Remove()
        return
    }
    
    // Success
    if scholar.KnowsSpell(spellID) {
        Send(scholar, "You already know that spell!")
        return
    }
    
    scholar.LearnSpell(spellID, 30)  // 30% proficiency
    scholar.Skills[STUDY_SKILL_ID] += 2  // +2% on success
    item.Remove()
    Send(scholar, 
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

// Scholar Study tests (NEW)
✓ Can study magical item in room
✓ Item must contain spell (value[3])
✓ Study success check works (DC vs proficiency)
✓ On success: Learn spell at 30%, item consumed
✓ On failure: No spell learned, item consumed
✓ Can't learn same spell twice
✓ Study proficiency increases with use
```

### Deliverable: MVP

**Two Classes**:
- **Scholar** - spell learning via Study mechanic (discovery-based)
- **Warrior** - melee combat foundation

**What Ships**:
- 8 core spells usable by both classes
- 8 magical items with embedded spell IDs
- **Scholar's Study ability** (Scholar-exclusive) - discover items, make checks, learn spells
- Study skill proficiency tracks and improves with use
- **Warrior Combat Maneuvers** - learn from trainers (academy, combat masters), improve via combat use
- **Equipment + Skills progression system** - no levels, growth via gear drops and skill proficiency
- Combat-ready spell system for testing
- Skill-based advancement (no traditional levels)
- Test client passes

**Why These Two**: 
- **Scholar**: Unique, lore-grounded exploration-based progression. Must hunt items to grow, drives world engagement.
- **Warrior**: Combat-focused maneuver progression. Low barrier to entry (grab weapon, learn basics, fight). Can explore for advanced trainers optionally.

**Class Philosophy**:
- **No levels** - skill-based advancement only
- Each class has a unique progression flavor:
  - Scholar: Exploration-driven (Study items → learn spells → improve Study skill)
  - Warrior: Combat-driven (learn maneuvers → use in combat → improve maneuver skills)
- Different incentive structures: Scholar must explore to stay viable. Warrior can fight immediately.

**Deferred Classes**:
- **Bard** (stretch goal for late MVP if time allows—has unique flavor tied to performance/inspiration)
- Spirit-communing classes (Ranger, Druid variants) tied to Colistani animism—design these separately post-MVP
- Legacy placeholders (Chevalier, Lancer, Paladin, Monk, Rogue, Illusionist, Elementalist, Enchanter) are removed—these were pre-design templates, not lore-grounded

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

### Growth Point 6: Class Progression Philosophy

**Warrior Advancement**:
- **Combat Maneuvers** - Tactics learned from trainers (Vojvoda Zsa's academy, combat masters in towns)
- Learn basic techniques (Power Attack, Defensive Stance, Guard)
- **Lower barrier to entry**: grab a weapon, learn basic moves, fight monsters
- Maneuver proficiency improves through combat use
- Can optionally explore to find distant trainers for advanced techniques
- Simple, self-directed loop: fight → improve → fight better

**Scholar Advancement**:
- **Spell Discovery** - only way to learn new spells is Study
- Must actively hunt for magical items
- Can't just cast Arcane Bolt forever and plateau—exploration is built-in progression
- Proficiency improves through successful studying
- Engaging, exploration-driven loop: hunt → study → expand repertoire
- Higher cognitive load but more rewarding for curious players

**Design Principle**: Both skill-based. Different incentive structures:
- Warrior: fight monsters first, optionally explore for trainers
- Scholar: must explore to stay viable

### Growth Point 7: Class Expansion (After MVP)
- Performance-based inspiration (buff allies)
- Unique lore tie: wandering performers in Colistani culture
- Spell learning: study performance items instead of scrolls

**Spirit-Touched Classes** (Post-MVP, designed from animism):
Tie to Colistani animal spirits (Fox/Owl/Deer/Snake):
- Fox-touched: hunting/tracking bonuses
- Owl-touched: wisdom/utility (Scholar already owns learning)
- Deer-touched: blessing/warding
- Snake-touched: underworld/shadow magic

**Design Philosophy**: Each new class needs:
1. Lore grounding (tied to world, not generic archetype)
2. Unique mechanic (not just "different spell list")
3. Actual design document (not template placeholders)

**What We're NOT Doing**: Legacy classes like Chevalier/Lancer/Paladin/Monk/Rogue/Illusionist were pre-design templates. Don't ship them. Only add classes that fit NJATA's world.

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
- [ ] Spell definitions (8 spells, minimal JSON)
- [ ] Spell loader
- [ ] Cast command (very simple)
- [ ] **NEW: Study command (Scholar-specific)**
- [ ] **NEW: 8 magical item definitions (wands/scrolls with spell IDs)**
- [ ] Damage calculation (basic dice + attribute)
- [ ] Proficiency tracking & Study skill tracking
- [ ] Test: All 8 spells castable
- [ ] Test: Scholar can find, study, and learn from items

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

