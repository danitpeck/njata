````markdown
# NJATA UI Design Framework

**Phase**: Post-MVP Visual Polish | **Priority**: Medium | **Effort**: 3-4 weeks

---

## Vision

Transform NJATA from functional text output to **visually polished classic MUD interface** that feels alive and immersive. The reference is the original njata.c: clean layout, color-coded information, spatial awareness, and persistent status display.

**Design Goal**: Player can instantly understand:
- Where they are (room name + exits)
- What's happening (combat, NPCs, objects)
- How they're doing (health/mana at a glance)
- What to do next (prompt is clear)

---

## Phase Overview

**Phase 3: UI Polish** (Post-MVP, parallel with Week 2-3 gameplay testing)

- Redesign room display (exits, description, NPCs/objects, state)
- Create output formatter system (colors, alignment, templates)
- Build status bar (health/mana/endurance, always visible)
- Polish combat messages (clear, color-coded, readable)
- Add visual feedback (action results, state changes)
- Create help/info display system

**Dependent On**: Nothing—can start after MVP foundation works

**Blocks**: Smooth gameplay feel (but not core functionality)

---

## Part 1: Output Formatter Architecture

### Current State (MVP)
```
Simple text -> raw output to client
No formatting, no structure, walls of text
```

### Target State
```
Formatted output with:
- ANSI color codes
- Line wrapping
- Alignment (centered headers, indented lists)
- Message templates (combat, status, errors)
- Visual separators
```

### Core System: Output Formatter

```go
// internal/game/output_formatter.go

type OutputFormatter interface {
    // Low-level rendering
    ColorCode(color string) string
    AlignCenter(text string, width int) string
    DrawLine(width int, char string) string
    
    // Mid-level templates
    RoomHeader(roomName string, exits []string) string
    RoomDescription(description string) string
    MobPresence(mobs []string) string
    ObjectPresence(objects []string) string
    
    // High-level complete displays
    RoomDisplay(room *Room) string
    StatusBar(player *Player) string
    CombatMessage(action string, actor, target string, damage int) string
}

// Standard implementation
type DefaultFormatter struct {
    Width int  // Terminal width (80 by default)
}

// Usage in session output
func (s *Session) Output(format string, args ...interface{}) {
    formatted := s.formatter.Format(format, args...)  // Apply colors/alignment
    s.Send(formatted)
}
```

### Color Constants

```go
// internal/text/colors.go - Extend existing

const (
    // Named colors (semantic meaning)
    ColorExit     = Green       // Exits
    ColorCommand  = Cyan        // Player commands
    ColorDamage   = Red         // Damage taken
    ColorHealing  = BrightGreen // Healing received
    ColorStatus   = White       // Important info
    ColorAmbient  = Dim         // Descriptions
    ColorError    = BrightRed   // Errors
    ColorSuccess  = BrightGreen // Success messages
    
    // ANSI codes
    Green       = "\x1b[32m"
    BrightGreen = "\x1b[92m"
    Red         = "\x1b[31m"
    BrightRed   = "\x1b[91m"
    Cyan        = "\x1b[36m"
    White       = "\x1b[37m"
    Dim         = "\x1b[90m"
    Bold        = "\x1b[1m"
    Reset       = "\x1b[0m"
)
```

---

## Part 2: Display Elements

### Element 1: Room Header (Exits + Name)

**Current Output:**
```
You are in a sunny field.
Exits: north, east, west
```

**Target Output:**
```
[Exits: North East West]
The sunny field
```

**Implementation:**

```go
func (f *DefaultFormatter) RoomHeader(room *Room) string {
    // Build exit line
    var exits []string
    if _, ok := room.Exits["north"]; ok { exits = append(exits, "North") }
    if _, ok := room.Exits["east"]; ok { exits = append(exits, "East") }
    if _, ok := room.Exits["south"]; ok { exits = append(exits, "South") }
    if _, ok := room.Exits["west"]; ok { exits = append(exits, "West") }
    if _, ok := room.Exits["up"]; ok { exits = append(exits, "Up") }
    if _, ok := room.Exits["down"]; ok { exits = append(exits, "Down") }
    
    exitLine := fmt.Sprintf("[Exits: %s]\n", strings.Join(exits, " "))
    if len(exits) == 0 {
        exitLine = "[Exits: None]\n"
    }
    
    // Room name line
    nameLine := fmt.Sprintf("%-s%s%s\n", 
        f.ColorCode(ColorStatus), 
        room.Name, 
        f.ColorCode(Reset))
    
    return exitLine + nameLine
}
```

