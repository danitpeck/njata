# NJATA Skills System: Visual Overview

## The Complete Picture

```
                     NJATA SKILLS SYSTEM
                          
        ┌─────────────────────────────────────┐
        │      DESIGN LIBRARY (Reference)     │
        │   Use when answering design Qs      │
        ├─────────────────────────────────────┤
        │  • SKILLS_SYSTEM_DESIGN.md (15 sec)│
        │  • SPELL_CATALOG.md (28 spells)    │
        │ • IMPLEMENTATION_STRATEGY.md (tech) │
        └─────────────────────────────────────┘
                         ▲
                         │
        ┌────────────────┴────────────────┐
        │                                 │
        │    GROW ORGANICALLY             │
        │   (Week 2-8+)                   │
        │                                 │
        │  • Play extensively             │
        │  • Observe what's missing       │
        │  • Reference design docs        │
        │  • Add next feature             │
        │                                 │
        └────────────────┬────────────────┘
                         │
        ┌────────────────▼────────────────┐
        │                                 │
        │    MVP LAYER (Week 1)           │
        │   What you build first          │
        │                                 │
        │  • 8 spells                     │
        │  • 8 magical items              │
        │  • Scholar Study (core feature) │
        │  • Basic casting                │
        │                                 │
        │  USE THESE DOCS:                │
        │  • MVP_ROADMAP.md               │
        │  • SCHOLAR_STUDY_DESIGN.md      │
        │                                 │
        └─────────────────────────────────┘
```

---

## Week 1: The Four Documents You Need

### 1️⃣ MVP_ROADMAP.md
**"Tell me exactly what to build"**

```
Part 1: MVP Phase
├─ What goes in (8 spells, Study system)
├─ What stays out (components, effects, etc.)
├─ Data structures (spell JSON, item JSON)
├─ Go code examples
└─ MVP testing checklist

Part 2: Growth Points
├─ When to add components
├─ When to add effects
├─ When to add more spells
└─ How organic growth works

Part 3: Design for Growth
└─ Why the architecture supports growth
```

**Action**: Read Part 1, start coding.

---

### 2️⃣ SCHOLAR_STUDY_DESIGN.md
**"Build the core Scholar feature"**

```
Part 1: Scholar Vision
├─ Why Study exists
├─ Creates exploration incentive
├─ Creates item economy
├─ Unique NJATA content

Part 2: The Study Mechanic
├─ Find item in room
├─ Study skill check (DC = 55 - prof*0.8)
├─ Success: Learn spell at 30%
├─ Failure: No spell, item consumed
├─ Study proficiency improves with use

Part 3-4: Implementation Details
├─ Code examples
├─ Item types (wand, staff, scroll, furniture)
├─ 8 MVP items to create
└─ Study command pseudocode

Part 5: Item Strategy
├─ Where to place items
├─ How to make them respawn
├─ How to create item economy
└─ Future enhancements (fusion, crafting)
```

**Action**: Build the Study command and 8 items.

---

### 3️⃣ IMPLEMENTATION_STRATEGY.md (Skim Key Sections)
**"How do I implement this in Go?"**

Key sections:
- Part 3: Enhanced JSON schema (what fields spells need)
- Part 3.2: Go implementation structure (package organization)
- Part 3.3: Casting flow pseudocode (the main loop)
- Part 4: Balance spreadsheet (mana efficiency targets)

**Action**: Reference when implementing casting logic.

---

### 4️⃣ SPELL_CATALOG.md (Skim Section 1)
**"What are example spells?"**

Section: Tier 1 Spells (Levels 1-5)
- 8 detailed spell examples
- Copy their structure for your 8 MVP spells

**Action**: Use as template for spell definitions.

---

## What NOT to Do in Week 1

```
❌ Don't read SPELL_CATALOG.md all 28 spells
   (Use it later for inspiration)

❌ Don't read SKILLS_SYSTEM_DESIGN.md in full
   (Way too much detail for MVP)

❌ Don't try to implement components system
   (Store reference, skip for now)

❌ Don't design professionalization
   (Focus on Scholar only)

❌ Don't plan Phase 2-7 in advance
   (You'll know what's missing after playing)
```

---

## The 8 MVP Spells

```
1001: Magic Missile
      └─ Wand of Magic Missile (vnum 5001, value[3]=1001)

1002: Fireball
      └─ Wand of Fireball (vnum 5002, value[3]=1002)

1003: Heal
      └─ Scroll of Healing (vnum 5003, value[3]=1003)

1004: Blindness
      └─ Wand of Blindness (vnum 5004, value[3]=1004)

1005: Invisibility
      └─ Potion of Invisibility (vnum 5005, value[3]=1005)

1006: Teleport
      └─ Wand of Teleportation (vnum 5006, value[3]=1006)

1007: Frost Bolt
      └─ Wand of Frost Bolt (vnum 5007, value[3]=1007)

1008: Identify
      └─ Scroll of Identify (vnum 5008, value[3]=1008)
```

Each spell:
- Has a name, mana cost, cooldown
- Has a damage formula: "XdY + I" (plus attributes)
- Has messages (cast, hit, miss, save)
- Makes saving throw if applicable
- Has one item Scholars can find and study

---

## The Scholar Study Loop

