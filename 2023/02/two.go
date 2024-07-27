package main

import (
	_ "embed"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	redMax   = 12
	greenMax = 13
	blueMax  = 14
)

//go:embed input.txt
var input string
var grabRe = regexp.MustCompile(`(?:(?P<count>\d+) (?P<color>red|green|blue),?){1,3}`)
var gameRe = regexp.MustCompile(`Game (?P<id>\d+):(?P<grabs>.*)`)

type Grab struct {
	red, blue, green int
}

type Game struct {
	id    int
	grabs []Grab
}

func (game *Game) isPossible() bool {
	for _, grab := range game.grabs {
		if grab.red > redMax || grab.blue > blueMax || grab.green > greenMax {
			return false
		}
	}
	return true
}

func (game *Game) minimumSet() (red, green, blue int) {
	for _, grab := range game.grabs {
		if grab.red > red {
			red = grab.red
		}
		if grab.blue > blue {
			blue = grab.blue
		}
		if grab.green > green {
			green = grab.green
		}
	}
	return
}

func (game *Game) minimumPowerSet() int {
	red, green, blue := game.minimumSet()
	return red * green * blue
}

func main() {
	games := parseGames(input)
	possibleGameIdSum := 0
	minimumPowerSetSum := 0
	for _, game := range games {
		fmt.Printf("%+v -> %v, minimum power set = %v\n", game, game.isPossible(), game.minimumPowerSet())
		if game.isPossible() {
			possibleGameIdSum += game.id
		}
		minimumPowerSetSum += game.minimumPowerSet()
	}
	fmt.Printf("Possible game id sum: %v\n", possibleGameIdSum)
	fmt.Printf("Minimum power set sum: %v\n", minimumPowerSetSum)
}

func parseGames(input string) []Game {
	lines := strings.Split(input, "\n")
	games := make([]Game, 0, len(lines))
	for _, line := range lines {
		game := parseGame(line)
		games = append(games, game)
	}
	return games
}

func parseGame(line string) Game {
	matches := gameRe.FindStringSubmatch(line)
	id, _ := strconv.Atoi(matches[gameRe.SubexpIndex("id")])
	grabs := parseGrabs(matches[gameRe.SubexpIndex("grabs")])
	return Game{
		id,
		grabs,
	}
}

func parseGrabs(line string) []Grab {
	grabs := make([]Grab, 0)
	for _, grab := range strings.Split(line, ";") {
		grabs = append(grabs, parseGrab(grab))
	}
	return grabs
}

func parseGrab(grabStr string) Grab {
	grab := grabRe.FindAllStringSubmatch(grabStr, -1)
	var red, blue, green int
	for _, oneColor := range grab {
		switch oneColor[grabRe.SubexpIndex("color")] {
		case "red":
			red, _ = strconv.Atoi(oneColor[grabRe.SubexpIndex("count")])
		case "blue":
			blue, _ = strconv.Atoi(oneColor[grabRe.SubexpIndex("count")])
		case "green":
			green, _ = strconv.Atoi(oneColor[grabRe.SubexpIndex("count")])
		}
	}
	return Grab{
		red,
		blue,
		green,
	}
}
