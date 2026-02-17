# NJATA Design Framework

## Quick Navigation

**Just want to build?** → Read [MVP_ROADMAP.md](MVP_ROADMAP.md)

**Want design context?** → Read below, then [SCHOLAR_STUDY_DESIGN.md](SCHOLAR_STUDY_DESIGN.md)

---

## Design Philosophy

Build small. Add meaningful. Like the original njata.c approach:
- Start with only what's needed (8 core spells + Scholar Study + Warrior Maneuvers)
- Create clean extension points for custom additions
- Grow organically based on actual gameplay needs, not pre-designed systems
- Each addition should have purpose within the game world

---

## The Two Active Design Documents

### [MVP_ROADMAP.md](MVP_ROADMAP.md) — **The Specification**

**What you're actually building in Week 1.**

- Two classes: Scholar (Study-based learning) + Warrior (Maneuver-based combat)
- 8 core spells usable by both
- Equipment + Skills progression (no levels)
- Self-reinforcing gear loops for both classes
- Full technical specifications and examples

**Read this first.** It's the authoritative spec.

---

### [SCHOLAR_STUDY_DESIGN.md](SCHOLAR_STUDY_DESIGN.md) — **The Core Feature Deep Dive**

**Detailed design of Scholar's unique mechanic.**

- Study ability mechanics (command, checks, proficiency)
- Magical item types and discovery
- Scholar gameplay loop and progression
- Why Study creates unique class identity

**Read this when:** Building the Study command or designing magical items.

---

## Archived Documents

The `_legacy_design_archive/` folder contains earlier design iterations that are superseded by the new Equipment+Skills, no-levels approach:
- IMPLEMENTATION_STRATEGY.md (old level-based design)
- SPELL_CATALOG.md (comprehensive but pre-MVP)
- SKILLS_SYSTEM_DESIGN.md (legacy framework)
- MIGRATION.md (development status from earlier phase)
- ARCHITECTURE_DIAGRAMS.md (reference only)
- VISUAL_OVERVIEW.md (old doc navigation)

These are kept for reference but **not authoritative** for current work.

---

### 3. [SKILLS_SYSTEM_DESIGN.md](SKILLS_SYSTEM_DESIGN.md) — **Reference Library**
Comprehensive design for all spell mechanics. **Don't implement from this—use it to answer questions.**

Sections:
- Skill taxonomy (8 types)
- Spell categories (offensive, defensive, healing, etc.)
- Complete spell definition schema
- Component system (for future expansion)
- Damage & healing formulas
- Saving throws system
- Targeting modes
- Effects & affects system
- Proficiency & learning
- Cooldowns & cast queue
- Failure conditions (9 types)
- Message system
- Data structures
- Implementation roadmap (Phase 1-7)

**When to use**: 
- "How should components work?" → See Part 3
- "What are all the damage types?" → See Part 11
- "How is proficiency structured?" → See Part 8
- Answering any design question about spells

---

### 4. [SPELL_CATALOG.md](SPELL_CATALOG.md) — **Spell Reference Library**
Complete spell library (28 spells, 5 tiers). **For reference and inspiration when adding new spells.**

Sections:
- Tier 1-5 spells (levels 1-20)
- 28 complete spell definitions
- Organized by school (Evocation, Abjuration, etc.)
- Organized by function (Combat, Healing, Utility, etc.)
- Profession-specific spell trees
- Spell progression tables
- Damage type resistances

**When to use**:
- "What spell should I add next?" → Browse catalog
- "What does a high-level spell look like?" → See Tier 5
- "How many spells per level?" → See progression table
- "What professions need what spells?" → See profession trees

---

### 5. [IMPLEMENTATION_STRATEGY.md](IMPLEMENTATION_STRATEGY.md) — **Technical Reference**
Technical decisions, architecture, balance framework. **For implementation decisions and balance questions.**

Sections:
- 7 key design decisions (all answered)
- Technical implementation roadmap (Milestone 1-7)
- Enhanced JSON schema
- Go implementation structure
- Full casting flow pseudocode
- Proficiency persistence
- Balance spreadsheet
- Testing strategy
- Future enhancement ideas

**When to use**:
- "How should proficiency scale?" → See Part 2
- "What's the mana efficiency target?" → See Part 4
- "How do I structure the code?" → See Part 3
- "How should I test this?" → See Part 5

---

## Document Navigation

**"I'm starting Week 1. What do I build?"**
→ Read [MVP_ROADMAP.md](MVP_ROADMAP.md) Part 1 fully  
→ Read [SCHOLAR_STUDY_DESIGN.md](SCHOLAR_STUDY_DESIGN.md) Part 2 (Mechanics)

**"I need to implement the Study command."**
→ [SCHOLAR_STUDY_DESIGN.md](SCHOLAR_STUDY_DESIGN.md) - Part 1 & 4

