package api

import (
	"github.com/avast/retry-go"
	"github.com/xyaren/gw2api"
	"log"
)

var api = gw2api.NewGW2Api()

const language = "de"

var lastMatchId *string = nil

func GetData(worldId int) (gw2api.Match, map[int]string, gw2api.MatchStats) {
	var match *gw2api.Match
	if lastMatchId != nil {
		match = getMatchById(*lastMatchId)
	}
	if match == nil || !isInMatch(worldId, *match) {
		//log.Printf("Not in match: %s", match.ID)
		matches := getMatches()
		for _, elem := range matches {
			if isInMatch(worldId, elem) {
				match = &elem
				lastMatchId = &elem.ID
				break
			}
		}
	}

	if match == nil {
		panic("No match found for world id")
	}

	stats := getStats(match.ID)
	worldNameMap := getWorldNames(getIds(*match))
	return *match, worldNameMap, stats
}

func isInMatch(id int, match gw2api.Match) bool {
	return intInSlice(id, match.AllWorlds.Green) || intInSlice(id, match.AllWorlds.Blue) || intInSlice(id, match.AllWorlds.Red)
}

func intInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

var retryLog = retry.OnRetry(func(n uint, err error) {
	log.Printf("#%d: %s\n", n, err)
})

func getMatches() []gw2api.Match {
	var matchWorld []gw2api.Match
	err := retry.Do(
		func() error {
			var err error
			matchWorld, err = api.MatchIds("all")

			if err != nil && err.Error() != "Endpoint returned error: " {
				return err
			}
			return nil
		},
		retryLog,
	)
	if err != nil {
		panic(err)
	}
	return matchWorld
}

func getMatch(worldId int) gw2api.Match {
	var matchWorld gw2api.Match
	err := retry.Do(
		func() error {
			var err error
			matchWorld, err = api.MatchWorld(worldId)

			if err != nil && err.Error() != "Endpoint returned error: " {
				return err
			}
			return nil
		},
		retryLog,
	)
	if err != nil {
		panic(err)
	}
	return matchWorld
}

func getMatchById(matchId string) *gw2api.Match {
	var matchWorld gw2api.Match
	err := retry.Do(
		func() error {
			var err error
			var matchWorlds []gw2api.Match
			matchWorlds, err = api.MatchIds(matchId)
			matchWorld = matchWorlds[0]

			if err != nil && err.Error() != "Endpoint returned error: " {
				return err
			}
			return nil
		},
		retryLog,
	)
	if err != nil {
		panic(err)
	}
	return &matchWorld
}

func getStats(matchId string) gw2api.MatchStats {
	var stats gw2api.MatchStats
	err := retry.Do(
		func() error {
			var err error
			stats, err = api.MatchStats(matchId)
			if err != nil {
				return err
			}
			return nil
		},
		retryLog,
	)
	if err != nil {
		panic(err)
	}
	return stats
}

func getWorldNames(worldIds []int) map[int]string {
	var worlds []gw2api.World
	err := retry.Do(
		func() error {
			var err error
			worlds, err = api.WorldIds(language, worldIds...)
			if err != nil {
				return err
			}
			return nil
		},
		retryLog,
	)
	if err != nil {
		panic(err)
	}
	return toMap(worlds)
}

func getIds(matchWorld gw2api.Match) []int {
	var worlds []int
	worlds = append(worlds, matchWorld.AllWorlds.Green...)
	worlds = append(worlds, matchWorld.AllWorlds.Blue...)
	worlds = append(worlds, matchWorld.AllWorlds.Red...)
	return worlds
}

func toMap(ids []gw2api.World) map[int]string {
	elementMap := make(map[int]string)
	for _, entry := range ids {
		elementMap[entry.ID] = entry.Name
	}
	return elementMap
}
