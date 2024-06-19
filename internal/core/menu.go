package core

import (
	"fmt"

	"github.com/mattn/go-tty"
)

type MenuItem struct {
	DisplayString string
	ReturnString  string
}

type Menu struct {
	Header        string
	MenuItems     []MenuItem
	selectedIndex int
}

// ---------- AddItem ----------
func (m *Menu) AddItem(displayString, returnString string) {
	m.MenuItems = append(m.MenuItems, MenuItem{DisplayString: displayString, ReturnString: returnString})
}

// ---------- Display ----------
func (m *Menu) Display() string {
	tty, err := tty.Open()
	if err != nil {
		panic(err)
	}
	defer tty.Close()

	defer func() {
		// show cursor again
		fmt.Printf("\033[?25h")
	}()

	// hide cursor
	fmt.Printf("\033[?25l")

	for {
		m.printMenu()

		r, err := tty.ReadRune()
		if err != nil {
			panic(err)
		}

		if r == 13 { // enter key
			fmt.Println()
			return m.MenuItems[m.selectedIndex].ReturnString
		} else if r == 27 { // escape sequence
			r1, _ := tty.ReadRune()
			r2, _ := tty.ReadRune()

			if r1 == 91 { // '['
				switch r2 {
				case 65: // up arrow
					m.moveUp()
				case 66: // down arrow
					m.moveDown()
				}
			}
		}
	}
}

// ---------- printMenu ----------
func (m *Menu) printMenu() {
	fmt.Println(textColor(m.Header, 72))

	for i, item := range m.MenuItems {
		if i == m.selectedIndex {
			fmt.Println(textColor(fmt.Sprintf("> %s", item.DisplayString), 172))
		} else {
			fmt.Printf("  %s\n", item.DisplayString)
		}
	}
}

// ---------- clearMenu ----------
func (m *Menu) clearMenu() {
	for i := 0; i < len(m.MenuItems)+1; i++ {
		fmt.Print("\033[1A") // move cursor up one line
		fmt.Print("\033[2K") // clear entire line
	}
}

// ---------- moveUp ----------
func (m *Menu) moveUp() {
	if m.selectedIndex > 0 {
		m.selectedIndex--
	}
	m.clearMenu()
}

// ---------- moveDown ----------
func (m *Menu) moveDown() {
	if m.selectedIndex < len(m.MenuItems)-1 {
		m.selectedIndex++
	}
	m.clearMenu()
}
