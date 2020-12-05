package intro

import "fantasia/view"

func Intro() {
	var text = []string{
		"\033[01;34mMach Dich auf den gefahrenreichen Weg in",
		"das zauberhafte Land Fangotasia und suche",
		"nach märchenhaften Schätzen.",
		"Führe mich mit einfachen Kommandos in",
		"einem oder zwei Worten, z.B.:",
		"",
		"\033[01;33mNORDEN      BENUTZE TARNKAPPE      ENDE",
		"",
		"LEGE RUBIN     TÖTE DRACHE     INVENTAR",
		"",
		"\033[01;34mMit  \033[01;97mSAVE  \033[01;34mkannst Du den aktuellen Stand",
		"des Spieles abspeichern,",
		"mit  \033[01;97mLOAD  \033[01;34mwieder einlesen.",
		"\033[0mWeiter \u23CE",
	}
	view.PrintScreen(text)
	view.Scanner("once: true")
}

func Prelude() {
	var text = []string{
		"\033[01;31mF A N G O T A S I A",
		"",
		"\033[01;34m- Ein Adventure von Klaus Hartmuth -",
		"",
		"\033[01;33m- GO Version von Tom Hutter -",
		"",
		"\033[0mWeiter \u23CE",
	}
	view.PrintScreen(text)
	view.Scanner("once: true")
}
