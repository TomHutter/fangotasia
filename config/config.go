package config

const (
	// https://en.wikipedia.org/wiki/ANSI_escape_code
	RED     = "\033[01;31m"
	GREEN   = "\033[01;32m"
	YELLOW  = "\033[01;33m"
	BLUE    = "\033[01;34m"
	WHITE   = "\033[01;97m"
	NEUTRAL = "\033[0m"
)

const (
	BTL = "\u250F"
	BTR = "\u2513"
	BBL = "\u2517"
	BBR = "\u251B"
	HL  = "\u2501"
	VL  = "\u2503"
	LC  = "\u252B"
	RC  = "\u2523"
	TC  = "\u253B"
	BC  = "\u2533"
	AR  = "\u2BC8"
	AL  = "\u2BC7"
	AT  = "\u2BC5"
	AB  = "\u2BC6"
)

var Areas = [52][4]int{
	{}, {0, 0, 2, 0}, {8, 0, 3, 1}, {9, 0, 0, 2}, {10, 0, 5, 0},
	{11, 0, 0, 4}, {13, 0, 7, 0}, {0, 8, 0, 0}, {7, 2, 9, 0},
	{15, 3, 10, 8}, {0, 4, 11, 9}, {0, 5, 12, 10}, {0, 0, 0, 11},
	{16, 6, 0, 0}, {17, 0, 15, 0}, {0, 9, 0, 14}, {0, 13, 17, 0},
	{0, 14, 18, 16}, {24, 0, 19, 17}, {25, 0, 20, 18}, {0, 0, 21, 19},
	{26, 0, 0, 20}, {27, 0, 23, 0}, {0, 0, 24, 22}, {29, 18, 0, 23},
	{30, 19, 0, 0}, {32, 21, 0, 0}, {0, 22, 28, 0}, {0, 0, 0, 27},
	{0, 0, 30, 0}, {0, 25, 0, 0}, {0, 0, 0, 0}, {0, 26, 0, 33},
	{0, 34, 32, 0}, {33, 36, 37, 35}, {0, 0, 34, 0}, {34, 0, 5, 0},
	{37, 38, 39, 34}, {37, 38, 38, 38}, {40, 39, 39, 37}, {42, 39, 0, 0},
	{43, 41, 42, 41}, {44, 40, 42, 41}, {47, 41, 43, 43}, {48, 42, 45, 44},
	{49, 45, 45, 44}, {50, 51, 0, 0}, {47, 43, 48, 47}, {48, 44, 48, 47},
	{0, 45, 50, 0}, {0, 46, 0, 49}, {46, 0, 0, 0},
}

var AreaMap = [12][10]int{
	{0, 0, 0, 0, 0, 0, 47, 48, 49, 50},
	{0, 0, 0, 0, 0, 0, 43, 44, 45, 46},
	{0, 0, 0, 0, 0, 0, 41, 42, 0, 51},
	{0, 0, 0, 0, 33, 32, 0, 40, 0, 0},
	{0, 0, 0, 35, 34, 0, 37, 39, 0, 0},
	{27, 28, 29, 30, 36, 0, 38, 0, 0, 0},
	{22, 23, 24, 25, 0, 26, 0, 0, 0, 0},
	{16, 17, 18, 19, 20, 21, 0, 0, 0},
	{13, 14, 15, 0, 0, 0, 0, 0, 0, 0},
	{6, 7, 52, 0, 0, 0, 0, 0, 0, 0},
	{0, 8, 9, 10, 11, 12, 0, 0, 0, 0},
	{1, 2, 3, 4, 5, 0, 0, 0, 0, 0},
}

/*
var MapOverwrite = [52]MapSpecials{
	{},{},{},{},{},
	{},{},{{"blah", "fahsel", "Blubb"},{},
	{11, 0, 0, 4}, {13, 0, 7, 0}, {0, 8, 0, 0}, {7, 2, 9, 0},
	{15, 3, 10, 8}, {0, 4, 11, 9}, {0, 5, 12, 10}, {0, 0, 0, 11},
	{16, 6, 0, 0}, {17, 0, 15, 0}, {0, 9, 0, 14}, {0, 13, 17, 0},
	{0, 14, 18, 16}, {24, 0, 19, 17}, {25, 0, 20, 18}, {0, 0, 21, 19},
	{26, 0, 0, 20}, {27, 0, 23, 0}, {0, 0, 24, 22}, {29, 18, 0, 23},
	{30, 19, 0, 0}, {32, 21, 0, 0}, {0, 22, 28, 0}, {0, 0, 0, 27},
	{0, 0, 30, 0}, {0, 25, 0, 0}, {0, 0, 0, 0}, {0, 26, 0, 33},
	{0, 34, 32, 0}, {33, 36, 37, 35}, {0, 0, 34, 0}, {34, 0, 5, 0},
	{37, 38, 39, 34}, {37, 38, 38, 38}, {40, 39, 39, 37}, {42, 39, 0, 0},
	{43, 41, 42, 41}, {44, 40, 42, 41}, {47, 41, 43, 43}, {48, 42, 45, 44},
	{49, 45, 45, 44}, {50, 51, 0, 0}, {47, 43, 48, 47}, {48, 44, 48, 47},
	{0, 45, 50, 0}, {0, 46, 0, 49}, {46, 0, 0, 0},
}

MapOverwrite[7] = {
*/

var ObjectsInArea = [45][2]int{
	{-1, 0}, {}, {}, {}, {}, {}, {}, {}, {},
	{28, 0}, {29, 0}, {8, 0}, {24, 0}, {2, 0}, {26, 0}, {20, 10}, {19, 0}, {19, 22},
	{18, 0}, {18, 20}, {16, 0}, {13, 0}, {14, 0}, {0, 0}, {0, 26}, {15, 5}, {8, 5},
	{9, 0}, {11, 0}, {12, 0}, {1, 0}, {1, 0}, {1, 0}, {3, 7}, {4, 0}, {4, 0},
	{4, 0}, {0, 0}, {0, 0}, {0, 18}, {30, 0}, {33, 10}, {51, 0}, {40, 0}, {51, 47},
}
