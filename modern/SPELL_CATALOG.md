# NJATA Comprehensive Spell Catalog

This document defines the complete spell list, organized by category and level, drawing from the legacy SMAUG system and adapted for modern gameplay.

---

## Spell Organization

- **Levels**: 1-20 (progression through game)
- **Schools**: Magic (8 schools of magic)
- **Categories**: Offensive, Defensive, Healing, Control, Utility, Debuff, Buff
- **Casting Time**: Instant (0ms) or channeled (100-5000ms)
- **Range**: Touch (1 square) or Ranged (5-100+ squares)

---

# TIER 1: NOVICE (Levels 1-4)

## Combat Foundation

### 1001: Magic Missile
```
Category: Offensive (Evocation)
Mana Cost: 15
Cooldown: 2 seconds
Level Required: 1
Range: 25 squares
Damage: "1d6 + (I/2) + l"     // 1-6 + INT/2 + caster level
Save: None (unerring)
Proficiency Gain: +0.5% per cast
```
**Purpose**: First combat spell, no save, always hits. Spam-able for learning magic fundamentals.

### 1002: Spark
```
Category: Offensive (Evocation)
Mana Cost: 8
Cooldown: 1 second
Level Required: 1
Range: 20 squares
Damage: "1d4 + (I/3)"         // 1-4 + INT/3
Save: None
Proficiency Gain: +0.3% per cast
```
**Purpose**: Minimal resource consumption, fastest cooldown for learning combat rotation.

### 1003: Weak Heal
```
Category: Healing (Abjuration)
Mana Cost: 12
Cooldown: 3 seconds
Level Required: 1
Range: Touch
Healing: "1d4 + 3 + (W/3)"    // 1-4 + 3 + WIS/3
Target: Self or ally (single target)
Proficiency Gain: +0.5% per cast
```
**Purpose**: Early healing for survivability, low cost, supports teamwork.

### 1004: Mage Armor
```
Category: Buff (Abjuration)
Mana Cost: 10
Cooldown: 10 seconds
Duration: 60 seconds + (W * 5)
Level Required: 1
Expected AC Improvement: 2 points
Effect: Provides magical armor
Message: "You feel shielded by a magical aura."
```
**Purpose**: Passive defense, teaches buff mechanics.

---

## TIER 2: APPRENTICE (Levels 5-9)

### 1005: Fireball
```
Category: Offensive (Evocation)
Mana Cost: 30
Cooldown: 5 seconds
Level Required: 5
Range: 30 squares
Damage: "4d8 + I + (proficiency / 20)"
Save: Reflex DC 14 (half damage on success)
Area: 3 square radius (all enemies in area)
Components: 
  - Mana (30, consumed)
Proficiency Gain: +1.0% per cast (successful), +0.3% (if saved)
```
**Purpose**: First area spell, teaches save mechanics and multi-target tactics.

### 1006: Frost Bolt
```
Category: Offensive (Evocation)
Mana Cost: 25
Cooldown: 4 seconds
Level Required: 5
Range: 28 squares
Damage: "3d8 + I + (proficiency / 25)"
Save: Reflex DC 13
Special: Slows target 1d3 rounds (move -25%)
Proficiency Gain: +1.0% per cast
```
**Purpose**: Single-target control spell with crowd control component.

### 1007: Heal
```
Category: Healing (Abjuration)
Mana Cost: 20
Cooldown: 6 seconds
Level Required: 5
Range: Touch
Healing: "2d8 + 5 + (W * 1.5)"  // Wisdom-based
Target: Self or ally
Proficiency Gain: +1.0% per cast
```
**Purpose**: Significant healing, core spell for support players.

### 1008: Lightning Strike
```
Category: Offensive (Evocation)
Mana Cost: 40
Cooldown: 8 seconds
Level Required: 5
Range: 35 squares
Damage: "5d10 + I + (l * 0.5)"   // Scales with level heavily
Save: Reflex DC 15
Special: 10% chance to stun target 1 round
Proficiency Gain: +1.5% per cast (expensive spell, faster learning)
```
**Purpose**: Highest damage in tier, teaches resource management (high cost).

### 1009: Invisibility
```
Category: Buff
Mana Cost: 20
Cooldown: 30 seconds
Duration: 120 + (W * 4)
Level Required: 5
Target: Self or one ally
Effect: Grants stealth until attacking
State: invisible = true
Message: "You fade from sight."
```
**Purpose**: Utility/escape, teaches duration-based effects.

