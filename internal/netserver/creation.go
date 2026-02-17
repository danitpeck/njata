package netserver

import (
	"fmt"
	"strconv"
	"strings"

	"njata/internal/game"
	"njata/internal/kits"
	"njata/internal/races"
	"njata/internal/skills"
)

// Age categories
const (
	AgeChild = iota
	AgeYouth
	AgeAdult
	AgeMiddleAged
	AgeElderly
)

// Sex categories
const (
	SexNeutral = iota
	SexMale
	SexFemale
)

var ageNames = []string{"Child", "Youth", "Adult", "Middle-Aged", "Elderly"}
var ageDescriptions = []string{
	"Children are weak and young, barely experienced in the world and physically overwhelmed easily. Though they are extremely lucky, they are not very skilled.",
	"Youths are exuberant and physically fit, but not very experienced in the world, and rarely skilled.",
	"Adults are balanced, well-experienced but still in their physical prime. Though they are somewhat skilled, they still have much to learn in life.",
	"The middle-aged are experienced individuals who are highly skilled. They may not be quite as physically impressive as their younger counterparts, but they are still quite capable.",
	"While extremely skilled and usually of great wisdom or intelligence, the elderly are nearer to the end of their life-cycle and suffer from the beginnings of the failing of a mortal body.",
}

var sexNames = []string{"Neutral", "Male", "Female"}
var sexDescriptions = []string{
	"Neutral - Beyond the boundaries of typical gender, your essence transcends the physical form.",
	"Male - A masculine identity resonates within your being.",
	"Female - A feminine essence defines your presence.",
}

// CharacterCreation handles the character creation flow for new players
type CharacterCreation struct {
	session      *Session
	player       *game.Player
	selectedRace *races.RaceJSON
	selectedKit  *kits.StarterKit
}

// NewCharacterCreation creates a new character creation session
func NewCharacterCreation(session *Session, player *game.Player) *CharacterCreation {
	return &CharacterCreation{
		session: session,
		player:  player,
	}
}

// Run executes the full character creation flow
func (cc *CharacterCreation) Run() error {
	cc.displayWelcome()

	if err := cc.selectRace(); err != nil {
		return err
	}

	if err := cc.selectStarterKit(); err != nil {
		return err
	}

	if err := cc.selectAge(); err != nil {
		return err
	}

	if err := cc.selectSex(); err != nil {
		return err
	}

	cc.displayFinalStats()
	return nil
}

func (cc *CharacterCreation) displayWelcome() {
	cc.session.WriteLine("")
	cc.session.WriteLine("As you begin to concentrate not on who you are now but whom you would like")
	cc.session.WriteLine("to become, a number of questions begin to rattle about in your mind...")
	cc.session.WriteLine("")
}

func (cc *CharacterCreation) selectRace() error {
	for {
		cc.session.WriteLine("")
		cc.session.WriteLine("=== SELECT YOUR RACE ===")
		cc.session.WriteLine("")
		cc.session.WriteLine(races.MenuString())
		cc.session.WriteLine("")
		cc.session.WriteLine("If you would like to ruminate at greater length on the unique characteristics")
		cc.session.WriteLine("of one of these creatures, simply ask for [help]...")
		cc.session.WriteLine("")
		cc.session.Write("Choice (1-" + fmt.Sprintf("%d", races.Count()) + "): ")

		line, err := cc.session.ReadLine()
		if err != nil {
			return err
		}

		trimmed := strings.TrimSpace(line)
		if strings.ToLower(trimmed) == "help" {
			cc.session.WriteLine("Help is available for races (not fully implemented yet).")
			continue
		}

		choice, err := strconv.Atoi(trimmed)
		if err != nil {
			cc.session.WriteLine("Invalid input. Please enter a number.")
			continue
		}

		race := races.GetByMenuChoice(choice)
		if race == nil {
			cc.session.WriteLine("Invalid choice. Please try again.")
			continue
		}

		// Display flavor text and ask for confirmation
		if race.FlavorText != "" {
			cc.session.WriteLine("")
			cc.session.WriteLine(race.FlavorText)
		}
		cc.session.WriteLine("")
		cc.session.Write("Would you like to be " + race.Name + "? [Yes/No]: ")

		confirmLine, err := cc.session.ReadLine()
		if err != nil {
			return err
		}

		if strings.ToLower(strings.TrimSpace(confirmLine))[0:1] == "y" {
			cc.selectedRace = race
			cc.player.Race = race.RaceID
			cc.applyRaceModifiers()
			cc.session.WriteLine("")
			cc.session.WriteLine(fmt.Sprintf("Excellent! You have chosen to be a %s.", race.Name))
			return nil
		}

		cc.session.WriteLine("If you would not like to be one of those, what race would better suit you?")
	}
}

