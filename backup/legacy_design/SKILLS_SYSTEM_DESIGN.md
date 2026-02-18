# NJATA Skills & Spells System Design

## Executive Summary

The legacy SMAUG system offers decades of proven mechanics for skills and spellcasting. This design proposes a modernized, robust system that maintains proven gameplay mechanics while improving code clarity and extensibility.

**Key Design Principles:**
- **Learnable**: Players progressively learn skills; not all are available at start
- **Costly**: Spells require mana, components, or other resources
- **Failurable**: Casting can fail based on various conditions (mana, components, resistance, levels)
- **Scalable**: Damage/healing scales with caster level, stats, and proficiency
- **Meaningful**: Every spell type has distinct mechanical and narrative impacts
- **Persistent**: Spell effects last (buffs/debuffs, enchantments, damage)

---

## 1. Core Skill/Spell Taxonomy

### Skill Types (from legacy: skills.c line ~100)

```
SKILL_SPELL      - Magical spells (Fireball, Heal, Gate, etc.)
SKILL_SKILL      - Trained abilities (Shield Block, Parry, etc.)
SKILL_WEAPON     - Weapon proficiencies (Sword, Axe, Bow, etc.)
SKILL_RACIAL     - Racial abilities (Elven sight, Dwarven throw, etc.)
SKILL_TONGUE     - Languages (Common, Elvish, Draconic, etc.)
SKILL_CRAFT      - Crafting abilities (Blacksmithing, Alchemy, etc.)
```

### Spell Categories

**Offensive** - Direct damage to target
- Example: Fireball, Magic Missile, Lightning Bolt
- Mechanics: Damage formula, save vs damage, targeting hostile only

**Defensive** - Protect caster/allies
- Example: Shield Spell, Stoneskin, Mirror Image
- Mechanics: AC/damage reduction, duration effects, stacking rules

**Healing** - Restore HP/vitals
- Example: Heal, Cure Light Wounds, Revive
- Mechanics: Touch range or target, scales with caster level/WIS

**Control** - Affect movement/actions
- Example: Gate, Teleport, Mist Walk, Plant Pass
- Mechanics: Position changes, area restrictions, travel costs

**Utility** - Information gathering
- Example: Identify, Farsight, Locate Object, Detect Poison
- Mechanics: Non-combat, information reveal, cost efficiency

**Debuff** - Harm enemy state (not direct damage)
- Example: Curse, Poison, Blindness, Weaken, Sleep
- Mechanics: Saving throws, duration, stacking restrictions

**Buff** - Enhance ally state
- Example: Strength, Invisibility, Invis, Bless
- Mechanics: Self or target, duration, non-stacking rules

---

## 2. Complete Spell Definition Schema

```json
{
  "id": 1001,
  "name": "Fireball",
  "category": "offensive",
  "type": "spell",
  "level_required": 5,
  "mana_cost": 30,
  "cooldown_seconds": 5,
  
  "applicability": {
    "classes": [1, 2, 5],     // Mage, Enchanter, Augurer
    "races": [],               // All races can learn
    "restricted_for": [6]      // Can't learn: Warrior
  },
  
  "components": {
    "required": [
      {
        "type": "mana",
        "amount": 30,
        "consumed": true
      },
      {
        "type": "item_type",
        "item_type": 12,       // e.g., "gem"
        "amount": 1,
        "consumed": true,
        "optional": false
      }
    ],
    "restricted": [
      {
        "type": "affliction",
        "affliction": "silenced",
        "description": "Cannot cast while silenced"
      }
    ]
  },
  
  "targeting": {
    "mode": "hostile_single",  // self, ally_single, hostile_single, area
    "max_range": 25,
    "requires_line_of_sight": true,
    "affects_caster": false
  },
  
  "success_rate": {
    "base": 100,
    "scales_with": ["proficiency", "intelligence"],
    "resisted_by": ["magic_resistance"],
    "save_type": "reflex",
    "save_difficulty": 14
  },
  
  "effects": {
    "damage": {
      "type": "fire",
      "formula": "4d8 + int_bonus + (proficiency * 0.5)",
      "modifiers": {
        "per_level": 0.25,
        "per_intelligence": 0.3,
        "per_proficiency": 0.5
      }
    },
    "area_of_effect": {
      "radius": 3,
      "affected_groups": "hostile"
    }
  },
  
  "cooldown": {
    "seconds": 5,
    "affects": ["cast_any", "cast_fireball"],
    "shared_with": ["frost_bolt", "fire_breath"]
  },
  
  "messaging": {
    "cast_by_caster": "You hurl a massive fireball at $target.",
    "cast_by_observer": "$actor hurls a massive fireball at $target.",
    "hit_target": "$actor's fireball engulfs you in flames!",
    "hit_observer": "$target is engulfed in $actor's flames!",
    "miss_target": "$actor's fireball misses you!",
    "miss_observer": "$actor's fireball sails past $target.",
    "immune_to_caster": "$target is immune to fire.",
    "immune_to_observer": "$target is immune to fire."
  },
  
  "logistics": {
    "teaching_cost_gold": 100,
    "learning_cost_experience": 500,
    "proficiency_cap": 100,
    "proficiency_gain_per_cast": 1
  }
}
```