### 1010: Blindness
```
Category: Debuff (Enchantment)
Mana Cost: 18
Cooldown: 4 seconds
Level Required: 5
Range: 25 squares
Duration: 20 + (L * 2) seconds
Target: Single enemy
Save: Will DC 13
Effects:
  - AC penalty: +4 (worse AC)
  - Attack penalty: -20%
Message: "$target stumbles blindly!"
```
**Purpose**: Debuff/control, teaches negative affects.

---

## TIER 3: JOURNEYMAN (Levels 10-14)

### 1011: Teleport
```
Category: Control
Mana Cost: 50
Cooldown: 30 seconds
Level Required: 10
Range: 50 squares (to known location)
Target: Self or one ally
Effect: Transport to target location
Message: "$actor vanishes in a flash of light!"
Save: Will DC 16 (unwilling target resists)
```
**Purpose**: Escape/mobility, expensive resource.

### 1012: Identify
```
Category: Utility (Divination)
Mana Cost: 15
Cooldown: 2 seconds
Level Required: 10
Range: Touch
Target: Object or character
Effect: Reveal full description, stats, enchantments
Message: "You examine $target carefully..."
```
**Purpose**: Information gathering, helps with economy.

### 1013: Cure Poison
```
Category: Healing (Abjuration)
Mana Cost: 25
Cooldown: 5 seconds
Level Required: 10
Range: Touch
Target: Self or ally
Effect: Remove poison status
Message: "The poison drains from $target's body."
```
**Purpose**: Utility/cure for poison affiliation.

### 1014: Poison
```
Category: Debuff (Necromancy)
Mana Cost: 22
Cooldown: 6 seconds
Level Required: 10
Range: 20 squares
Duration: (L * 10) seconds
Target: Single enemy
Save: Fortitude DC 14
Effect: Poison damage 1 HP per 5 seconds
Message: "$target writhes in agony as poison courses through them!"
```
**Purpose**: DoT (Damage Over Time) mechanics.

### 1015: Cure Disease
```
Category: Healing
Mana Cost: 25
Cooldown: 5 seconds
Level Required: 10
Range: Touch
Target: Self or ally
Effect: Remove disease status
Message: "The plague leaves $target."
```
**Purpose**: Utility, cure-focused.

### 1016: Curse
```
Category: Debuff (Enchantment)
Mana Cost: 28
Cooldown: 8 seconds
Level Required: 10
Range: 25 squares
Duration: 30 + (l) seconds
Target: Single enemy
Save: Will DC 15
Effects:
  - All damage done: reduced 50%
  - All damage taken: increased 25%
Message: "$actor curses $target with dark magic!"
```
**Purpose**: Strategic debuff for group play.

### 1017: Shield Spell
```
Category: Buff (Evocation)
Mana Cost: 18
Cooldown: 10 seconds
Duration: 45 + (I * 2) seconds
Level Required: 10
Target: Self
Effect: Magical barrier
- AC improvement: 3 points
- Reflect 10% of damage back to attacker
Message: "A magical barrier forms around you!"
```
**Purpose**: Active defense spell.

---

## TIER 4: EXPERT (Levels 15-19)

### 1018: Gate/Portal
```
Category: Control (Conjuration)
Mana Cost: 75
Cooldown: 60 seconds
Level Required: 15
Range: 60 squares
Target: Destination location
Effect: Open magical gateway
- Transport self
- Transport willing allies
- Stay open 30 seconds
Message: "A swirling portal opens!"
```
**Purpose**: Group transport, expensive endgame mobility.

### 1019: Summon Familiar
```
Category: Control (Conjuration)
Mana Cost: 60
Cooldown: 120 seconds
Level Required: 15
Duration: (l * 2) minutes
Target: Self
Effect: Summon magical creature to aid caster
- Familiar stats scale with caster
- Familiar attacks on caster's turn
- Familiar limited AI
Message: "$actor summons a magical creature!"
```
**Purpose**: Offense multiplier through summons.

### 1020: Dispel Magic
```
Category: Utility (Abjuration)
Mana Cost: 30
Cooldown: 5 seconds
Level Required: 15
Range: 30 squares
Target: Single character or object
Effect: Attempt to remove one active buff/enchantment
Save: Will DC 16 (caster of target effect resists)
Message: "$actor's dispel magic breaks $target's buffs!"
```
**Purpose**: Counter-play, utility.

### 1021: Telekinesis
```
Category: Utility (Transmutation)
Mana Cost: 35
Cooldown: 4 seconds
Level Required: 15
Range: 40 squares
Target: Object or creature
Effect: Lift/move object mentally
- Move light objects
- Can be used to open doors
- Can be used to trigger traps safely
Message: "$object moves on its own!"
```
**Purpose**: Puzzle solving, utility.

