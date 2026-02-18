package game

import "strings"

// CapitalizeName capitalizes the first letter of a player name
// "vex" -> "Vex"
func CapitalizeName(name string) string {
	if len(name) == 0 {
		return name
	}
	return strings.ToUpper(name[:1]) + name[1:]
}
