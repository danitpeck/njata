# Trainer System Architecture

**Goal**: NPCs teach combat maneuvers (parallel to item-based spell discovery)

---

## 1. Data Structures

### 1.1 Maneuver Definition (in skills.json)

Maneuvers are **not** spells—they're combat techniques. Store separately in a `maneuvers` array:

```json
{
  "id": 2001,
  "name": "Slash",
  "description": "Basic melee attack using your weapon. Damage scales with STR and proficiency.",
  "type": "maneuver",
  "cooldown_seconds": 2,
  "mana_cost": 0,
  "level_required": 1,
  "targeting": {
    "mode": "hostile_single",
    "range": 3,
    "radius": 0
  },
  "effects": {
    "damage": "1d6 + S/2",
    "damage_type": "physical",
    "save_type": "none",
    "save_dc": 0
  },
  "messages": {
    "cast": "You slash at $target.",
    "hit": "$actor slashes you for $damage damage!",
    "cast_room": "$actor slashes at $target."
  }
}
```

**New Maneuvers** (to add to skills.json):

```json
{
  "id": 2002,
  "name": "Power Attack",
  "description": "A devastating overhead strike. Deals 1.5x damage but leaves you vulnerable. 8 second cooldown.",
  "type": "maneuver",
  "cooldown_seconds": 8,
  "mana_cost": 0,
  "level_required": 1,
  "targeting": {
    "mode": "hostile_single",
    "range": 3,
    "radius": 0
  },
  "effects": {
    "damage": "2d6 + S",
    "damage_type": "physical",
    "save_type": "none",
    "save_dc": 0
  },
  "messages": {
    "cast": "You raise your weapon high and strike $target with all your might!",
    "hit": "$actor's devastating blow strikes you for $damage damage!",
    "cast_room": "$actor raises their weapon and strikes $target!"
  }
},
{
  "id": 2003,
  "name": "Defensive Stance",
  "description": "Adopt a defensive posture. Reduce incoming damage by 30% but your damage is reduced by 20%. Lasts 30 seconds or until you attack.",
  "type": "maneuver",
  "cooldown_seconds": 3,
  "mana_cost": 0,
  "level_required": 1,
  "targeting": {
    "mode": "self",
    "range": 0,
    "radius": 0
  },
  "effects": {
    "damage": "0",
    "damage_type": "physical",
    "affect": {
      "name": "Defensive Stance",
      "duration": "30",
      "description": "You are in a defensive stance.",
      "ac_penalty": -5,
      "stat_mods": {}
    }
  },
  "messages": {
    "cast": "You adopt a defensive stance, ready to block incoming attacks.",
    "cast_room": "$actor takes a defensive stance."
  }
},
{
  "id": 2004,
  "name": "Riposte",
  "description": "Counter-attack after dodging. Similar damage to Slash but requires timing. 5 second cooldown.",
  "type": "maneuver",
  "cooldown_seconds": 5,
  "mana_cost": 0,
  "level_required": 1,
  "targeting": {
    "mode": "hostile_single",
    "range": 3,
    "radius": 0
  },
  "effects": {
    "damage": "1d8 + D/2",
    "damage_type": "physical",
    "save_type": "none",
    "save_dc": 0
  },
  "messages": {
    "cast": "You see an opening and riposte at $target!",
    "hit": "$actor's riposte catches you off guard for $damage damage!",
    "cast_room": "$actor executes a riposte against $target!"
  }
},
{
  "id": 2005,
  "name": "Cleave",
  "description": "Swing your weapon in a wide arc, hitting all enemies in a small area. 1d4+2 damage to each. 6 second cooldown.",
  "type": "maneuver",
  "cooldown_seconds": 6,
  "mana_cost": 0,
  "level_required": 1,
  "targeting": {
    "mode": "hostile_area",
    "range": 5,
    "radius": 2
  },
  "effects": {
    "damage": "1d4 + S/4",
    "damage_type": "physical",
    "save_type": "none",
    "save_dc": 0
  },
  "messages": {
    "cast": "You swing your weapon in a wide arc!",
    "hit": "$actor's cleaving strike catches you for $damage damage!",
    "cast_room": "$actor swings their weapon in a wide, devastating arc!"
  }
}
```

---

### 1.2 Trainer Struct (in world.go)

Add to Mobile struct:

