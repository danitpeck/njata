# NJATA Implementation Strategy & Technical Decisions

## Overview

This document outlines the recommended implementation roadmap, technical decisions for the skills system, and answers to key design questions.

---

## Part 1: Key Design Decisions

### 1.1: Spell Learning Model

**DECISION: Auto-Learn + Optional Boost**

Players **automatically learn all class spells** at appropriate levels:
```
Level 1:  Magic Missile, Spark, Weak Heal, Mage Armor (if available for class)
Level 5:  Fireball, Frost Bolt, Heal, Lightning Strike, etc.
Level 10: Teleport, Identify, etc.
Level 15: Gate, Summon, etc.
Level 20: Resurrection, Wish, etc.
```

**Optional Paid Learning**: Players can hire NPCs to teach spells early:
```
Cost: 100 + (spell_level * 50) gold
Time: Instant
Requirement: Player level >= spell level
```

**Rationale**: 
- Reduces friction - players always have combat spells
- Rewards level progression
- Still offers cosmetic choice (hire early for advantages)
- Matches legacy SMAUG system design

### 1.2: Proficiency Starting Point

**DECISION: Start at 30% Proficiency**

New spells learned at: **30% proficiency**

```
Example: Vex learns Fireball at level 5
- Immediate proficiency: 30%
- Effective damage: 30% penalty
- Formula: "4d8 + I" → becomes "4d8 * 0.3 + I * 0.3" effective? 
  NO - we use: "4d8 + I + (proficiency / 20)" → scales naturally
```

**Proficiency Gains Per Cast**:
```
Successful Cast: +1.5% proficiency
Resisted Cast: +0.5% proficiency
Failed Cast: 0% (no learning)
```

**Time to Mastery**:
```
Starting: 30%
Reaching 100%: 70% / 1.5% per cast = ~47 successful casts
Typical play: 2-4 weeks per spell = 70-100 days to max all spells
```

**Rationale**:
- Encourages spell use (gain proficiency)
- New players not overpowered
- Incentivizes progression (better spells as proficiency grows)
- Matches legacy learned[] percentage system

### 1.3: Component Strictness

**DECISION: Components Only for "Special" Spells**

**Basic Spells** (Combat, Healing, Buffs): No components required
```
1001 Magic Missile - Just mana
1005 Fireball - Just mana
1007 Heal - Just mana
(Most common spells: mana-only)
```

**Special Spells** (Resource-intensive, world-altering):
```
1018 Gate - Mana + 50 gold (expensive transport)
1019 Summon Familiar - Mana only (components optional)
1024 Animate Dead - Requires corpse in room
1026 Wish - Admin control (no component system needed)
```

**Rationale**:
- Keeps combat fluid (no item farming mid-fight)
- Makes rare spells feel special
- Allows future quest items with components
- Reduces code complexity while keeping depth

### 1.4: Spell Scaling Model

**DECISION: Three-Factor Scaling**

```
Damage = BaseFormula * ProficiencyMultiplier + CasterBonus

Example: Fireball
BaseFormula: "4d8" = 18
CasterBonus: INT bonus (5) + Level scaling (1 per level)
CasterBonus at Level 10, INT 18: 5 + 10 = 15

Proficiency:
- 30%: 18 * 1.3 + 15 = 38 damage
- 50%: 18 * 1.5 + 15 = 42 damage
- 100%: 18 * 2.0 + 15 = 51 damage

Saves:
- Half damage on successful Reflex save
- So: 38 damage → 19 with save
```

**Formula Complexity**:
```
Simple: "1d6 + I"              (Magic Missile - reliable)
Medium: "4d8 + I + l/2"        (Fireball - scales with level)
Complex: "6d10 + I*1.5 + l"    (Earthquake - heavy caster reliance)

Proficiency always multiplies final result by (1 + prof/100)
```

**Rationale**:
- Stat bonuses make attributes matter
- Level progression feels rewarding
- Proficiency creates long-term goals
- Predictable scaling (not exponential)

### 1.5: Area of Effect Friendliness

**DECISION: Strict Hostile AoE (No Friendly Fire)**

```
Fireball cast in room with alliances:
├─ Caster: Unharmed
├─ Caster's party: Unharmed
├─ Enemy 1: Affected
├─ Enemy 2: Affected
└─ Neutral NPC: SAFE (not affected)
```

**Targeting Logic**:
```go
type AoETargeting struct {
    Mode string        // "hostile_area", "allied_area", "universal_area"
    Radius int         // squares affected
    AffectedGroups []string  // ["hostile"] or ["allied"] or ["all"]
}

// Hostile area: Only attack enemies (faction-based)
// Allied area: Only help allies (party-based)
// Universal: Affect ALL (Earthquake - affects everyone)
```