---

## 3. Spell Components System (legacy: magic.c ~1043)

**Purpose**: Spells can require items, gold, health, or be restricted by conditions

### Component Types

| Type | Description | Example |
|------|-------------|---------|
| `mana` | Mana cost (primary) | 30 mana for Fireball |
| `gold` | Gold payment | 50 gold for Gate spell |
| `health` | HP cost | 10 HP for Drain Touch |
| `item_type` | Specific item type | Gem (type 12) for Enchant |
| `item_vnum` | Specific item | Scroll of Fireball (vnum 5041) |
| `item_keyword` | Item with keyword | Anything with "crystal" keyword |
| `health_requirement` | Must have min HP | Must have >20 HP remaining |
| `movement_requirement` | Must have movement | Must have >10 movement remaining |

### Component Operators (legacy: magic.c ~1030-1043)

```
!  - Spell fails if player HAS this component
+  - Component is NOT consumed
@  - Decrease item value[0], extract at 0
#  - Decrease item value[1], extract at 0
$  - Decrease item value[2], extract at 0
%  - Decrease item value[3], extract at 0
^  - Decrease item value[4], extract at 0
&  - Decrease item value[5], extract at 0
```

**Example**: 
- `+G500` = Check for 500 gold, don't consume
- `@V3052` = Decrease value[0] of vnum 3052, extract when 0
- `!T5` = Spell fails if player has item of type 5

---

## 4. Damage & Healing Formulas

**Legacy System** (magic.c ~640-900): Uses `rd_parse()` for dice expression parsing with order of operations

### Damage Formula Variables

```
l = caster's level
L = victim's level
H = caster's HP
M = caster's Mana
V = caster's Movement
S/I/W/C/D/A/U = caster's attributes (STR/INT/WIS/CON/DEX/CHA/LCK)
```

### Example Formulas

```
Fireball:      "4d8 + I + (proficiency / 20)"       // 4-32 + INT + prof bonus
Magic Missile: "1d6 + I/2 + l"                      // Scales with level
Heal:          "2d8 + 5 + (W * 1.5)"               // Wisdom-based healing
Lightning:     "5d10 + I + (l * 0.5)"              // Level and INT scaling
Weak Heal:     "1d4 + 3"                           // Basic healing
```

### Proficiency Scaling

```
Damage/Healing = BaseFormula * (1 + (Proficiency / 100))

Example:
- Fireball at 50% proficiency: 4d8 base * 1.5 = 1.5x damage
- Fireball at 100% proficiency: 4d8 base * 2.0 = 2.0x damage
```

### Resistance Adjustment (legacy: magic.c ~635-641)

