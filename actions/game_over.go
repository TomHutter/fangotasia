package actions

import (
	"fantasia/setup"
	"fantasia/view"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GameOver() {
	var board []string
	sum := 0
	board = append(board, fmt.Sprintln("G A M E    O V E R"))
	//board = append(board, fmt.Sprint(""))
	inv := setup.ObjectsInArea(setup.GetAreaByID(setup.INVENTORY))
	if len(inv) > 0 {
		board = append(board, fmt.Sprintln("Du besitzt:"))
		//board = append(board, fmt.Sprint(""))
	}
	for _, o := range inv {
		val := o.Properties.Value
		if val > 0 {
			sum += val
			desc := strings.Replace(o.Properties.Description.Long, "::", "", -1)
			board = append(board, fmt.Sprintf("- %s: %d Punkte", desc, val))
		}
	}
	board = append(board, fmt.Sprint(""))
	// all valuable objects found
	if sum == 170 {
		switch {
		case moves < 500:
			board = append(board, fmt.Sprintf("- Du hast weniger als %d Züge gebraucht: 7 Punkte", 500))
			sum += 7
			fallthrough
		case moves < 400:
			board = append(board, fmt.Sprintf("- Du hast weniger als %d Züge gebraucht: 7 Punkte", 400))
			sum += 7
			fallthrough
		case moves < 300:
			board = append(board, fmt.Sprintf("- Du hast weniger als %d Züge gebraucht: 7 Punkte", 300))
			sum += 7
		}
	}
	board = append(board, fmt.Sprint(""))
	if setup.GetAreaByID(5).Properties.Visited {
		board = append(board, fmt.Sprintf("- Du bist im Moor gewesen: %d Punkte", 2))
		sum += 2
	}
	if setup.GetAreaByID(29).Properties.Visited {
		board = append(board, fmt.Sprintf("- Du hast die verlassene Burg besucht: %d Punkte", 3))
		sum += 3
	}
	if setup.GetAreaByID(31).Properties.Visited {
		board = append(board, fmt.Sprintf("- Du bist auf einen Baum geklettert: %d Punkte", 4))
		sum += 4
	}
	board = append(board, fmt.Sprint(""))
	board = append(board, fmt.Sprintf("Du hast %d von 200 Punkten!", sum))
	board = append(board, fmt.Sprint("Noch ein Spiel (j/n)?"))
	view.PrintScreen(board)
	res := view.Scanner("once: true")
	if strings.ToLower(res) == "j" {
		setup.Setup()
	}
	/*
		621 poke214,9:poke211,13:sysvd:printb$"-rang ";
		622 ifpu=0thenprint"10 -":goto632
		623 ifpu<25thenprint"9 -":goto632
		624 ifpu<50thenprint"8 -":goto632
		625 ifpu<75thenprint"7 -":goto632
		626 ifpu<100thenprint"6 -":goto632
		627 ifpu<125thenprint"5 -":goto632
		628 ifpu<150thenprint"4 -":goto632
		629 ifpu<175thenprint"3 -":goto632
		630 ifpu<200thenprint"2 -"goto632
		631 ifpu=200thenprint"1 -"
	*/
	//return true
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	os.Exit(0)
}

//func (obj object) hit(inv []string) {
//}
