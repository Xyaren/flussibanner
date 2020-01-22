package api

import (
	"github.com/avast/retry-go"
	"github.com/xyaren/gw2api"
	"log"
)

var api = gw2api.NewGW2Api()

const language = "de"

func GetData(worldId int) (gw2api.Match, map[int]string, gw2api.MatchStats) {
	matchWorld := getMatch(worldId)
	stats := getStats(matchWorld.ID)
	worldNameMap := getWorldNames(getIds(matchWorld))
	return matchWorld, worldNameMap, stats
}

var retryLog = retry.OnRetry(func(n uint, err error) {
	log.Printf("#%d: %s\n", n, err)
})

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