**Rationale**:
- No punishment for using AoE
- Encourages tactical positioning
- Prevents griefing in parties
- Matches modern MMORPG expectations
- Legacy SMAUG had limitations due to codebase

### 1.6: Cooldown Sharing

**DECISION: No Shared Cooldown Pool (Individual Per-Spell)**

```
Mage casts Fireball
├─ Fireball enters 5s cooldown
├─ Frost Bolt: READY (different cooldown)
├─ Lightning Strike: READY
├─ Magic Missile: READY
└─ After 5s: Fireball ready again
```

**Why not shared cooldown?**
```
Legacy SMAUG example: "spell_cast_wait = 1" = wait state
Modern approach: Per-spell cooldown
```

**Future Enhancement**: Cooldown categories (optional)
```
If implemented later:
├─ "fire_spells" = Fireball + Fire Breath (3s shared)
├─ "ice_spells" = Frost Bolt + Chill Touch (2s shared)
└─ "healing" = Heal + Cure Light (no pool, individual only)
```

**Rationale**:
- Simpler to implement
- Allows more action/rotation flexibility
- Can add pools later if balance requires
- Prevents "lockout" feeling

### 1.7: Message Randomization

**DECISION: Template-Based with Subtle Variance**

```
Option A (Simple): Same message every cast
"You hurl a fireball at Goblin."

Option B (Variance): Permutation templates
"You [hurl|throw|blast] a fireball [at|toward|upon] $target."
→ Hurl/throw/blast (3 options) × at/toward/upon (3 options) = 9 combos

Option C (Dynamic): Adjective insertion
"You hurl a [massive|blazing|scorching] fireball at $target."
```

**DECISION**: Implement **Option B - Permutation Templates**

```go
type MessageTemplate struct {
    Verb    []string   // ["hurl", "throw", "blast"]
    Prep    []string   // ["at", "toward", "upon"]
    Object  string     // "$target"
}

Example Message:
"You throw a fireball toward the Goblin."  // Different each cast
```

**Rationale**:
- Adds replay value (messages feel fresh)
- EZ to implement
- Not overwhelming (still recognizable)
- Matches player expectations from modern games

---

## Part 2: Implementation Roadmap

### Milestone 1: Enhanced Spell System (Week 1-2)

**Goal**: Expand from 6 spells to 15+ with full mechanics

**Tasks**:
- [ ] Expand `skills/spells.json` to 15 spells (level 1-15)
- [ ] Add component system to spell definition
- [ ] Add message templates to spell definition
- [ ] Add proficiency scaling to spell definition
- [ ] Update spell loader to parse all new fields
- [ ] Create spell catalog help text

**Deliverable**: 
- Playable spells 1001-1015
- All spell definitions complete
- Test: `spellbook` shows all spells with correct proficiency

### Milestone 2: Damage System (Week 2-3)

**Goal**: Implement damage calculations and apply actual damage

**Tasks**:
- [ ] Create formula parser for damage/healing expressions
- [ ] Implement variable substitution (I, l, proficiency, etc.)
- [ ] Create damage effect application
- [ ] Implement saving throw system (Reflex/Will)
- [ ] Add resistance/immunity checking
- [ ] Test proficiency scaling (30% → 100%)

**Deliverable**:
- Spells do real damage (currently they just message)
- Test: `cast fireball goblin` → damage applied, HP reduced
- Test: `cast heal self` → healing applied, HP increased

### Milestone 3: Combat Integration (Week 3-4)

**Goal**: Spells work in actual combat with enemies

**Tasks**:
- [ ] Create `fight` command
- [ ] Add mob HP/damage tracking
- [ ] Implement enemy attacking back
- [ ] Add spell targeting during combat
- [ ] Integrate spell damage into combat flow
- [ ] Test: Kill mob with spells, gain experience

**Deliverable**:
- Playable combat loop with spells
- Test: Fight goblin, cast spells, defeat enemy, gain XP

### Milestone 4: Effects/Affects System (Week 4-5)

**Goal**: Buffs/debuffs that persist and affect gameplay

**Tasks**:
- [ ] Create affect/buff tracking in Player
- [ ] Implement duration tracking
- [ ] Create stat modification application
- [ ] Implement affect message display
- [ ] Add affecting commands (show active buffs)
- [ ] Test: Apply blindness, verify AC penalty

**Deliverable**:
- Potion/buff effects work
- Test: Cast Blindness on enemy, verify AC penalty applied