```go
type Mobile struct {
	Vnum       int
	Keywords   []string
	Short      string
	Long       string
	Race       string
	Class      string
	Position   string
	Gender     string
	Level      int
	MaxHP      int
	HP         int
	Mana       int
	MaxMana    int
	Attributes [7]int
	
	// NEW: Trainer metadata (if this mobile is a trainer)
	IsTrainer         bool   // true if this mob teaches maneuvers
	TeachesSpellID    int    // maneuver ID (e.g., 2002 for Power Attack)
	RequiredStatName  string // "Strength", "Dexterity", "Constitution", etc.
	RequiredStatValue int    // minimum value player needs in that stat
	TrainerMessage    string // custom dialog from trainer
}
```

---

### 1.3 JSON Format for Trainers (in area files)

Extend Mobile definition in area JSON with stat-specific requirements:

```json
{
  "vnum": 1401,
  "keywords": ["Nissa", "Smith", "Blacksmith"],
  "short": "Nissa the Smith",
  "long": "Nissa the Smith is here, hammering away on her anvil.",
  "race": "dwarf",
  "class": "chevalier",
  "position": "standing",
  "gender": "female",
  "level": 0,
  "max_hp": 30,
  "hp": 30,
  "mana": 0,
  "max_mana": 0,
  "attributes": [18, 13, 13, 18, 13, 13, 13],
  "room_vnum": 0,
  
  "is_trainer": true,
  "teaches_spell_id": 2002,
  "required_stat_name": "Strength",
  "required_stat_value": 18,
  "trainer_message": "Nissa says: Only the truly strong can master Power Attack!"
}
```

---

## 2. Training Command Architecture

### 2.1 Command Flow

```
> train nissa power attack
```

**Processing:**
1. Find trainer in room (by keyword match)
2. Validate trainer has `is_trainer: true`
3. Check player hasn't already learned this maneuver
4. Check if player's required stat meets minimum
5. If below → explain deficit ("You need +3 Strength")
6. If meets → grant maneuver at `default_proficiency` (10%)
7. Update proficiency tracking

### 2.2 cmdTrain Implementation (pseudocode)

```go
func cmdTrain(p *Player, args []string) {
    if len(args) < 1 {
        Send(p, "Train with whom?")
        return
    }

    // Find trainer in room
    trainer := FindMobileInRoom(p.Location, args[0])
    if trainer == nil {
        Send(p, "That person is not here.")
        return
    }

    if !trainer.IsTrainer {
        Send(fmt.Sprintf("%s doesn't teach anything.", trainer.Short))
        return
    }

    spellID := trainer.TeachesSpellID
    spell := skills.GetSpell(spellID)
    if spell == nil {
        Send(p, "ERROR: Trainer teaches invalid maneuver.")
        return
    }

    // Check if already learned
    if p.Skills[spellID] != nil {
        Send(fmt.Sprintf("You have already learned %s.", spell.Name))
        return
    }

    // Get player's stat value
    playerStatValue := GetPlayerStat(p, trainer.RequiredStatName)
    
    // Check if meets requirement
    if playerStatValue < trainer.RequiredStatValue {
        deficit := trainer.RequiredStatValue - playerStatValue
        Send(p, fmt.Sprintf(
            "%s says: You're not ready yet. You need +%d %s to learn from me.\n",
            trainer.Short, deficit, trainer.RequiredStatName))
        return
    }

    // Meets requirement → auto-success
    proficiency := 10  // Start at 10% proficiency
    p.LearnSpell(spellID, proficiency)
    Send(p, fmt.Sprintf("You learn %s!", spell.Name))
    Send(p, spell.Description)
    Broadcast(p.Location, fmt.Sprintf("%s learns %s from %s.", 
        p.Name, spell.Name, trainer.Short))
}

// Helper to get stat value by name
func GetPlayerStat(p *Player, statName string) int {
    switch strings.ToLower(statName) {
    case "strength":
        return p.Strength
    case "dexterity":
        return p.Dexterity
    case "constitution":
        return p.Constitution
    case "intelligence":
        return p.Intelligence
    case "wisdom":
        return p.Wisdom
    case "charisma":
        return p.Charisma
    case "luck":
        return p.Luck
    default:
        return 0
    }
}
```

---

## 3. Trainer Placement Strategy

| Trainer | Location | Maneuver | Required Stat | Minimum | Notes |
|---------|----------|----------|---------------|---------|-------|
| **Nissa the Smith** | Immak Blacksmith (1403) | Power Attack (2002) | **Strength** | 18 | Raw devastating power, dwarven strength |
| **Harald the Barkeeper** | Immak Pub (1422) | Riposte (2004) | **Dexterity** | 18 | Timing and reflexes, former adventurer wisdom |
| **New: Aina Fighter** | Aina Fighter's Guild (new room ~4750) | Defensive Stance (2003) | **Constitution** | 18 | Endurance and toughness, defensive focus |
| **New: Desert Nomad** | Desert Nomad Camp (new room ~31220) | Cleave (2005) | **Strength** | 20 | Even more raw power, group tactics |