```go
// Damage reduced by resistance/immunity
if IsImmune(target, damageType) {
    damage = 0
    return "They are immune!"
}

if IsResistant(target, damageType) {
    damage = damage / 2
}

if IsSusceptible(target, damageType) {
    damage = damage * 1.5
}
```

---

## 5. Saving Throws System (legacy: magic.c ~957-976)

**Three Save Types**: Fortitude, Reflex, Will

```go
func SavingThrow(casterLevel int, victim *Character, saveType string) bool {
    // Base 50% chance, adjusted by level difference
    var baseSave int
    
    switch saveType {
    case "fortitude":
        baseSave = 50 + (victim.Level - casterLevel + victim.SaveFortitude) * 5
    case "reflex":
        baseSave = 50 + (victim.Level - casterLevel + victim.SaveReflex) * 5
    case "will":
        baseSave = 50 + (victim.Level - casterLevel + victim.SaveWill) * 5
    }
    
    // Clamp between 5-95%
    baseSave = Max(5, Min(95, baseSave))
    
    // Roll d100 vs base save
    return Rand(1, 100) <= baseSave
}
```

**Success Outcomes**:
- **Successful Save**: Damage reduced by 50%, effect negated, or duration halved
- **Failed Save**: Full damage/effect, full duration

---

## 6. Targeting System

### Targeting Modes

```
self              - Caster only (Invisibility, Strength, Heal Self)
ally_single       - One allied character (Heal)
ally_area         - All allies in room (Group Heal)
hostile_single    - One enemy (Fireball)
hostile_area      - All enemies in area (Earthquake)
object            - Find object by name (Locate Object)
any_character     - Any character (Identify)
location          - Ground location (Gate, Teleport)
```

### Range System

```go
type TargetInfo struct {
    Mode string
    MaxRange int          // 0 = touch, 25+ = ranged
    RequiresLineOfSight bool
    AreaRadius int        // For area spells (0 = single target)
}
```

---

## 7. Spell Effects & Affects System

### Damage Effects

```go
type DamageEffect struct {
    Type string                // fire, cold, acid, lightning, slashing, etc.
    Amount int                 // Direct damage
    Formula string             // "4d8 + I + (proficiency/20)"
    ScalesWithLevel bool
    ScalesWithStat string      // intelligence, wisdom, strength, etc.
}
```

### Affect/Buff Effects

```go
type AffectEffect struct {
    Name string                // "strength_boost", "blindness", etc.
    DurationSeconds int
    DurationFormula string     // "30 + (W * 5)" -> scales with wisdom
    
    // Stat modifications
    AttributeMods map[string]int  // {"strength": +2, "dexterity": -3}
    
    // Combat modifications
    ACBonus int                // -5 = 5 points better AC
    DamageBonus int
    DamageResistance map[string]int  // {"fire": 25, "cold": 10}
    
    // State modifications
    CanAct bool
    CanMove bool
    VisibleState string        // "invisible", "ethereal", etc.
    
    // Messages
    ApplyMessage string        // "You feel stronger!"
    RemoveMessage string       // "You feel weaker."
    
    // Stacking rules
    MaxStacks int              // 1 = exclusive buff, 3 = stackable
    StackType string           // "strength" - same type don't stack
}
```

### Duration Formulas

```
Strength Buff:     "30 + (W * 5)"      -> 30 seconds + 5 per WIS
Invisibility:      "120"               -> Fixed 2 minutes
Blindness:         "30 + (L * 2)"      -> 30 + 2 per target level
Poison:            "l * 10"            -> Duration based on caster level
```

---

## 8. Proficiency & Learning System

### Player Skill Progress

```json
{
  "player_id": "vex",
  "skills": {
    "1001": {
      "proficiency": 45,
      "learned": true,
      "learned_at_level": 5,
      "learned_at_time": 1708151400,
      "lifetime_casts": 127,
      "last_cast": 1708151890
    }
  }
}
```

### Proficiency Progression

