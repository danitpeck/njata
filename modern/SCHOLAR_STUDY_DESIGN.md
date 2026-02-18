# Scholar Starter Kit: The Study System

## Vision

The Scholar is the game's **discovery-focused caster**. Rather than being granted spells, Scholars must **actively seek out and study** magical items to learn spells. This creates:

- **Exploration incentive**: Players search areas for rare magical items
- **Economic gameplay**: Magical items become valuable trade goods
- **Skill progression**: Scholars strengthen their Study skill to learn spells more reliably
- **Unique identity**: Not just "caster with different loadout," but fundamentally different gameplay loop
- **Natural progression gate**: Can't just camp casting Arcane Bolt forever—must hunt for new spells to stay viable

## Scholar Progression Model (Equipment + Skills, No Levels)

Scholars grow stronger through:

**1. Equipment Discovery** (gear loops):
- Hunt monsters → find better gear (armor, rings, etc.)
- Better gear → fight tougher monsters → find rare magical items
- Magical items often contain spells to study

**2. Spell Repertoire Expansion** (via Study):
- Find magical items (wands, scrolls, staves)
- Study them → learn new spells at 30% proficiency
- Cast spells to improve proficiency (scales damage/effectiveness)
- Higher Study skill → easier to learn harder spells

**3. Spell Proficiency Growth** (usage-based):
- Each cast slightly increases spell proficiency
- Higher proficiency = better damage/reliability
- Different spells improve at different rates

**The Loop**: Hunt → Find scalp → Study spell → Expand roster → Cast → Improve proficiency → Hunt harder areas

---

## The Study Ability

### Mechanics

```
Command: study <item>

Prerequisites:
- Item must be in the room (not inventory, like examining loot)
- Item type: WAND, STAFF, SCROLL, or FURNITURE (with spell embedded)
- Item must contain a spell (vnum-based linking)

Process:
1. Scholar examines the magical item
2. Make proficiency check: 1d100 vs (55 + Study_Proficiency * 4/5)
   - Only need to beat the roll once per item
3. Success: Learn the spell at 30% proficiency, item consumed
4. Failure: No spell learned, item consumed
5. Can't learn same spell twice

Messages:
- Start: "You study the wand carefully..."
- Success: "Knowledge floods into your mind! You have learned the art of Fireball!"
- Failure: "The magic is too complex. You cannot glean any knowledge from it."
- Either way: "The item flares brightly and crumbles to dust."
```

### Study Proficiency

**Study Skill** (Skill ID 2306 in legacy system):
```
Starting: 0% proficiency
Gain: +2.5% on successful learn, +1% on failed attempt
Max: 100%

DC Calculation (inverted):
- 55% base difficulty
- Each 1% Study reduces DC by 0.8%
- At 100%: DC = 55 - (100 * 0.8) = -25 (99% success)
- At 50%: DC = 55 - (50 * 0.8) = 15 (85% success)
- At 20%: DC = 55 - (20 * 0.8) = 39 (61% success)
```

---

## Mandatory Items for Scholar Gameplay

To support Study, the world needs magical items to discover:

### Item Types That Can Be Studied

**Wands** (ITEM_WAND)
```
value[3] = spell_id (e.g., 1002 = Fireball)
value[0] = charges (cosmetic, can be empty)
Example: "Wand of Fireball" → study it → learn Fireball
```

**Staves** (ITEM_STAFF)
```
value[3] = spell_id
value[0] = charges (cosmetic)
Example: "Staff of Healing" → study it → learn Heal
```

**Scrolls** (ITEM_SCROLL)
```
value[3] = spell_id
One-time use, consumed
Example: "Scroll of Magic Missile" → study it → learn MM
```

**Furniture** (ITEM_FURNITURE)
```
value[3] = spell_id (if 0, can't be studied)
Reusable—multiple Scholars can study same furniture
Example: "Magical Codex" in a library → multiple scholars learn from it
```

### Core Discovery Items (MVP)

The 8 MVP spells correspond to these items:

```
1. Wand of Arcane Bolt          (value[3] = 1001) - Immak libraries
2. Wand of Leviathan's Fire     (value[3] = 1002) - Immak naiad cove
3. Scroll of Mend               (value[3] = 1003) - Sedna herbalist
4. Wand of Shadow Veil          (value[3] = 1004) - Hidden/dark locations
5. Feather of Ephemeral Step    (value[3] = 1005) - Aviatrix fairy items
6. Amulet of Path Shift         (value[3] = 1006) - Courier's secret item
7. Crystal of Winter's Whisper  (value[3] = 1007) - Mountain/cold areas
8. Tome of Knowing              (value[3] = 1008) - Sedna Bronze Artisan
```