**Visual Example:**
```
[Exits: North East West]
The sunny field
```

### Element 2: Room Description (Prose + Flavor)

**Display Strategy:**
- Prose description in default/dim color
- Line width at 75-80 chars for readability
- Empty line after description

**Implementation:**

```go
func (f *DefaultFormatter) RoomDescription(description string) string {
    // Wrap at 75 chars, maintain paragraph breaks
    wrapped := f.WordWrap(description, 75)
    
    // Add color (default/ambient)
    return f.ColorCode(ColorAmbient) + wrapped + f.ColorCode(Reset) + "\n"
}

// Helper: Word wrapping
func (f *DefaultFormatter) WordWrap(text string, width int) string {
    // Standard word wrapping algorithm
    // Preserve \n for paragraph breaks
}
```

**Visual Example:**
```
You are standing in the middle of a wide, summery, sunlit field. There is a 
fragrance of spring in the air, a sound of summer and a feeling of eternal 
Saturday afternoon. To the south is a clear, sparkling lake, and to the north 
and west is the holy grove in the wood's edge, to the east is a staticy mansion,
shimmering softly through the colors of the rainbow.
```

### Element 3: NPC/Mob Presence

**Display Strategy:**
- List NPCs with their status
- Important ones highlighted
- Format: `Name [status] - description`

**Implementation:**

```go
func (f *DefaultFormatter) MobPresence(mobs []*Mobile) string {
    if len(mobs) == 0 {
        return ""
    }
    
    var lines []string
    for _, mob := range mobs {
        health := "unharmed"
        if mob.HP < mob.MaxHP/2 {
            health = "[wounded]"
        }
        if mob.HP < mob.MaxHP/4 {
            health = "[critical]"
        }
        
        line := fmt.Sprintf("%s %s - %s%s",
            mob.ShortDescription,
            f.ColorCode(ColorStatus) + health + f.ColorCode(Reset),
            mob.LongDescription,
            Reset)
        
        lines = append(lines, line)
    }
    
    return strings.Join(lines, "\n") + "\n"
}
```

**Visual Example:**
```
The Hierophant glows with an aura of divine radiance.
The hiero is here, munching on grass.
```

### Element 4: Objects in Room

**Display Strategy:**
- List items that can be picked up
- Format: `A noun lies here.` or `Nouns are here.`
- Simple, inventory-focused

**Implementation:**

```go
func (f *DefaultFormatter) ObjectPresence(objects []*Object) string {
    if len(objects) == 0 {
        return ""
    }
    
    // Group by vnum for plurality
    grouped := make(map[int]int)
    var singles []*Object
    
    for _, obj := range objects {
        if obj.Count > 1 {
            grouped[obj.Vnum]++
        } else {
            singles = append(singles, obj)
        }
    }
    
    var lines []string
    
    // Plurals
    for vnum, count := range grouped {
        // Assume we have object template
        obj := GetObjectTemplate(vnum)
        pluralName := obj.PluralName // e.g. "gold coins"
        lines = append(lines, fmt.Sprintf("%d %s lie here.", count, pluralName))
    }
    
    // Singles
    for _, obj := range singles {
        lines = append(lines, fmt.Sprintf("A %s lies here.", obj.ShortDescription))
    }
    
    return strings.Join(lines, "\n") + "\n"
}
```

**Visual Example:**
```
A wand of Leviathan's Fire lies here.
3 gold coins lie here.
```

### Element 5: Status Bar (Health/Mana/Endurance)

**Display Strategy:**
- Single line, always visible after each output
- Compact format: `Health: X/Y Mana: X/Y Endurance: X/Y`
- Color-coded (red if low, yellow if mid, green if full)
- Optional: Add bar visualization

**Implementation:**