func (cc *CharacterCreation) selectStarterKit() error {
	for {
		cc.session.WriteLine("")
		cc.session.WriteLine("=== SELECT YOUR STARTING PATH ===")
		cc.session.WriteLine("")
		cc.session.WriteLine("In Wislyu, power comes from what you discover and who you learn from.")
		cc.session.WriteLine("Choose how you wish to begin your journey:")
		cc.session.WriteLine("")
		cc.session.WriteLine(kits.MenuString())
		cc.session.WriteLine("")
		cc.session.Write("Choice (1-3) or [help]: ")

		line, err := cc.session.ReadLine()
		if err != nil {
			return err
		}

		trimmed := strings.TrimSpace(line)
		if strings.ToLower(trimmed) == "help" {
			cc.session.WriteLine("")
			cc.session.WriteLine("=== SCHOLAR'S KIT ===")
			cc.session.WriteLine(kits.ScholarKit.FlavorText)
			cc.session.WriteLine("")
			cc.session.WriteLine("Starts with: Arcane Bolt spell (30%), Study skill (10%)")
			cc.session.WriteLine("")
			cc.session.WriteLine("=== WARRIOR'S KIT ===")
			cc.session.WriteLine(kits.WarriorKit.FlavorText)
			cc.session.WriteLine("")
			cc.session.WriteLine("Starts with: Slash maneuver (10%)")
			cc.session.WriteLine("")
			cc.session.WriteLine("=== WANDERER'S KIT ===")
			cc.session.WriteLine(kits.WandererKit.FlavorText)
			cc.session.WriteLine("")
			cc.session.WriteLine("Starts with: Arcane Bolt (20%), Study skill (5%), Slash maneuver (5%)")
			cc.session.WriteLine("")
			continue
		}

		choice, err := strconv.Atoi(trimmed)
		if err != nil {
			cc.session.WriteLine("Invalid input. Please enter a number.")
			continue
		}

		kit := kits.GetByMenuChoice(choice)
		if kit == nil {
			cc.session.WriteLine("Invalid choice. Please try again.")
			continue
		}

		// Display flavor text and ask for confirmation
		if kit.FlavorText != "" {
			cc.session.WriteLine("")
			cc.session.WriteLine(kit.FlavorText)
		}
		cc.session.WriteLine("")
		cc.session.Write("Would you like to begin with the " + kit.Name + "? [Yes/No]: ")

		confirmLine, err := cc.session.ReadLine()
		if err != nil {
			return err
		}

		if strings.ToLower(strings.TrimSpace(confirmLine))[0:1] == "y" {
			cc.selectedKit = kit
			cc.applyKitModifiers()
			cc.session.WriteLine("")
			cc.session.WriteLine(fmt.Sprintf("Excellent! You begin your journey with the %s.", kit.Name))
			return nil
		}

		cc.session.WriteLine("If that does not appeal to you, what path would better suit you?")
	}
}