### Milestone 5: Advanced Spells (Week 5-6)

**Goal**: Area spells, summoning, special effects

**Tasks**:
- [ ] Implement area targeting (Earthquake)
- [ ] Implement multi-target processing
- [ ] Create summon system (Familiar)
- [ ] Add channeled casting (Resurrection)
- [ ] Implement teleportation
- [ ] Test: Earthquake damages all enemies in area

**Deliverable**:
- All 15-20 spells functional
- Test: Complex spell scenarios work

### Milestone 6: Proficiency System (Week 6-7)

**Goal**: Spell learning, progression, mastery

**Tasks**:
- [ ] Implement proficiency-based learning
- [ ] Create spell teaching NPCs
- [ ] Add proficiency gain on cast
- [ ] Persist proficiency to player file
- [ ] Create proficiency display
- [ ] Test: Cast spell 100x, proficiency increases

**Deliverable**:
- Full learning arc
- Test: New spell starts 30%, gains 1-2% per cast, reaches 100%

### Milestone 7: Polish & Balance (Week 7-8)

**Goal**: Game feel, feedback, balance

**Tasks**:
- [ ] Add visual/text feedback (color, formatting)
- [ ] Implement spell animations (text-based)
- [ ] Balance damage numbers
- [ ] Test against various enemy levels
- [ ] Create spell help documentation
- [ ] Implement spell descriptions

**Deliverable**:
- Game feels good to play
- Spells feel balanced
- All 28 spells documented

---

## Part 3: Technical Implementation Details

### 3.1: Enhanced Spell Definition Schema

```json
{
  "id": 1005,
  "name": "Fireball",
  "category": "offensive",
  "description": "Hurl a massive ball of flame at your enemies",
  "school": "evocation",
  "level_required": 5,
  "mana_cost": 30,
  "cooldown_seconds": 5,
  "casting_time": 0,
  
  "components": {
    "required": [
      {
        "type": "mana",
        "amount": 30,
        "consumed": true
      }
    ],
    "restricted": [
      {
        "type": "status",
        "status": "silenced",
        "reason": "Cannot cast while silenced"
      }
    ]
  },
  
  "targeting": {
    "mode": "hostile_area",
    "max_range": 30,
    "area_radius": 3,
    "requires_line_of_sight": true
  },
  
  "effects": {
    "damage": {
      "type": "fire",
      "formula": "4d8 + I + (proficiency / 20)",
      "scales_with": ["intelligence", "proficiency"]
    },
    "save": {
      "type": "reflex",
      "dc": 14,
      "on_success": "half_damage"
    },
    "resistance": [
      {
        "type": "fire",
        "reduces_damage_by": 25
      }
    ]
  },
  
  "messages": {
    "cast": {
      "verbs": ["hurl", "throw", "blast"],
      "preps": ["at", "toward", "upon"],
      "template": "You {verb} a fireball {prep} $target."
    },
    "hit_target": "$actor's fireball scorches you for $damage damage!",
    "hit_room": "$target is engulfed in $actor's flames!",
    "miss_target": "$actor's fireball misses you!",
    "miss_room": "$actor's fireball sails past $target.",
    "immune": "$target is completely immune to fire!"
  },
  
  "logistics": {
    "teaching_cost_gold": 500,
    "proficiency_cap": 100,
    "starting_proficiency": 30,
    "proficiency_gain": 1.5
  }
}
```

### 3.2: Go Implementation Structure

```
internal/spells/
├── spell.go
│   ├── Spell struct
│   ├── SpellEffect interface
│   └── LoadSpells()
│
├── effects.go
│   ├── DamageEffect
│   ├── HealingEffect
│   ├── AffectEffect
│   └── Apply() methods
│
├── targeting.go
│   ├── TargetMode enum
│   ├── FindTargets()
│   └── IsValidTarget()
│
├── formula.go
│   ├── FormulaParser
│   ├── EvaluateFormula()
│   └── VariableSubstitution
│
├── effects/
│   ├── damage.go - ApplyDamage()
│   ├── healing.go - ApplyHealing()
│   ├── buff.go - ApplyBuff()
│   └── debuff.go - ApplyDebuff()
│
└── casting/
    ├── cast.go - CastSpell()
    ├── validation.go - ValidateCast()
    ├── cooldown.go - ManageCooldowns()
    └── proficiency.go - UpdateProficiency()
```

### 3.3: Casting Flow (Pseudocode)