```go
func (f *DefaultFormatter) StatusBar(player *Player) string {
    healthColor := f.HealthColor(player.HP, player.MaxHP)
    manaColor := f.ManaColor(player.Mana, player.MaxMana)
    enduranceColor := f.EnduranceColor(player.Stamina, player.MaxStamina)
    
    bar := fmt.Sprintf("%s<Health: %d/%d%s %s<Mana: %d/%d%s %s<Endurance: %d/%d%s",
        healthColor, player.HP, player.MaxHP, Reset,
        manaColor, player.Mana, player.MaxMana, Reset,
        enduranceColor, player.Stamina, player.MaxStamina, Reset)
    
    return bar
}

// Color based on percentage
func (f *DefaultFormatter) HealthColor(current, max int) string {
    pct := (current * 100) / max
    if pct <= 10 { return BrightRed }    // Critical
    if pct <= 30 { return Red }          // Low
    if pct <= 70 { return Yellow }       // Medium
    return BrightGreen                   // Good
}
```

**Visual Example (Low Health):**
```
<Health: 250/1000 <Mana: 1000/1000 <Endurance: 991/1000
```

### Element 6: Health Bar Visualization (Optional)

**Display Strategy:**
- After combat damage: show graphical HP bar
- Format: `[||||||||   ] 80%` or `[████████░░] 80%`
- Compact, informative, shows at-a-glance state

**Implementation:**

```go
func (f *DefaultFormatter) HealthBarVisual(current, max int) string {
    barLength := 20
    filled := (current * barLength) / max
    
    var bar strings.Builder
    bar.WriteString("[")
    for i := 0; i < barLength; i++ {
        if i < filled {
            bar.WriteString("█")
        } else {
            bar.WriteString("░")
        }
    }
    bar.WriteString("] ")
    
    pct := (current * 100) / max
    bar.WriteString(fmt.Sprintf("%d%%", pct))
    
    return bar.String()
}
```

**Visual Example:**
```
[████████░░░░░░░░░░░░] 35%
```

### Element 7: Combat Message System

**Display Strategy:**
- Clear, action-focused
- Color-coded (red for damage taken, cyan for your actions, green for healing)
- Timestamp optional
- Context (who did what to whom)

**Message Templates:**

```
Your action (Cyan):
  You slash the goblin for 12 damage.

Enemy action (Red):
  The goblin slashes you for 8 damage.

Status change (Yellow):
  The goblin is wounded!

Result (Green):
  The goblin dies.
```

**Implementation:**

```go
// internal/game/combat_messages.go

type CombatMessage struct {
    Type    string  // "attack", "heal", "status", "death"
    Actor   string
    Target  string
    Action  string
    Damage  int
}

func (f *DefaultFormatter) CombatMessage(msg CombatMessage) string {
    var output string
    
    switch msg.Type {
    case "attack":
        color := ColorDamage
        if msg.Damage < 0 { color = ColorHealing }
        
        output = fmt.Sprintf("%s%s strikes %s for %d damage.%s",
            color, msg.Actor, msg.Target, msg.Damage, Reset)
    
    case "your_attack":
        output = fmt.Sprintf("%sYou %s %s for %d damage.%s",
            ColorCommand, msg.Action, msg.Target, msg.Damage, Reset)
    
    case "enemy_attack":
        output = fmt.Sprintf("%s%s %s you for %d damage.%s",
            ColorDamage, msg.Actor, msg.Action, msg.Damage, Reset)
    
    case "healing":
        output = fmt.Sprintf("%s%s heals you for %d HP.%s",
            ColorHealing, msg.Actor, msg.Damage, Reset)
    
    case "status":
        output = fmt.Sprintf("%s%s is %s.%s",
            ColorStatus, msg.Target, msg.Action, Reset)
    
    case "death":
        output = fmt.Sprintf("%s%s dies.%s",
            ColorSuccess, msg.Target, Reset)
    }
    
    return output
}
```

---

## Part 3: Complete Room Display

**Putting It All Together**:

```go
func (f *DefaultFormatter) RoomDisplay(room *Room, player *Player) string {
    var output strings.Builder
    
    // 1. Exits + Room name
    output.WriteString(f.RoomHeader(room))
    output.WriteString("\n")
    
    // 2. Description
    output.WriteString(f.RoomDescription(room.Description))
    output.WriteString("\n")
    
    // 3. NPCs/Mobs
    if len(room.Mobiles) > 0 {
        output.WriteString(f.MobPresence(room.Mobiles))
        output.WriteString("\n")
    }
    
    // 4. Objects
    if len(room.Objects) > 0 {
        output.WriteString(f.ObjectPresence(room.Objects))
        output.WriteString("\n")
    }
    
    // 5. Status bar (after each look)
    output.WriteString(f.StatusBar(player))
    output.WriteString("\n")
    
    return output.String()
}
```