func (cc *CharacterCreation) applyRaceModifiers() {
	if cc.selectedRace == nil {
		return
	}

	// Initialize base attributes to 10 if not set
	if cc.player.Attributes[0] == 0 {
		for i := range cc.player.Attributes {
			cc.player.Attributes[i] = 10
		}
	}

	// Apply race modifiers (STR, INT, WIS, DEX, CON, LCK, CHM)
	cc.player.Attributes[0] += cc.selectedRace.StrPlus // STR
	cc.player.Attributes[1] += cc.selectedRace.IntPlus // INT
	cc.player.Attributes[2] += cc.selectedRace.WisPlus // WIS
	cc.player.Attributes[3] += cc.selectedRace.DexPlus // DEX
	cc.player.Attributes[4] += cc.selectedRace.ConPlus // CON
	cc.player.Attributes[5] += cc.selectedRace.LckPlus // LCK
	cc.player.Attributes[6] += cc.selectedRace.ChaPlus // CHA

	// Apply racial stat bonuses
	cc.player.HP += cc.selectedRace.Hit
	cc.player.MaxHP += cc.selectedRace.Hit
	cc.player.Mana += cc.selectedRace.Mana
	cc.player.MaxMana += cc.selectedRace.Mana
	cc.player.Armor += cc.selectedRace.ACPlus
}

func (cc *CharacterCreation) applyKitModifiers() {
	if cc.selectedKit == nil {
		return
	}

	// Apply starter kit stats
	cc.player.HP = cc.selectedKit.HP
	cc.player.MaxHP = cc.selectedKit.HP
	cc.player.Mana = cc.selectedKit.Mana
	cc.player.MaxMana = cc.selectedKit.Mana
	cc.player.Move = cc.selectedKit.Move
	cc.player.MaxMove = cc.selectedKit.Move

	// Grant starting skills/spells
	if cc.player.Skills == nil {
		cc.player.Skills = make(map[int]*skills.PlayerSkillProgress)
	}

	for skillID, proficiency := range cc.selectedKit.StartingSkills {
		cc.player.Skills[skillID] = &skills.PlayerSkillProgress{
			SpellID:       skillID,
			Proficiency:   proficiency,
			Learned:       true,
			LifetimeCasts: 0,
		}
	}
}

func (cc *CharacterCreation) selectAge() error {
	for {
		cc.session.WriteLine("")
		cc.session.WriteLine("Next we must answer the question of your experience in the world. A more")
		cc.session.WriteLine("experienced person knows more of the world, has had more time to practice")
		cc.session.WriteLine("their talents and to learn of the world. A younger person, by contrast, may")
		cc.session.WriteLine("be full of youthful exuberance, but is unlikely to have dedicated the same")
		cc.session.WriteLine("amount of time to self improvement as their elder peers.")
		cc.session.WriteLine("")
		cc.session.WriteLine("[Child | Youth | Adult | Middle-Aged | Elderly]")
		cc.session.WriteLine("")
		cc.session.WriteLine("If you wish to know more about the specific differences between these age")
		cc.session.WriteLine("categories, feel free to ask for [help]. To that end, how old are you?")
		cc.session.WriteLine("")
		cc.session.Write("> ")

		line, err := cc.session.ReadLine()
		if err != nil {
			return err
		}

		trimmed := strings.TrimSpace(line)
		if strings.ToLower(trimmed) == "help" {
			cc.session.WriteLine("")
			for i, desc := range ageDescriptions {
				cc.session.WriteLine(fmt.Sprintf("%s: %s", ageNames[i], desc))
				cc.session.WriteLine("")
			}
			continue
		}

		ageChoice := -1
		firstChar := strings.ToLower(trimmed)[0:1]
		switch firstChar {
		case "c":
			ageChoice = AgeChild
		case "y":
			ageChoice = AgeYouth
		case "a":
			ageChoice = AgeAdult
		case "m":
			ageChoice = AgeMiddleAged
		case "e":
			ageChoice = AgeElderly
		default:
			cc.session.WriteLine("Please type [Child], [Youth], [Adult], [Middle-Aged], or [Elderly].")
			continue
		}

		// Display confirmation
		cc.session.WriteLine("")
		cc.session.WriteLine(ageDescriptions[ageChoice])
		cc.session.WriteLine("")
		cc.session.Write("Is a " + ageNames[ageChoice] + " acceptable? [Yes/No]: ")

		confirmLine, err := cc.session.ReadLine()
		if err != nil {
			return err
		}

		if strings.ToLower(strings.TrimSpace(confirmLine))[0:1] == "y" {
			cc.player.Age = ageChoice
			cc.session.WriteLine("")
			return nil
		}

		cc.session.WriteLine("Then please select another age category.")
	}
}