### 1022: Earthquake
```
Category: Offensive (Evocation)
Mana Cost: 80
Cooldown: 20 seconds
Level Required: 15
Range: Self (affects entire area)
Area: 5 square radius
Damage: "6d12 + (I * 2) + (l * 0.75)"
Save: Reflex DC 17 (half damage)
Special: 25% chance to knock down (can't move 2 rounds)
Message: "The ground shakes violently!"
```
**Purpose**: Highest damage AoE spell.

### 1023: Meteor Storm
```
Category: Offensive (Evocation)
Mana Cost: 100
Cooldown: 25 seconds
Level Required: 15
Range: 50 squares
Area: 4 square radius
Damage: "8d10 + (I * 1.5) + l"
Save: Reflex DC 18
Duration: 3 seconds (falling meteors) 
Message: "The sky erupts with falling meteors!"
```
**Purpose**: Extreme area damage, long wind-up.

### 1024: Animate Dead
```
Category: Control (Necromancy)
Mana Cost: 70
Cooldown: 180 seconds
Level Required: 15
Duration: (l * 5) minutes or until destroyed
Target: Corpse
Effect: Raise zombie to fight for caster
- Zombie stats: 50% of original
- Zombie lasts until HP depleted or duration expires
- Can control / despawn at will
Message: "A corpse rises as your undead servant!"
```
**Purpose**: Advanced summoning, requires corpse.

### 1025: Time Stop
```
Category: Utility (Transmutation)
Mana Cost: 150
Cooldown: 300 seconds
Level Required: 15
Duration: 6 + (I / 5) seconds
Target: Self
Effect: Freeze time around caster
- Caster can act normally
- All enemies frozen in place
- Can't attack frozen enemies
- Excellent for escape
Message: "Time stops for everyone except you!"
Special: Costs 1 experience per second in effect (prevents spam)
```
**Purpose**: Ultimate escape, very limited use.

---

## TIER 5: LEGENDARY (Levels 20+)

### 1026: Wish
```
Category: Utility (Universal)
Mana Cost: 200
Cooldown: 600 seconds
Level Required: 20
Target: Self or ally
Effect: Grant one special request (GM-controlled)
Message: "The universe bends to your will..."
Restrictions: Limited uses, prevents game-breaking
```
**Purpose**: Easter egg, admin tool for quest rewards.

### 1027: Divine Intervention
```
Category: Healing (Abjuration)
Mana Cost: None (Life force cost)
Cooldown: 300 seconds
Level Required: 20
Target: Dying ally
Effect: Restore ally from death with full HP
Cost: Caster loses 50% current HP
Message: "$actor sacrifices their life force to save $target!"
```
**Purpose**: Ultimate save, high cost.

### 1028: Resurrection
```
Category: Healing (Abjuration)
Mana Cost: 100
Cooldown: 120 seconds
Level Required: 20
Target: Dead player (corpse in room)
Effect: Bring dead character back to life
- Restored with 10% HP/Mana
- No experience loss
- Takes 10 seconds to cast (channeled)
Message: "$actor chants ancient words of power..."
```
**Purpose**: Death recovery, high-level utility.

---

# Spell Categories By School

## Evocation (Blast Magic)
- **1002**: Spark
- **1005**: Fireball
- **1006**: Frost Bolt
- **1008**: Lightning Strike
- **1017**: Shield Spell
- **1022**: Earthquake
- **1023**: Meteor Storm

## Abjuration (Protection & Healing)
- **1003**: Weak Heal
- **1004**: Mage Armor
- **1007**: Heal
- **1009**: Invisibility
- **1013**: Cure Poison
- **1015**: Cure Disease
- **1020**: Dispel Magic
- **1024**: Divine Intervention
- **1028**: Resurrection

## Enchantment (Control & Charm)
- **1010**: Blindness
- **1016**: Curse
- **1026**: Wish (partial)

## Transmutation (Shape-Shifting & Movement)
- **1011**: Teleport
- **1021**: Telekinesis
- **1025**: Time Stop

## Conjuration (Creation & Summoning)
- **1018**: Gate/Portal
- **1019**: Summon Familiar
- **1024**: Animate Dead

## Divination (Information)
- **1012**: Identify
- **1001**: Magic Missile (detection aspect)

## Necromancy (Death Magic)
- **1014**: Poison
- **1024**: Animate Dead

---

# Spell Categories By Game Function

## Leveling/First Spells (1-4): Learn to Cast
1001, 1002, 1003, 1004

## Combat Foundation (5-9): Learn to Fight
1005, 1006, 1007, 1008, 1009, 1010

## Utility Spells (10-14): Solve Problems
1011, 1012, 1013, 1014, 1015, 1016, 1017