```go
func CastSpell(caster *Player, spellID int, targetArg string) {
    // 1. Load spell
    spell := GetSpell(spellID)
    
    // 2. Validate casting
    if !caster.Learned[spell.ID] {
        return "You don't know that spell."
    }
    if caster.Mana < spell.ManaCost {
        return "Insufficient mana."
    }
    if !IsReadyToCast(caster, spell.ID) {
        return fmt.Sprintf("This spell is not ready yet. [%.1fs remaining]", 
            GetCooldownRemaining(caster, spell.ID))
    }
    
    // 3. Find target
    targets := FindTargets(caster, spell, targetArg)
    if len(targets) == 0 {
        return "Target not found."
    }
    
    // 4. Apply spell
    for _, target := range targets {
        // Calculate damage
        baseDamage := RollDice(spell.Effects.Damage.Formula)
        profMult := 1.0 + (float32(caster.Proficiency[spell.ID]) / 100.0)
        damage := int(float32(baseDamage) * profMult)
        
        // Saving throw
        if target.MakeSave(spell.Effects.Save) {
            damage = damage / 2
            SendMessage(target, "You feel the effects reduced!")
        }
        
        // Resistances
        resistance := target.GetResistance(spell.Effects.Damage.Type)
        damage = damage * (100 - resistance) / 100
        
        // Apply damage
        target.TakeDamage(damage)
        
        // Messages
        SendMessage(caster, spell.Messages.Cast)
        SendMessage(target, spell.Messages.HitTarget)
        SendRoomMessage(caster, spell.Messages.HitRoom)
    }
    
    // 5. Cleanup
    caster.Mana -= spell.ManaCost
    SetCooldown(caster, spell.ID)
    UpdateProficiency(caster, spell.ID)
    
    return "Spell cast successfully!"
}
```

### 3.4: Proficiency Persistence (JSON)

```json
{
  "name": "vex",
  "race": 6,
  "class": 12,
  "skills": {
    "1001": { "proficiency": 45, "learned": true, "lifetime_casts": 127 },
    "1005": { "proficiency": 78, "learned": true, "lifetime_casts": 256 },
    "1007": { "proficiency": 30, "learned": true, "lifetime_casts": 3 },
    "1008": { "proficiency": 0, "learned": false, "lifetime_casts": 0 }
  }
}
```

---

## Part 4: Balance Spreadsheet

### Mana Per Damage Ratio (Efficiency)

```
Spell               Mana  DamageRange  Avg    Efficiency  Notes
Magic Missile       15    1-6 + mods   8      0.53        Spam spell
Spark               8     1-4 + mods   4      0.50        Very spammable
Fireball            30    4-32 + mods  32     1.07        Core spell
Frost Bolt          25    3-24 + mods  28     1.12        Control element
Lightning Strike    40    5-50 + mods  45     1.13        High damage high cost
Earthquake          80    6-72 + mods  60     0.75        Area efficiency
Meteor Storm        100   8-80 + mods  85     0.85        Extreme damage area
```

**Goal**: Mana efficiency between 0.5-1.2 (fair trade-off)

### Spell Difficulty vs Reward

```
Level  Spell                Difficulty  Reward    Balance
1      Magic Missile        Easy        Low       ✓ Learning spell
5      Fireball             Medium      Medium    ✓ Key spell
10     Teleport             Hard        High      ✓ Mobility
15     Earthquake           Hard        High      ✓ Power spike
20     Resurrection         Hard        High      ✓ Rare ability
```

---

## Part 5: Testing Strategy

### Unit Tests
- Formula parser (dice, variables)
- Proficiency calculations
- Cooldown tracking
- Message template generation

### Integration Tests
- Full spell cast flow
- Damage application
- Effect persistence
- Proficiency gains

### Balance Tests
- Damage vs mana efficiency
- Proficiency progression speed
- Cooldown fairness
- Effect duration balance

### User Acceptance Tests
- New player can learn to cast in <5 min
- Spells feel powerful and useful
- Progression feels rewarding
- No game-breaking exploits

---

## Part 6: Future Enhancement Ideas

### Level 2: Advanced Magic
- Spell interruption system
- Casting failures
- Mana shield mechanic
- Spell mastery tier (100%+ proficiency bonuses)

### Level 3: Complex Systems
- Spell combinations (synergies)
- Magic schools with bonus mechanics
- Spell attunement (pick 5 of 20 spells)
- Spell mutation (randomized effects)

### Level 4: Social Systems
- Spell trading between players
- Spell creation system
- Community spell voting
- Ranked spell tournaments

---

**END OF IMPLEMENTATION STRATEGY**

**NEXT RECOMMENDED STEP**: Review design decisions 1.1-1.7, discuss any changes, then proceed with Milestone 1 implementation.
