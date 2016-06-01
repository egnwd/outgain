/// <reference path="sse.d.ts" />

import { IWorldState, ILogEvent } from "./protocol";
import { GameRenderer } from './renderer'
import { UserPanel } from './gameUI'
import * as $ from 'jquery'

var userPanel = new UserPanel("#user-id", "#user-resources")

// Move to GameUI
$(function() {
  $('#collapse-arrow').click(function() {
    // Change style
    $(this).toggleClass('small')
    $('#game-log').toggleClass('small')

    let swap = $(this).data("text-swap")
    $(this).data("text-swap", $(this).text());
    $(this).text(swap);

    let gameLog = document.getElementById("game-log")
    gameLog.scrollTop = gameLog.scrollHeight;

    return false
  })
})

function getLobbyId() {
  let url = window.location.href.toString()
  let re = /([0-9]+)$/g
  return url.match(re)[0]
}

$(function() {

    userPanel.setUserID()

    let gameLog = document.getElementById("game-log")
    let canvas = <HTMLCanvasElement> document.getElementById("game-view")

    let renderer = new GameRenderer(canvas)

    let source = new EventSource("/updates/" + getLobbyId())

    source.addEventListener("state", function(event) {
        let data = JSON.parse((<sse.IOnMessageEvent>event).data)

        let update = <IWorldState>data

        renderer.pushState(update)
    })

    source.addEventListener("shutdown", function() {
        window.location.href = "http://" + window.document.location.host + "/lobbies";
    })

    source.addEventListener("log", function(lEvent) {
	let data = JSON.parse((<sse.IOnMessageEvent>lEvent).data)
        let logEvent = <ILogEvent>data

        let scrollUpdate = gameLog.scrollHeight - gameLog.clientHeight <= gameLog.scrollTop + 1
	switch (logEvent.logType) {
	    case 0:
		gameLog.innerHTML = "A new game has started, good luck!\n"
	        break
	    case 1:
		gameLog.innerHTML = gameLog.innerHTML + "Yum, creature "
		    + logEvent.protagID + " ate a resource\n"
		break
	    case 2:
		gameLog.innerHTML = gameLog.innerHTML + "Creature "
		    + logEvent.protagID + " ate creature " + logEvent.antagID + "\n"
		break
	}
        if (scrollUpdate) {
            gameLog.scrollTop = gameLog.scrollHeight - gameLog.clientHeight
        }

    })

    window.addEventListener("resize", () => renderer.onResize())
    window.requestAnimationFrame(function draw() {
        renderer.render()
        window.requestAnimationFrame(draw)
    })
})