func (cc *CharacterCreation) selectSex() error {
	for {
		cc.session.WriteLine("")
		cc.session.WriteLine("Now we must address the matter of your physical form. In Njata, beings")
		cc.session.WriteLine("may manifest in many ways. You may choose to express yourself as:")
		cc.session.WriteLine("")
		cc.session.WriteLine("[Male | Female | Neutral]")
		cc.session.WriteLine("")
		cc.session.WriteLine("If you wish to understand these choices more deeply, feel free to ask for [help].")
		cc.session.WriteLine("So then, what form calls to you?")
		cc.session.WriteLine("")
		cc.session.Write("> ")

		line, err := cc.session.ReadLine()
		if err != nil {
			return err
		}

		trimmed := strings.TrimSpace(line)
		if strings.ToLower(trimmed) == "help" {
			cc.session.WriteLine("")
			for _, desc := range sexDescriptions {
				cc.session.WriteLine(desc)
				cc.session.WriteLine("")
			}
			continue
		}

		sexChoice := -1
		firstChar := strings.ToLower(trimmed)[0:1]
		switch firstChar {
		case "m":
			sexChoice = SexMale
		case "f":
			sexChoice = SexFemale
		case "n":
			sexChoice = SexNeutral
		default:
			cc.session.WriteLine("Please type [Male], [Female], or [Neutral].")
			continue
		}

		// Display confirmation
		cc.session.WriteLine("")
		cc.session.WriteLine(sexDescriptions[sexChoice])
		cc.session.WriteLine("")
		cc.session.Write("Does this resonate with your essence? [Yes/No]: ")

		confirmLine, err := cc.session.ReadLine()
		if err != nil {
			return err
		}

		if strings.ToLower(strings.TrimSpace(confirmLine))[0:1] == "y" {
			cc.player.Sex = sexChoice
			cc.session.WriteLine("")
			return nil
		}

		cc.session.WriteLine("Then perhaps one of the other forms would suit you better.")
	}
}

func (cc *CharacterCreation) displayFinalStats() {
	cc.session.WriteLine("")
	cc.session.WriteLine("=== FINAL CHARACTER STATS ===")
	cc.session.WriteLine("")
	cc.session.WriteLine(fmt.Sprintf("Name:        %s", cc.player.Name))
	cc.session.WriteLine(fmt.Sprintf("Race:        %s", cc.selectedRace.Name))
	cc.session.WriteLine(fmt.Sprintf("Starter Kit: %s", cc.selectedKit.Name))
	cc.session.WriteLine(fmt.Sprintf("Age:         %s", ageNames[cc.player.Age]))
	cc.session.WriteLine(fmt.Sprintf("Sex:         %s", sexNames[cc.player.Sex]))
	cc.session.WriteLine("")
	attrNames := []string{"STR", "INT", "WIS", "DEX", "CON", "LCK", "CHA"}
	for i, name := range attrNames {
		cc.session.WriteLine(fmt.Sprintf("%s: %d", name, cc.player.Attributes[i]))
	}

	cc.session.WriteLine("")
	cc.session.WriteLine(fmt.Sprintf("HP:   %d", cc.player.MaxHP))
	cc.session.WriteLine(fmt.Sprintf("Mana: %d", cc.player.MaxMana))
	cc.session.WriteLine("")
	cc.session.WriteLine("Welcome to Njata, adventurer.")
	cc.session.WriteLine("")
}