```
Category         Min Level  Starting   Per Cast  Max at Level 20
--------         ---------  ---------  --------  ----------------
Basic Spell      1          30%        0.5%      85%
Intermediate     5          20%        1.0%      95%
Advanced         10         10%        1.5%      100%
Mastery          15         5%         2.5%      100%
```

### Learning Mechanics

**Auto-Learn**: Player learns all class spells at appropriate levels
```
Mage learns Fireball at level 5
Cleric learns Heal at level 1
Ranger learns Track at level 3
```

**Manual Learn**: Hire NPC to teach spell
```
Cost: 100 + (spell_level * 50) gold
Time: Instant
Requirement: Character level >= spell level
```

**Natural Progression**: Spells gain proficiency through use
```
Each cast gives: 0.5-2.5% proficiency (varies by spell level)
Successful cast: +1.5x proficiency gain
Resisted cast: +0.5x proficiency gain
Failed cast: No proficiency gain
```

---

## 9. Cooldown & Cast Queue System

### Cooldown Types

```go
type Cooldown struct {
    SpellID int
    ExpiresAt int64        // Unix timestamp (nanoseconds)
    Category string        // "firemagic", "healing", "combat", etc.
    SharedWith []int       // Other spells on same cooldown
}
```

### Cooldown Examples

```
Magic Missile:   2 seconds  (spam capability)
Fireball:        5 seconds  (moderate cooldown)
Teleport Gate:   30 seconds (long cooldown)
Revive:          120 seconds (rare powerful spell)

Shared Cooldown: All fire spells share 3-second cooldown
```

### Cast Queue

```go
// High-level casters can queue next spell during cast time
type CastState struct {
    CurrentCast SpellID
    CastTime int            // milliseconds
    QueuedSpell SpellID
    QueuedTarget Target
}
```

---

## 10. Failure Conditions

**Spell fails if:**

1. **Mana**: Insufficient mana for cast
   ```
   Required: 30 mana for Fireball
   Available: 15 mana
   → "You don't have enough mana."
   ```

2. **Components**: Missing required items/gold
   ```
   Gem required: Not in inventory
   → "You lack the necessary components."
   ```

3. **Status**: Afflictions prevent casting
   ```
   Silenced status active
   → "You cannot speak the incantation!"
   ```

4. **Target**: Invalid target selection
   ```
   "cast fireball self" while needing hostile
   → "Fireball can only target enemies."
   ```

5. **Line of Sight**: Ranged spells blocked
   ```
   Hidden behind wall
   → "You cannot see that target."
   ```

6. **Save**: Target makes saving throw
   ```
   Lightning Bolt cast (Reflex save)
   Target rolls save successfully
   → Damage reduced by 50%
   ```

7. **Resistance**: Target is immune
   ```
   Fire Breath vs fire-immune dragon
   → "The dragon is immune to fire!"
   ```

8. **Cooldown**: Spell recently cast
   ```
   Fireball on 3.2 second cooldown (5s total)
   → "This spell is not ready yet. [3.2s remaining]"
   ```

---

## 11. Message System

### Message Contexts

```
TO_CASTER       - Only the spell caster sees this
TO_TARGET       - Only the target of the spell sees this
TO_ROOM         - Everyone in the room except caster/target
TO_ALL          - Everyone including caster/target
```

### Message Variables

```
$actor     - Name of caster
$target    - Name of spell target
$spell     - Name of spell
$damage    - Amount of damage dealt
$healing   - Amount of healing done
$object    - Name of object affected
$class     - Target's class
$his/her   - Pronoun (his/her)
```

### Example Messages

```
Cast by Caster:
  "You hurl a massive fireball at $target."

Damage to Target:
  "$actor's fireball scorches you for $damage damage!"
  
Damage to Observers:
  "$actor hurls a fireball at $target, dealing $damage damage!"

Miss Message:
  "$actor's fireball sails past $target harmlessly."

Resistance:
  "$target's resistance to fire reduces the damage to $damage."

Immunity:
  "$target is completely immune to fire magic!"
```

---

## 12. Data Structures in Go

### Spell Definition