**"I need to understand the full system design?"**
→ [SKILLS_SYSTEM_DESIGN.md](SKILLS_SYSTEM_DESIGN.md) - Sections 1-5 gives you overview

**"I want to answer a design question?"**
→ Check [MVP_ROADMAP.md](MVP_ROADMAP.md) first (quick answers)  
→ Check [IMPLEMENTATION_STRATEGY.md](IMPLEMENTATION_STRATEGY.md) (technical decisions)  
→ Check [SKILLS_SYSTEM_DESIGN.md](SKILLS_SYSTEM_DESIGN.md) (comprehensive reference)

**"I want to add new spells after Week 1?"**
→ [SPELL_CATALOG.md](SPELL_CATALOG.md) for inspiration  
→ [SKILLS_SYSTEM_DESIGN.md](SKILLS_SYSTEM_DESIGN.md) Part 2 for schema

**"I want to balance damage numbers?"**
→ [IMPLEMENTATION_STRATEGY.md](IMPLEMENTATION_STRATEGY.md) Part 4

**"I want to plan Week 3+ features?"**
→ [MVP_ROADMAP.md](MVP_ROADMAP.md) Part 2 (Growth Points)

---

## Timeline

```
WEEK 1 (MVP) — Use MVP_ROADMAP.md + SCHOLAR_STUDY_DESIGN.md
├─ 8 spells + 8 items + Study command
├─ Get playable in 3-4 days
└─ Start testing Day 5

WEEK 2 (Observation)
├─ Play 20+ hours
├─ Watch what works/doesn't
└─ List growth priorities

WEEK 3+ (Organic Growth) — Use design docs as reference
├─ Pick #1 priority
├─ Reference docs for how to implement
├─ Code it, test it
└─ Repeat
```

---

## The Philosophy: Why This Works

### Don't Do This
```
❌ Implement all 28 spells
❌ Build components system upfront
❌ Design all professionalization systems
❌ Create damage type resistances you don't need
❌ Plan 7 phases of development before starting
```

### Do This Instead
```
✅ Build 8 spells + Scholar Study in Week 1
✅ Get it PlayABLE (not perfect)
✅ Play for 20+ hours
✅ Observe what's actually missing
✅ Add that next, not what you predicted
✅ Use design docs as reference when growing
```

### Why Scholar is the MVP Focus
1. **Unique custom content** — Not borrowed from SMAUG, specific to NJATA
2. **Drives other systems** — Need items, need placing, need item economy
3. **Different from auto-learn** — Creates achievement feeling
4. **Balances power** — Scholars progress slower (discovery-based)
5. **Feels NJATA** — Like the original njata.c (custom race-specific abilities)

---

## Success Criteria

**Week 1 MVP Success:**
- [ ] 8 spells castable
- [ ] Scholar can study items and learn spells
- [ ] Study proficiency tracks
- [ ] Non-Scholars cast normally
- [ ] No crashes in 2+ hours gameplay
- [ ] Test client passes

**Week 2 Success:**
- [ ] Played 20+ hours
- [ ] Made observations
- [ ] Have clear growth priorities

**Ongoing Success:**
- Each growth feature builds cleanly
- No regressions
- Game feels better + more fun

---

## Quick Reference Table

| Document | Size | Purpose | Start with Section |
|----------|------|---------|-------------------|
| MVP_ROADMAP.md | 3000 words | **Week 1 spec** | Part 1 |
| SCHOLAR_STUDY_DESIGN.md | 3500 words | Scholar deep dive | Part 1 (Mechanics) |
| SKILLS_SYSTEM_DESIGN.md | 7000 words | Design reference | Part 2-5 (overview) |
| SPELL_CATALOG.md | 4000 words | Spell examples | Tier 1-2 (start small) |
| IMPLEMENTATION_STRATEGY.md | 5000 words | Technical decisions | Part 2-3 (code info) |

---

## Next Step

**Ready to code?**

1. Open [MVP_ROADMAP.md](MVP_ROADMAP.md)
2. Read Part 1 completely
3. Read [SCHOLAR_STUDY_DESIGN.md](SCHOLAR_STUDY_DESIGN.md) Part 1-2
4. Start Week 1 implementation

**Questions while coding?**

1. Check what document matches your question (from table above)
2. Find the relevant section
3. Answer informed

**Week 2+ growth?**

1. Play the MVP extensively
2. Note what's missing
3. Check [MVP_ROADMAP.md](MVP_ROADMAP.md) Part 2 for growth point ideas
4. Reference relevant design doc to implement

---

**This framework separates:**
- **What you build now** (MVP_ROADMAP + SCHOLAR_STUDY)
- **How you build it** (IMPLEMENTATION_STRATEGY)
- **What to build later** (SPELL_CATALOG, SKILLS_SYSTEM_DESIGN)

**No information overload. No decision fatigue. Just build Week 1.**

