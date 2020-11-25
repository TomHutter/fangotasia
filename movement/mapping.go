package movement

import (
	"fantasia/config"
	"fmt"
	"strings"
)

func drawBox(a int, boxLen int) (box [3]string) {
	// draw emty field, if area == 0
	if a == 0 {
		// boxlen + left an right connection
		spacer := strings.Repeat(" ", boxLen+2)
		for l := 0; l < 3; l++ {
			box[l] = fmt.Sprintf("%s", spacer)
		}
		return
	}
	// we have an overwrite for this box?
	ov := config.GetOverwriteByArea(a)
	if ov != (config.MapOverwrites{}) {
		var dummy [3]string
		for i, v := range ov.Content {
			dummy[i] = v
		}
		box = dummy
		return
	}
	var leftCon, rightCon, topCon, bottomCon string
	area := config.GetAreaByID(a)
	// get first line of area from locations
	text := area.Properties.Description.Short
	textLen := len([]rune(text)) + 2 // two space left and right
	leftSpacer := strings.Repeat(" ", (boxLen-textLen)/2)
	rightSpacer := strings.Repeat(" ", boxLen-len(leftSpacer)-textLen)
	// horizontal line - left/right corner and middle connection element
	horLine := strings.Repeat(config.HL, (boxLen-3)/2)
	// can we walk to the north?
	if area.Properties.Directions[0] == 0 {

		// no => draw a hoizontal line
		topCon = config.HL
	} else {
		// yes => draw a connection to north
		topCon = config.TC
	}
	// can we walk to the south?
	if area.Properties.Directions[1] == 0 {
		// no => draw a hoizontal line
		bottomCon = config.HL
	} else {
		// yes => draw a connection to south
		bottomCon = config.BC
	}
	// can we walk to the east?
	if area.Properties.Directions[2] == 0 {
		// no => draw a vertical line
		rightCon = fmt.Sprintf("%s ", config.VL)
	} else {
		// yes => draw a connection to west
		rightCon = fmt.Sprintf("%s%s", config.RC, config.HL)
	}
	// can we walk to the west?
	if area.Properties.Directions[3] == 0 {
		// no => draw a vertical line
		leftCon = fmt.Sprintf(" %s", config.VL)
	} else {
		// yes => draw a connection to west
		leftCon = fmt.Sprintf("%s%s", config.HL, config.LC)
	}
	box[0] = fmt.Sprintf(" %s%s%s%s%s ", config.BTL, horLine, topCon, horLine, config.BTR)
	box[1] = fmt.Sprintf("%s%s%s%s%s", leftCon, leftSpacer, text, rightSpacer, rightCon)
	box[2] = fmt.Sprintf(" %s%s%s%s%s ", config.BBL, horLine, bottomCon, horLine, config.BBR)
	return
}

func DrawMap(area config.Area) (text []string) {
	x := area.Properties.Coordinates.X
	y := area.Properties.Coordinates.Y
	// max x = 9, don't go further east than 8
	/*
		if x > 8 {
			x = 8
		}
	*/
	boxLen := config.BoxLen
	var boxes [5][3]string
	for i := 0; i < 6; i++ {
		iy := y + i - 2
		// outside y range => draw empty boxes
		if iy < 0 || iy > 11 {
			for j := 0; j < 5; j++ {
				boxes[j] = drawBox(0, boxLen)
			}
		} else {
			for j := 0; j < 5; j++ {
				ix := x + j - 2
				if ix < 0 || ix > 9 {
					boxes[j] = drawBox(0, boxLen)
				} else {
					v := config.Map[iy][ix]
					boxes[j] = drawBox(v, boxLen)
				}
			}
		}
		for l := 0; l < 3; l++ {
			if iy == y {
				text = append(text, fmt.Sprintf("%s%s%s%s%s%s%s%s", config.NEUTRAL, boxes[0][l], boxes[1][l],
					config.YELLOW, boxes[2][l],
					config.NEUTRAL, boxes[3][l], boxes[4][l]))
			} else {
				text = append(text, fmt.Sprintf("%s%s%s%s%s%s", config.NEUTRAL, boxes[0][l],
					boxes[1][l], boxes[2][l],
					boxes[3][l], boxes[4][l]))
			}
		}
	}
	//printScreen(text)
	return
}