**Why Different Stats?**
- **Strength-based** trainers teach power techniques (devastating single/group attacks)
- **Dexterity-based** trainers teach precision techniques (timing, counterattacks)
- **Constitution-based** trainers teach resilience techniques (defensive, mitigation)
- Encourages stat specialization: STR warrior hunts STR gear, DEX rogue hunts DEX gear, CON tank hunts CON gear

---

## 4. Implementation Checklist

### Phase 1: Data Structure
- [ ] Add 4 new maneuver definitions to skills.json (2002, 2003, 2004, 2005)
- [ ] Extend Mobile JSON schema with `is_trainer`, `teaches_spell_id`, `trainer_dc`, `trainer_message`
- [ ] Add IsTrainer, TeachesSpellID, TrainerDC, TrainerMessage fields to Mobile struct
- [ ] Update area JSON files (immak.json) with trainer metadata for existing NPCs

### Phase 2: Loading & Lookup
- [ ] Update area JSON parser to load trainer metadata from mobiles
- [ ] Update Mobile creation to include trainer fields
- [ ] Add lookup helper: `GetTrainerByMobile(mobile *Mobile) *Trainer` or similar

### Phase 3: Command Implementation
- [ ] Create cmdTrain in player_commands.go
- [ ] Register "train" command in command router
- [ ] Implement learning check (d100 vs DC, with WIS/INT modifier)
- [ ] Create LearnManeuver on Player (mirrors LearnSpell)
- [ ] Add proficiency tracking for maneuvers (same system as spells)

### Phase 4: Testing
- [ ] Test learning from existing trainers (Nissa, Harald)
- [ ] Test failed learning checks
- [ ] Test can't learn same maneuver twice
- [ ] Test proficiency improves with use (existing system)
- [ ] Test maneuver usage in combat (should already work via Slash system)

### Phase 5: Content Expansion (optional for MVP+1)
- [ ] Create "Aina Fighter's Guild" room
- [ ] Create "Desert Nomad Camp" room
- [ ] Place new trainers
- [ ] Add new maneuvers to their areas

---

## 5. Edge Cases & Design Decisions

**Q: What stat does each trainer require?**
- A: Designed to match teaching style. See trainer placement table. Different stats encourage different gear specializations.

**Q: What if player is 1 point below requirement?**
- A: Check is `playerStat >= requiredStat`. 17 STR fails vs 18 requirement. Player knows exactly what to do (find +1 STR item).

**Q: Can player learn multiple maneuvers?**
- A: Yes. Each maneuver is independent in p.Skills map. Player can be STR/Melee specialist and still learn DEX/Riposte if they find +DEX gear.

**Q: What happens if trainer NPC dies?**
- A: Trainer respawns with area reset, same as any mob. Learning is permanent once granted.

**Q: Can maneuver be learned multiple times?**
- A: No. p.Skills[spellID] != nil check prevents it. One-time learn, then proficiency improves via use.

**Q: Why is Power Attack 20 STR but Cleave also 20 STR?**
- A: Cleave is MORE powerful (AoE, more damage), so higher STR requirement makes sense. Can tune individual requirements per trainer.

**Q: Should I be able to learn trainers out of order?**
- A: Yes! If you find +STR gear early, you could learn Power Attack before Riposte. No gate-keepingjust stat-based gating.

---

## 6. Integration Points

- **Skills System**: Maneuvers use same `Spell` struct as spells (both stored in skills.json with `type` field)
- **Combat System**: Maneuvers execute same way as spells (via existing cast/damage system)
- **Proficiency System**: Players track maneuver proficiency in `p.Skills` map (same as spells)
- **Save/Load**: Player JSON already includes `Skills` map, so maneuvers persist automatically

---

## 7. Future Extensions

**Post-MVP Additions** (not in scope, but architecture supports):
- Multiple maneuvers per trainer (trainer teaches entire "school")
- Skill prerequisites ("Learn Slash before Power Attack")
- Trainer reagent costs ("10 gold + scroll of mastery")
- Trainer reputation system ("Train more to unlock advanced techniques")
- Combat chains ("Power Attack into Cleave into Riposte")

All of these can be added without restructuring core system.
