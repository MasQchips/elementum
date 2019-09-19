package api

import (
	"fmt"

	"github.com/elgatito/elementum/database"
	"github.com/elgatito/elementum/xbmc"

	"github.com/gin-gonic/gin"
)

// History ...
func History(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	infohash := ctx.Query("infohash")
	if torrent := InTorrentsHistory(infohash); torrent != nil {
		xbmc.PlayURLWithTimeout(URLQuery(
			URLForXBMC("/play"), "uri", torrent.URI,
		))
		return
	}

	items := []*xbmc.ListItem{}
	var ths []database.TorrentHistory
	if err := database.GetStormDB().All(&ths); err != nil {
		log.Infof("Could not get list of history items: %s", err)
	}
	if err := database.GetStormDB().AllByIndex("Dt", &ths); err != nil {
		log.Infof("Could not get list of history items: %s", err)
	}

	for _, th := range ths {
		items = append(items, &xbmc.ListItem{
			Label: th.Name,
			Path:  torrentHistoryGetXbmcURL(th.InfoHash),
			ContextMenu: [][]string{
				[]string{"LOCALIZE[30406]", fmt.Sprintf("XBMC.RunPlugin(%s)",
					URLQuery(URLForXBMC("/history/remove"),
						"infohash", th.InfoHash,
					))},
			},
		})
	}

	ctx.JSON(200, xbmc.NewView("", items))
}

func torrentHistoryEmpty() bool {
	count, err := database.GetStormDB().Count(&database.TorrentHistory{})
	if err != nil {
		log.Infof("Could not get count for torrent history: %s", err)
	}

	return err != nil || count == 0
}

// HistoryRemove ...
func HistoryRemove(ctx *gin.Context) {
	infohash := ctx.DefaultQuery("infohash", "")

	if len(infohash) == 0 {
		return
	}

	log.Debugf("Removing infohash '%s' with torrent history", infohash)
	database.GetStormDB().Delete(database.TorrentHistoryBucket, infohash)
	xbmc.Refresh()

	ctx.String(200, "")
	return
}

// HistoryClear ...
func HistoryClear(ctx *gin.Context) {
	log.Debugf("Cleaning queries with torrent history")
	if err := database.GetStormDB().Drop(&database.TorrentHistory{}); err != nil {
		log.Infof("Could not clean torrent history: %s", err)
	}
	xbmc.Refresh()

	ctx.String(200, "")
	return
}

func torrentHistoryGetXbmcURL(infohash string) string {
	return URLQuery(URLForXBMC("/history"), "infohash", infohash)
}

func torrentHistoryGetHTTPUrl(infohash string) string {
	return URLQuery(URLForHTTP("/history"), "infohash", infohash)
}
