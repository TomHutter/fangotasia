package movement

import (
	"fantasia/config"
	"fmt"
	"strings"
)

var visibleMap [12][10]int

func initVisibleAreas() {
	// set all areas to invisible
	for y := 0; y < 12; y++ {
		for x := 0; x < 10; x++ {
			visibleMap[y][x] = 0
		}
	}
	// show first area
	visibleMap[11][0] = 1
}

func AreaVisible(area int) bool {
	coordinates := config.AreaCoordinates[area]
	return visibleMap[coordinates.Y][coordinates.X] != 0
}

func drawBox(area int, boxLen int) (box [3]string) {
	// draw emty field, if area == 0
	if area == 0 {
		// boxlen + left an right connection
		spacer := strings.Repeat(" ", boxLen+2)
		for l := 0; l < 3; l++ {
			box[l] = fmt.Sprintf("%s", spacer)
		}
		return
	}
	// we have an overwrite for this box?
	if len(config.Overwrites) >= area && len(config.Overwrites[area][0]) > 0 {
		var dummy [3]string
		for i, v := range config.Overwrites[area] {
			dummy[i] = v
		}
		box = dummy
		return
	}
	var leftCon, rightCon, topCon, bottomCon string
	// get first line of area from locations
	text := strings.Split(config.Locations[area-1].Short, "\n")[0]
	textLen := len([]rune(text)) + 2 // two space left and right
	leftSpacer := strings.Repeat(" ", (boxLen-textLen)/2)
	rightSpacer := strings.Repeat(" ", boxLen-len(leftSpacer)-textLen)
	// horizontal line - left/right corner and middle connection element
	horLine := strings.Repeat(config.HL, (boxLen-3)/2)
	// can we walk to the north?
	if config.Areas[area][0] == 0 {
		// no => draw a hoizontal line
		topCon = config.HL
	} else {
		// yes => draw a connection to north
		topCon = config.TC
	}
	// can we walk to the south?
	if config.Areas[area][1] == 0 {
		// no => draw a hoizontal line
		bottomCon = config.HL
	} else {
		// yes => draw a connection to south
		bottomCon = config.BC
	}
	// can we walk to the east?
	if config.Areas[area][2] == 0 {
		// no => draw a vertical line
		rightCon = fmt.Sprintf("%s ", config.VL)
	} else {
		// yes => draw a connection to west
		rightCon = fmt.Sprintf("%s%s", config.RC, config.HL)
	}
	// can we walk to the west?
	if config.Areas[area][3] == 0 {
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

func DrawMap(area int) (text []string) {
	coordinates := config.AreaCoordinates[area]
	x := coordinates.X
	y := coordinates.Y
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
					v := visibleMap[iy][ix]
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
	coordinates := config.AreaCoordinates[area]
	visibleMap[coordinates.Y][coordinates.X] = area
	switch area {
	case 5:
		if AreaVisible(36) {
			visibleMap[11][4] = 57
			visibleMap[11][5] = 58
			visibleMap[5][5] = 59
		}
	case 6:
		if AreaVisible(7) {
			visibleMap[9][1] = 52
		}
	case 7:
		if AreaVisible(6) {
			visibleMap[9][1] = 52
		}
	case 9:
		if AreaVisible(31) {
			visibleMap[10][2] = 54
		}
	case 15:
		if AreaVisible(31) {
			visibleMap[9][2] = 55
		} else {
			visibleMap[9][2] = 53
		}
	case 31:
		visibleMap[10][2] = 54
		if AreaVisible(15) {
			visibleMap[9][2] = 55
		} else {
			visibleMap[9][2] = 56
		}
	case 32:
		if AreaVisible(37) {
			visibleMap[4][5] = 60
		} else {
			visibleMap[4][5] = 53
		}
		if visibleMap[11][5] != 0 {
			visibleMap[5][5] = 59
		} else {
			visibleMap[5][5] = 53
		}
	case 37:
		visibleMap[4][5] = 60
		if AreaVisible(40) {
			visibleMap[3][6] = 61
			visibleMap[4][6] = 62
		}
	case 38:
		if AreaVisible(40) {
			visibleMap[5][6] = 63
			visibleMap[6][6] = 64
		} else {
			visibleMap[5][6] = 0
		}
	case 39:
		if AreaVisible(40) {
			visibleMap[4][7] = 65
			visibleMap[5][7] = 64
		} else {
			visibleMap[4][7] = 0
		}
	case 40:
		visibleMap[3][6] = 61
		visibleMap[4][6] = 62
		visibleMap[5][6] = 63
		visibleMap[6][6] = 64
		visibleMap[4][7] = 65
		visibleMap[5][7] = 64
	case 41:
		if AreaVisible(51) {
			visibleMap[2][6] = 66
			if AreaVisible(40) {
				visibleMap[3][6] = 64
			} else {
				visibleMap[3][6] = 67
			}
		}
	case 42:
		if AreaVisible(51) {
			visibleMap[2][7] = 68
		} else {
			visibleMap[2][7] = 42
		}
	case 43:
		if AreaVisible(51) {
			visibleMap[1][6] = 69
		} else {
			visibleMap[1][6] = 0
		}
	case 44:
		if AreaVisible(51) {
			visibleMap[1][7] = 70
		} else {
			visibleMap[1][7] = 0
		}
	case 45:
		if AreaVisible(51) {
			visibleMap[1][8] = 71
			visibleMap[2][8] = 64
		} else {
			visibleMap[1][8] = 0
			visibleMap[2][8] = 0
		}
	case 46:
		if !AreaVisible(51) {
			visibleMap[1][9] = 0
		}
	case 47:
		if AreaVisible(51) {
			visibleMap[0][6] = 72
		} else {
			visibleMap[0][6] = 0
		}
	case 48:
		if AreaVisible(51) {
			visibleMap[0][7] = 73
		} else {
			visibleMap[0][7] = 0
		}
	case 49:
		if !AreaVisible(51) {
			visibleMap[0][8] = 0
		}
	case 50:
		if !AreaVisible(51) {
			visibleMap[0][9] = 0
		}
	case 51:
		visibleMap[2][6] = 66
		visibleMap[3][6] = 67
		visibleMap[2][7] = 68
		visibleMap[2][7] = 68
		visibleMap[1][6] = 69
		visibleMap[1][7] = 70
		visibleMap[1][8] = 71
		visibleMap[2][8] = 64
		visibleMap[0][6] = 72
		visibleMap[0][7] = 73
		visibleMap[1][9] = 46
		visibleMap[0][8] = 49
		visibleMap[0][9] = 50
	}
}