```
┌──────────────────────────────────────────┐
│ Scholar Creates Character (Dryad)        │
└──────────────┬───────────────────────────┘
               │
               ▼
┌──────────────────────────────────────────┐
│ Explore World, Find Wand of Fireball    │
└──────────────┬───────────────────────────┘
               │
               ▼
┌──────────────────────────────────────────┐
│ Command: study wand                      │
├──────────────────────────────────────────┤
│ Roll 1d100 vs DC (55 - Study_Prof * 0.8)│
│  • 0% Study proficiency → DC 55          │
│  • Roll 1d100, need beat 55              │
│  • 55%+ chance of success (1st attempt)  │
└──────────────┬───────────────────────────┘
               │
     ┌─────────┴──────────┐
     │                    │
     ▼ SUCCESS            ▼ FAILURE
  (55%+)              (<55%)
     │                    │
     ├─ Learn spell  ├─ No spell learned
     ├─ 30% prof     ├─ Study +1%
     ├─ +2.5% Study  ├─ Item destroyed
     ├─ Item gone    └─ Learn from failure
     └─ Message
     
               │
               ▼
┌──────────────────────────────────────────┐
│ Cast Fireball in Combat                  │
├──────────────────────────────────────────┤
│ • 30 mana cost                           │
│ • 5 second cooldown                      │
│ • 4d8 base damage + INT bonus            │
│ • Target's Reflex save for half damage   │
│ • +1.5% proficiency per cast             │
└──────────────┬───────────────────────────┘
               │
               ▼
     ┌─────────────────────────┐
     │ Find More Wands         │
     │ Study More              │
     │ Proficiency Improves    │
     │ Become Spell Master     │
     └──────────────┬──────────┘
                    │
                    ▼
           (Repeat forever)
```

---

## The Design Reference Flow

```
┌────────────────────────────────────────────┐
│ You're Coding Week 1, Question Arises      │
└──────────────┬─────────────────────────────┘
               │
      ┌────────┴────────┐
      │                 │
      ▼                 ▼
"How do I structure    "What should
 the code?"           Blindness spell do?"
      │                 │
      ▼                 ▼
IMPLEMENTATION_         SPELL_CATALOG
STRATEGY.md             .md
Part 3                  Tier 1
(Go structure)          (Example spells)
      │                 │
      ▼                 ▼
  Code it            Copy Schema


     ┌────────────────────────┐
     │ More Complex Question  │
     └──────────┬─────────────┘
                │
      ┌─────────┴──────────┐
      │                    │
      ▼                    ▼
 "How do saving          "What's proficiency?"
  throws work?"
      │                    │
      ▼                    ▼
 MVP_ROADMAP           IMPLEMENTATION_STRATEGY
 Part 1                Part 2
 (Quick answer)        (Full answer)
      │                    │
      ▼                    ▼
   Implement            Understand System


     ┌────────────────────────┐
     │ Design Question        │
     │ (Not about code)       │
     └──────────┬─────────────┘
                │
      ┌─────────┴──────────┐
      │                    │
      ▼                    ▼
  "Should spells have    "How many tiers
   components?"          of spells?"
      │                    │
      ▼                    ▼
  SKILLS_SYSTEM_       SPELL_CATALOG
  DESIGN.md            .md
  Part 3               Part 1-5
  (Components)         (Spell progression)
      │                    │
      ▼                    ▼
  Answer: No (for MVP)  Answer: Check Part 2
  Store for later       (8 spells in Tier 1-2)
```

---

## File Sizes (What to Expect)

```
MVP_ROADMAP.md              ~3,000 words   Easy read (1hr)
SCHOLAR_STUDY_DESIGN.md     ~3,500 words   Easy read (1hr)
IMPLEMENTATION_STRATEGY.md  ~5,000 words   Reference (browse sections)
SKILLS_SYSTEM_DESIGN.md     ~7,000 words   Full reference (scan for sections)
SPELL_CATALOG.md            ~4,000 words   Skim Tier 1 (30min)
FRAMEWORK.md                ~2,000 words   Navigation aid (15min)
```

**Total read time for Week 1 prep: ~3-4 hours**
**Then start coding!**

---

## Success Checklist

### Before You Start
- [ ] Read MVP_ROADMAP.md Part 1 fully
- [ ] Read SCHOLAR_STUDY_DESIGN.md Part 1-2
- [ ] Skim IMPLEMENTATION_STRATEGY.md Part 3
- [ ] Skim SPELL_CATALOG.md Tier 1 spells
- [ ] Understand the Study mechanic
- [ ] Understand the 8 MVP spells

### During Development
- [ ] Check SPELL_CATALOG.md for spell structure
- [ ] Reference IMPLEMENTATION_STRATEGY.md for code layout
- [ ] Use SKILLS_SYSTEM_DESIGN.md only if stuck on design question
- [ ] Don't try to implement growth features

### After MVP Works
- [ ] Play for 20+ hours
- [ ] Make observations
- [ ] Check MVP_ROADMAP.md Part 2 for growth points
- [ ] Decide which growth point to tackle next
- [ ] Reference relevant design doc for implementation

---

## Remember

```
Week 1: Small, focused scope
        8 spells + Scholar Study
        Get it PLAYABLE
        
Week 2: Play and observe
        What works?
        What's missing?
        
Week 3+: Grow based on observations
        Not predictions
        Real player experience
        
Design Docs: Reference library
        Not a spec
        Answer questions as they come up
        Check them when you need answers
```

**Start small. Add meaningful. Keep NJATA unique.**