## Mastery Spells (15-19): Change the Game
1018, 1019, 1020, 1021, 1022, 1023, 1024, 1025

## Legendary (20+): Endgame Power
1026, 1027, 1028

---

# Profession-Specific Spell Trees

## Mage (Intelligence-based Caster)
**All Spells**: 1001-1025 (all magic spells available)
**Bonus**: +2 mana cost reduction per spell
**Penalty**: Reduced dodge AC

## Cleric (Wisdom-based Healer)
**Spells**: 1001, 1003, 1007, 1013, 1015, 1028, 1027
**Bonus**: Healing spells: +25% effectiveness
**Penalty**: Offensive spells: -25% damage

## Ranger (Utility Caster)
**Spells**: 1002, 1006, 1009, 1011, 1012, 1014, 1021, 1025
**Bonus**: Teleport/stealth spells: -20% mana cost
**Penalty**: Area spells restricted (single-target only)

## Rogue (Stealth Caster)
**Spells**: 1001, 1002, 1009, 1011, 1012, 1018
**Bonus**: Invisibility spell grants +2 attack bonus
**Penalty**: Limited to lower-level spells

## Warrior (Combat Focus)
**Spells**: 1001, 1002, 1004, 1008, 1017
**Bonus**: None (partial caster)
**Penalty**: Limited spell selection, -25% spell damage

---

# Spell Progression Table

| Level | Spell ID | Name | Mana | Cooldown | Type |
|-------|----------|------|------|----------|------|
| 1 | 1001 | Magic Missile | 15 | 2s | Offensive |
| 1 | 1002 | Spark | 8 | 1s | Offensive |
| 1 | 1003 | Weak Heal | 12 | 3s | Healing |
| 1 | 1004 | Mage Armor | 10 | 10s | Buff |
| 5 | 1005 | Fireball | 30 | 5s | Offensive |
| 5 | 1006 | Frost Bolt | 25 | 4s | Offensive |
| 5 | 1007 | Heal | 20 | 6s | Healing |
| 5 | 1008 | Lightning Strike | 40 | 8s | Offensive |
| 5 | 1009 | Invisibility | 20 | 30s | Buff |
| 5 | 1010 | Blindness | 18 | 4s | Debuff |
| 10 | 1011 | Teleport | 50 | 30s | Control |
| 10 | 1012 | Identify | 15 | 2s | Utility |
| 10 | 1013 | Cure Poison | 25 | 5s | Healing |
| 10 | 1014 | Poison | 22 | 6s | Debuff |
| 10 | 1015 | Cure Disease | 25 | 5s | Healing |
| 10 | 1016 | Curse | 28 | 8s | Debuff |
| 10 | 1017 | Shield Spell | 18 | 10s | Buff |
| 15 | 1018 | Gate | 75 | 60s | Control |
| 15 | 1019 | Summon Familiar | 60 | 120s | Control |
| 15 | 1020 | Dispel Magic | 30 | 5s | Utility |
| 15 | 1021 | Telekinesis | 35 | 4s | Utility |
| 15 | 1022 | Earthquake | 80 | 20s | Offensive |
| 15 | 1023 | Meteor Storm | 100 | 25s | Offensive |
| 15 | 1024 | Animate Dead | 70 | 180s | Control |
| 15 | 1025 | Time Stop | 150 | 300s | Utility |
| 20 | 1026 | Wish | 200 | 600s | Utility |
| 20 | 1027 | Divine Intervention | 0 | 300s | Healing |
| 20 | 1028 | Resurrection | 100 | 120s | Healing |

---

# Damage Type Resistances

Players/creatures can have resistances/immunities to:

- **Fire**: Fireball, Fire Breath
- **Cold**: Frost Bolt, Chill Touch
- **Lightning**: Lightning Strike, Lightning Breath
- **Acid**: Acid Splash, Acid Breath
- **Magic**: All spells (magic resistance)
- **Physical**: Non-magical damage
- **Poison**: Poison spell, Poison gas
- **Holy**: Divine spells
- **Unholy**: Curse, dark magic

**Example**: Dragon with 50% fire resistance takes half damage from Fireball and Fire Breath effects.

---

# Next Steps for Implementation

1. **Expand spell definitions** in `skills/spells.json` to include 20-28 total spells
2. **Implement damage calculation engine** for formula evaluation
3. **Add spell effects/affects** system (buffs/debuffs with duration)
4. **Implement saving throw system** (DC-based)
5. **Add spell targeting modes** (self, area, hostile, etc.)
6. **Create spell descriptions** in help system
7. **Add spell learning/teaching** NPCs
8. **Implement spell proficiency progression**
9. **Create combat integration** (enemies affected by spells)
10. **Balance spell damage** based on testing