**Visual Output**:
```
[Exits: North East West]
The sunny field

You are standing in the middle of a wide, summery, sunlit field. There is a 
fragrance of spring in the air, a sound of summer and a feeling of eternal 
Saturday afternoon. To the south is a clear, sparkling lake, and to the north 
and west is the holy grove in the wood's edge, to the east is a staticy mansion,
shimmering softly through the colors of the rainbow.

The Hierophant glows with an aura of divine radiance.
The hiero is here, munching on grass.

A wand of Leviathan's Fire lies here.
3 gold coins lie here.

<Health: 1000/1000 <Mana: 1000/1000 <Endurance: 989/1000
```

---

## Part 4: Message Types & Examples

### Action Messages (Player Does Something)

```
You cast Arcane Bolt at the goblin.
You dodge to the side.
You study the wand carefully.
```

### Combat Messages (Interaction)

```
Your Arcane Bolt hits the goblin for 12 damage!
The goblin slashes you for 8 damage.
The goblin dies.
```

### Status Messages (State Changes)

```
You gain 1% Arcane Bolt proficiency.
The goblin blocks your attack.
You restore full health.
```

### Error Messages (Problems)

```
You don't see that target.
Not enough mana.
Your spell is on cooldown.
```

### Information Messages (Requests)

```
Known Spells (5):
  1. Arcane Bolt (45% proficiency)
  2. Leviathan's Fire (30% proficiency)
  ...
```

---

## Part 5: Input/Output Flow

### Current Session Loop (MVP)

```
1. Player connects
2. Player sees prompt "# "
3. Player types command
4. Handle command
5. Output result (plain text)
6. Go to step 2
```

### Improved Session Loop (with UI)

```
1. Player connects
2. Show welcome + room display (formatted)
3. Show status bar
4. Show prompt [> ]
5. Player types command
6. Handle command
7. Show formatted output (with colors, alignment)
8. Show status bar (if changed)
9. Go to step 4
```

**Code Example:**

```go
func (s *Session) RunGameLoop() {
    // Initial room display
    room := s.player.Location
    s.Send(s.Formatter.RoomDisplay(room, s.player))
    
    for {
        // Prompt
        s.Send(ColorCommand + "[> ]" + Reset)
        
        // Read input
        line, err := s.ReadLine()
        if err != nil { break }
        
        // Execute
        s.HandleCommand(line)
        
        // Output formatter handles coloring
        // Status bar shows automatically
    }
}
```

---

## Part 6: Implementation Roadmap

### Phase 3.1: Foundation (Week 1)

- [ ] Create OutputFormatter interface
- [ ] Implement color constants
- [ ] Build RoomHeader() and RoomDescription()
- [ ] Create StatusBar()
- [ ] Hook into Session.Send() for formatting
- [ ] Test: Room display looks correct

### Phase 3.2: Combat (Week 1-2)

- [ ] Implement CombatMessage() system
- [ ] Add colored damage/healing output
- [ ] Create health bar visualization
- [ ] Update combat command output
- [ ] Test: Combat messages are clear and color-coded

### Phase 3.3: Polish (Week 2)

- [ ] Add MobPresence() and ObjectPresence()
- [ ] Create message templates for all action types
- [ ] Add help display formatting
- [ ] Implement word wrapping utility
- [ ] Test: All outputs are readable and aligned

### Phase 3.4: Iteration (Week 3)

- [ ] Gather feedback from gameplay
- [ ] Adjust colors if needed
- [ ] Add optional features (health bar, etc.)
- [ ] Performance testing (ensure no lag from formatting)

---

## Part 7: Extensible Design Points

### Adding New Message Types

**Process**:
1. Define new type in CombatMessage enum
2. Add case statement to CombatMessage()
3. Use existing color constants
4. Test output

**Example** (adding critical hit):
```go
case "critical":
    output = fmt.Sprintf("%s*** CRITICAL HIT! ***%s %s strikes %s for %d damage!",
        BrightRed, Reset, msg.Actor, msg.Target, msg.Damage)
```

### Adding New Colors

