package api

import (
	"encoding/json"
	"fmt"
	"github.com/avast/retry-go"
	"log"
	"net/http"
	"time"

	"github.com/xyaren/gw2api"
)

var api = gw2api.NewGW2Api()

const language = "en"

var lastMatchId *string = nil

func GetData(worldId int) (gw2api.Match, map[int]string, gw2api.MatchStats) {
	var match *gw2api.Match
	if lastMatchId != nil {
		matchById := getMatchById(*lastMatchId)
		match = &matchById
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

func getMatchById(matchId string) gw2api.Match {
	var match gw2api.Match
	err := retry.Do(
		func() error {
			var err error
			var matches []gw2api.Match
			matches, err = api.MatchIds(matchId)

			if len(matches) < 1 {
				return fmt.Errorf("response without matches received for query by id %v: %v", matchId, matches)
			}
			match = matches[0]

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
	return match
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
	//err := retry.Do(
	//	func() error {
	//		var err error
	//		worlds, err = api.WorldIds(language, worldIds...)
	//		if err != nil {
	//			return err
	//		}
	//		return nil
	//	},
	//	retryLog,
	//)
	//if err != nil {
	//	panic(err)
	//}

	err := retry.Do(
		func() error {
			var err error
			var myClient = &http.Client{Timeout: 10 * time.Second}
			r, err := myClient.Get("https://next.werdes.net/json/temp_worlds.json")
			if err != nil {
				return err
			}
			defer r.Body.Close()

			var nextWorlds []gw2api.World
			err = json.NewDecoder(r.Body).Decode(&nextWorlds)
			if err != nil {
				return err
			}
			worlds = append(worlds, nextWorlds...)

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