**Placement**: Items are scattered across the world, hidden in rooms or rewarded from early quests. Scholars hunt them down to expand their spell repertoire.

---

## Scholar Gameplay Loop

### Early Game (New Scholar)

1. Create Scholar, receive starting gear (basic robe, dagger)
2. Receive Arcane Bolt via training (starting spell)
3. Fight goblins in starting area → gain basic combat experience
4. Loot simple gear and first magical item (Wand of Leviathan's Fire)
5. Bring wand back to study location
6. Study → success → learn Leviathan's Fire at 30% proficiency
7. Use Leviathan's Fire in combat → proficiency gradually improves to 50%+
8. Find harder area with better gear and rare items

### Mid Game (Exploration Phase)

1. Actively hunt for specific spells needed for play style
2. Explore deeper dungeons for rare wands/scrolls
3. Trade with other players for hard-to-find items
4. Study skill approaches 50-70%, success rate increases
5. Can now tackle tougher monsters solo
6. Spell roster expands: can cast 5-6 different spells effectively
7. Each new area has new items to discover

### Late Game (Mastery Phase)

1. Seek out legendary items with powerful spells
2. Study proficiency 80%+, rarely fail
3. Spell proficiency 80%+, damage is substantial
4. Equipment at tier 3-4 (best gear available)
5. Can mentor newer Scholars or trade knowledge
6. Continually hunt for undiscovered items
7. Build unique spell combinations from exploration finds

---

## MVP Implementation (Study System Only)

### Required Components

**1. Spell ID Embedding in Items** (Already exists)
```go
// In object definition, value[3] = spell ID
type Object struct {
    Values [4]int  // value[3] for spell_id
}
```

**2. Study Command**
```go
func StudyItem(scholar *Player, itemName string) {
    // 1. Find item in room
    item := GetItemInRoom(scholar, itemName)
    
    // 2. Validate item has spell
    spellID := item.Values[3]
    if spellID == 0 { error }
    
    // 3. Make Study check
    studyDC := 55 - (scholar.Skills[STUDY] * 0.8)
    roll := Rand(1, 100)
    
    if roll < studyDC {
        // Failure
        Send(scholar, "You cannot glean knowledge from it.")
        scholar.Skills[STUDY] += 1  // +1% on failure
        item.Extract()
        return
    }
    
    // 4. Success - Learn spell
    if scholar.HasSpell(spellID) {
        Send(scholar, "You already know that spell!")
        return
    }
    
    scholar.LearnSpell(spellID, 30)  // 30% proficiency
    scholar.Skills[STUDY] += 2.5      // +2.5% on success
    item.Extract()
    Send(scholar, "You have learned the art of [spell name]!")
}
```

**3. Study Skill Tracking**
```go
type Player struct {
    Skills map[int]int  // Add STUDY skill (ID 2306)
    // skill[STUDY] = 0-100 proficiency
}
```

**4. Persistence**
```json
{
  "name": "vex",
  "skills": {
    "1001": {"proficiency": 45, "learned": true},
    "1002": {"proficiency": 30, "learned": true},
    "2306": {"proficiency": 35, "learned": true}  // STUDY skill
  }
}
```

---

## Item Creation Strategy

### MVP: 8 Core Items (Minimal)

Create object definitions for:
- Wand of Magic Missile (vnum 5001, value[3]=1001)
- Wand of Fireball (vnum 5002, value[3]=1002)
- Scroll of Healing (vnum 5003, value[3]=1003)
- Wand of Blindness (vnum 5004, value[3]=1004)
- etc.

**Where to place**:
- Some in starting area (easy discovery)
- Some in hidden rooms (incentive to explore)
- Some that respawn (repeatable learning opportunity)
- Some that are quest rewards

### Growth: Item Economy

**Week 2-3**: Add 8 more items (Wand of Poison, Scroll of Teleport variations, etc.)

**Week 3-4**: Create crafting to produce wands/scrolls

**Week 4+**: Merchants sell magical items for gold

---

## Sample Test: Scholar Gameplay

```
1. Create character: Dryad Scholar
2. Start in New Acad (Lyceum)
3. look around
4. see: "Wand of Fireball lies here"
5. study wand
   → Roll 1d100, DC 55 (Study 0%)
   → Fail? "You cannot glean..."
   → Success? "You have learned Fireball!"
6. spellbook
   → Shows: Fireball (30% proficiency)
7. go hunt goblins
8. cast fireball goblin
   → Damage dealt, proficiency increases
9. find another wand
10. study wand
    → Roll 1d100, DC 54.2 (Study now 1%)
    → Success rate improving
```

---

## Why This Design Works

### For Scholar Players
- **Active discovery**: Not boring auto-learning, but hunting for spells
- **Skill progression**: Study skill gets better over time
- **Exploration**: Incentive to visit new areas
- **Economic gameplay**: Trade wands and scrolls with other players
- **Identity**: "I discovered this spell" feels different than "I got it at level 5"

### For Game Balance
- **Gating mechanism**: Scholars can't spam powerful spells if items are rare
- **Pacing**: Gradual spell discovery vs. all spells available at once
- **Value**: Magical items worth keeping, trading, crafting
- **Difficulty curve**: Early spells easy to find, late spells require hunting

### For World Building
- **Magic feels real**: Spells aren't gifts, they're discovered knowledge
- **Lore**: Libraries, tomes, magical artifacts matter
- **Customization**: Players find different items, learn different spells
- **Uniqueness**: Two Scholars might have learned different spell selections

---

## Future Expansions (Growth Points)

### Growth 1: Spell Variations

```
Instead of one "Scroll of Heal," have multiple:
- Scroll of Weak Healing (restores 1d4 + 2)
- Scroll of Greater Healing (restores 2d8 + 5)
- Scroll of Divine Healing (restores 3d10 + 10, Cleric-only)

Scholar learns different versions with different power levels
```

### Growth 2: Spell Fusion

```
Scholar studies two wands together
"study wand1 with wand2"
→ Learns a combined spell or upgraded version
→ Requires high Study proficiency (80%+)
```

### Growth 3: Spell Creation

```
At 100% Study proficiency, Scholar can:
- Create wands from blank wands + spell knowledge
- Inscribe scrolls from vellum + spell knowledge
- Teach spells to other Scholars (trade knowledge)
```

### Growth 4: Spell Libraries

```
Specific locations with reusable furniture
- "Ancient Library in Forgotten Ruins"
- Multiple Scholars can study the same codex
- Rare and powerful spells hidden in libraries
- Quest to find the key/gain access
```

---

## Implementation Roadmap

### MVP Phase (Week 1)
- [x] Create 8 magical items with spell IDs
- [ ] Implement Study command
- [ ] Implement Study skill tracking
- [ ] Implement spell learning from study
- [ ] Test: Scholar can find and study items, learn spells
- [ ] Test: Study proficiency increases on success/failure
- [ ] Test: Can't learn same spell twice

### Polish Phase (Week 2)
- [ ] Add more items (8-12 additional)
- [ ] Place items throughout the world
- [ ] Create help text for Study command
- [ ] Test: Items respawn correctly
- [ ] Test: Study feels rewarding

### Growth Phase (Week 3+)
- [ ] Spell variations (weak/strong versions)
- [ ] Item crafting (create wands)
- [ ] NPC merchants selling items
- [ ] Study proficiency scaling
- [ ] Advanced mechanics (fusion, creation)

---

## Why Scholar First

**Scholar's place in the Equipment + Skills system:**

Unlike Warriors who gain strength through gear and combat experience, Scholars have a *unique progression gate* built into their starter kit identity:
- Gear improves like everyone else (hunt monsters → loot equipment)
- **But spells only come from discovered items**
- Can't just spam Arcane Bolt forever—must actively hunt to expand spell roster
- This natural gate keeps Scholars engaged in exploration

Example: 
- Warrior can reach mid-game power by fighting the same monsters and improving weapon proficiency
- Scholar hitting the same monsters eventually realizes "I need new spells" and must venture out to hunt items
- Different incentive systems drive different playstyles

**Three reasons to prioritize Scholar over generic spells**:

1. **Unique gameplay**: Study is completely different from auto-learn
2. **Drives item creation**: Need to create 8+ magical items anyway
3. **Custom NJATA feeling**: Not generic spells, but specific to Scholar identity

---

**Scholar + Study should be the centerpiece of the MVP, not an afterthought.**