```go
type Spell struct {
    ID              int
    Name            string
    Category        string    // offensive, defensive, healing, etc.
    LevelRequired   int
    ManaCost        int
    CooldownSeconds int
    
    Applicability   SpellApplicability
    Components      []Component
    Targeting       TargetInfo
    SuccessRate     SuccessInfo
    Effects         SpellEffects
    Cooldown        CooldownInfo
    Messaging       MessageSet
    Logistics       LogisticsInfo
}

type Component struct {
    Type        string      // mana, gold, item_type, item_vnum, etc.
    Amount      int
    Consumed    bool
    Optional    bool
}

type TargetInfo struct {
    Mode                string
    MaxRange            int
    RequiresLineOfSight bool
    AreaRadius          int
}

type SuccessInfo struct {
    BaseRate        int
    ScalesWith      []string   // proficiency, intelligence, etc.
    ResistedBy      []string   // magic_resistance, etc.
    SaveType        string     // fortitude, reflex, will
    SaveDifficulty  int
}

type SpellEffects struct {
    Damage          DamageEffect
    Healing         HealingEffect
    Affect          AffectEffect
    AreaOfEffect    AreaInfo
}
```

### Player Spell Progress

```go
type PlayerSkillProgress struct {
    SkillID        int
    Proficiency    int       // 0-100
    Learned        bool
    LearnedAt      time.Time
    LifetimeCasts  int
    LastCast       time.Time
    CurrentCooldown int64    // Unix nanoseconds
}
```

---

## 13. Implementation Roadmap

### Phase 1: Foundation ✅ (Current)
- [x] Basic spell definitions (6 spells)
- [x] Mana costs and cooldowns
- [x] spellbook/cast commands
- [ ] **Add**: Components system loading
- [ ] **Add**: Damage/healing formulas

### Phase 2: Combat Integration (Recommended Next)
- [ ] Fight command with spell damage
- [ ] Targeting and area effects
- [ ] Saving throws implementation
- [ ] Damage type resistances
- [ ] Enemy spell effects (affects/buffs)

### Phase 3: Effects System
- [ ] Buff/debuff stacking
- [ ] Duration tracking
- [ ] Status effect display
- [ ] Affect removal/dispel

### Phase 4: Advanced Features
- [ ] Multi-target spells
- [ ] Cast sequencing/queue
- [ ] Proficiency-based scaling
- [ ] Component consumption
- [ ] Spell teaching system

### Phase 5: Polish
- [ ] Message variability
- [ ] Emote system integration
- [ ] Spell animations
- [ ] Sound effects hints
- [ ] Spell descriptions in help system

---

## 14. Comparison: Legacy vs Modern Design

| Feature | SMAUG Legacy | NJATA Modern | Benefit |
|---------|-------------|-------------|---------|
| Spell Defs | Code array | JSON files | Hot reload, editor-friendly |
| Components | Cryptic codes | Structured types | Maintainable, clear intent |
| Formulas | String parsing | Type-safe eval | Faster, safer, extensible |
| Affects | Global linked list | Per-player map | Cleaner persistence |
| Cooldowns | Wait state bits | Timestamp map | More flexible, precise |
| Messaging | Hard-coded | Template vars | Reusable, consistent |
| Persistence | Binary structs | JSON | Human-readable, portable |
| Tests | Manual @test mobs | Automated suite | Regression-free |

---

## 15. Next Design Decisions

**For your review:**

1. **Proficiency System**: Start at 30% for new spells or 0%?
2. **Component Strictness**: Require components for all spells or only "resource spells"?
3. **Spell Scaling**: Per-proficiency step (1% per 1% prof) or threshold-based (30/60/90%)?
4. **Area Spells**: Should Fireball damage friendlies or hit-friendlies-safe?
5. **Cooldown Sharing**: Should wizard's fire spells share a cooldown pool?
6. **Message Customization**: Random variance per cast ("The fireball is powerful" vs "The fireball is VERY powerful")?

---

**END OF DESIGN DOCUMENT**