**Process**:
1. Add ANSI code constant to colors.go
2. Define semantic color name (ColorBuff, ColorDebuff, etc.)
3. Use in formatter

**Example** (adding buff indicator):
```go
const (
    ColorBuff    = BrightGreen
    ColorDebuff  = BrightRed
    ColorNeutral = Cyan
)
```

### Custom Formatting Per Object

**Process**:
1. Objects can have DisplayFormatter method
2. Session uses object's formatter if available
3. Falls back to default formatter

**Example** (boss mob with special display):
```go
func (m *BossMob) Format(f OutputFormatter) string {
    return fmt.Sprintf("!!! %s !!! [%d/%d HP]",
        m.Name, m.HP, m.MaxHP)
}
```

---

## Part 8: Compatibility Notes

### Terminal Requirements

- **Width**: 80 characters (standard, adjustable)
- **Height**: No hard requirement (dynamic wrapping)
- **Colors**: ANSI 256-color support (most terminals, last 20 years)
- **Special**: Assumes `\x1b[` escape sequences work

### Telnet Client Support

- Most modern telnet clients support ANSI colors
- If client disables colors: text still readable (reset codes ignored)
- Alignment preserved even without color support

### Web Client (Future)

- ANSI codes convert to HTML span tags
- Colors map to CSS classes
- Alignment preserved with responsive design

---

## Part 9: Testing Strategy

### Unit Tests

```go
func TestColorCode(t *testing.T) {
    f := NewDefaultFormatter()
    code := f.ColorCode(ColorRed)
    assert.Equal(t, "\x1b[31m", code)
}

func TestStatusBar(t *testing.T) {
    p := &Player{HP: 50, MaxHP: 100, Mana: 75, MaxMana: 100}
    bar := f.StatusBar(p)
    assert.Contains(t, bar, "50/100")
    assert.Contains(t, bar, "75/100")
}

func TestRoomHeader(t *testing.T) {
    room := &Room{
        Name: "Test Room",
        Exits: map[string]int{"north": 123, "east": 456},
    }
    header := f.RoomHeader(room)
    assert.Contains(t, header, "North")
    assert.Contains(t, header, "East")
    assert.NotContains(t, header, "South")
}
```

### Integration Tests

```
1. Player connects -> room display has colors
2. Player casts spell -> message is color-coded
3. Player takes damage -> health bar updates
4. Player heals -> green healing message
5. Mob dies -> death message in bright red
```

### Manual Testing

```
1. Connect with telnet/MUD client
2. Check colors display correctly
3. Check alignment looks good at 80 chars
4. Check with color-blind enabled (if possible)
5. Check on different terminal emulators
```

---

## Part 10: Success Criteria

**Phase 3 Complete When:**
- [ ] Room display is formatted and color-coded
- [ ] Status bar appears after each action
- [ ] Combat messages are clear and readable
- [ ] Health bar visualization works (optional)
- [ ] No performance degradation vs MVP
- [ ] Tested on 3+ terminal clients
- [ ] Feels closer to original njata.c UI
- [ ] Player feedback: "Much better to look at!"

---
## Part 9b: Visual Reference

**See**: `lore/njata-gameplay.bmp` for the original njata.c UI in action.

This screenshot is the target aesthetic for Phase 3—study it while designing output formatting.

---
## Reference: Original njata.c UI Style

**What Made It Good:**
- ✅ Exits always visible and formatted
- ✅ Clear room descriptions with prose
- ✅ Important info (NPCs, objects) highlighted
- ✅ Combat feedback immediate and clear
- ✅ Status always visible (no scrolling up to check)
- ✅ Color-coding reduced cognitive load
- ✅ Visual hierarchy (important info bright, flavor dim)

**NJATA Will Match**:
- ✅ Same exits display style
- ✅ Prose room descriptions
- ✅ Mob/object presence clear
- ✅ Status bar always visible
- ✅ Color-coded messages
- ✅ Clean, readable output

---

## Next Steps

**After MVP is Stable (Week 2)**:
1. Read Part 1-4 of this document
2. Create OutputFormatter interface
3. Build default formatter
4. Hook into Session output
5. Test with telnet client
6. Iterate based on feel

**Timeline**: 3-4 weeks post-MVP, done in parallel with gameplay testing.

**Why Separate Phase**: UI is polish, not core gameplay. But it makes the game feel professional and is worth a dedicated phase.

````
