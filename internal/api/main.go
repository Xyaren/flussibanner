package api

import (
	"github.com/xyaren/gw2api"
)

var api = gw2api.NewGW2Api()

const language = "de"

func GetData(worldId int) (gw2api.Match, map[int]string, gw2api.MatchStats) {
	matchWorld, _ := api.MatchWorld(worldId)
	stats, _ := api.MatchWorldStats(worldId)
	worldNameMap := getWorldMap(matchWorld)
	return matchWorld, worldNameMap, stats
}

func getWorldMap(matchWorld gw2api.Match) map[int]string {
	var worlds []int
	worlds = append(worlds, matchWorld.AllWorlds.Green...)
	worlds = append(worlds, matchWorld.AllWorlds.Blue...)
	worlds = append(worlds, matchWorld.AllWorlds.Red...)
	ids, _ := api.WorldIds(language, worlds...)
	return toMap(ids)
}

func toMap(ids []gw2api.World) map[int]string {
	elementMap := make(map[int]string)
	for _, entry := range ids {
		elementMap[entry.ID] = entry.Name
	}
	return elementMap
}