func RevealArea(area int) {
	a := config.GetAreaByID(area)
	config.Map[a.Properties.Coordinates.Y][a.Properties.Coordinates.X] = area
	switch area {
	case 5:
		if config.AreaVisible(36) {
			config.Map[11][4] = 57
			config.Map[11][5] = 58
			config.Map[5][5] = 59
		}
	case 6:
		if config.AreaVisible(7) {
			config.Map[9][1] = 52
		}
	case 7:
		if config.AreaVisible(6) {
			config.Map[9][1] = 52
		}
	case 9:
		if config.AreaVisible(31) {
			config.Map[10][2] = 54
		}
	case 15:
		if config.AreaVisible(31) {
			config.Map[9][2] = 55
		} else {
			config.Map[9][2] = 53
		}
	case 31:
		config.Map[10][2] = 54
		if config.AreaVisible(15) {
			config.Map[9][2] = 55
		} else {
			config.Map[9][2] = 56
		}
	case 32:
		if config.AreaVisible(37) {
			config.Map[4][5] = 60
		} else {
			config.Map[4][5] = 53
		}
		if config.Map[11][5] != 0 {
			config.Map[5][5] = 59
		} else {
			config.Map[5][5] = 53
		}
	case 37:
		config.Map[4][5] = 60
		if config.AreaVisible(40) {
			config.Map[3][6] = 61
			config.Map[4][6] = 62
		}
	case 38:
		if config.AreaVisible(40) {
			config.Map[5][6] = 63
			config.Map[6][6] = 64
		} else {
			config.Map[5][6] = 0
		}
	case 39:
		if config.AreaVisible(40) {
			config.Map[4][7] = 65
			config.Map[5][7] = 64
		} else {
			config.Map[4][7] = 0
		}
	case 40:
		config.Map[3][6] = 61
		config.Map[4][6] = 62
		config.Map[5][6] = 63
		config.Map[6][6] = 64
		config.Map[4][7] = 65
		config.Map[5][7] = 64
	case 41:
		if config.AreaVisible(51) {
			config.Map[2][6] = 66
			if config.AreaVisible(40) {
				config.Map[3][6] = 64
			} else {
				config.Map[3][6] = 67
			}
		}
	case 42:
		if config.AreaVisible(51) {
			config.Map[2][7] = 68
		} else {
			config.Map[2][7] = 42
		}
	case 43:
		if config.AreaVisible(51) {
			config.Map[1][6] = 69
		} else {
			config.Map[1][6] = 0
		}
	case 44:
		if config.AreaVisible(51) {
			config.Map[1][7] = 70
		} else {
			config.Map[1][7] = 0
		}
	case 45:
		if config.AreaVisible(51) {
			config.Map[1][8] = 71
			config.Map[2][8] = 64
		} else {
			config.Map[1][8] = 0
			config.Map[2][8] = 0
		}
	case 46:
		if !config.AreaVisible(51) {
			config.Map[1][9] = 0
		}
	case 47:
		if config.AreaVisible(51) {
			config.Map[0][6] = 72
		} else {
			config.Map[0][6] = 0
		}
	case 48:
		if config.AreaVisible(51) {
			config.Map[0][7] = 73
		} else {
			config.Map[0][7] = 0
		}
	case 49:
		if !config.AreaVisible(51) {
			config.Map[0][8] = 0
		}
	case 50:
		if !config.AreaVisible(51) {
			config.Map[0][9] = 0
		}
	case 51:
		config.Map[2][6] = 66
		config.Map[3][6] = 67
		config.Map[2][7] = 68
		config.Map[2][7] = 68
		config.Map[1][6] = 69
		config.Map[1][7] = 70
		config.Map[1][8] = 71
		config.Map[2][8] = 64
		config.Map[0][6] = 72
		config.Map[0][7] = 73
		config.Map[1][9] = 46
		config.Map[0][8] = 49
		config.Map[0][9] = 50
	}
}
