/// <reference path="sse.d.ts" />

import { IWorldState, ILogEvent } from "./protocol";
import { GameRenderer } from './renderer'
import { UserPanel, ModalPopUp } from './gameUI'
import * as $ from 'jquery'

var userPanel = new UserPanel("#user-id", "#user-resources")

$(function() {
    if (!userPanel.isUserAuthenticated()) {
      ModalPopUp.mainModal()
    } else {
      userPanel.setUserID()
    }
})

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

$(function() {
    // Do not render the game if the user is not authenticated
    if (!userPanel.isUserAuthenticated()) {
        return
    }

    let gameLog = document.getElementById("game-log")
    let canvas = <HTMLCanvasElement> document.getElementById("game-view")

    let renderer = new GameRenderer(canvas)

    let source = new EventSource("/updates")

    source.addEventListener("state", function(event) {
        let data = JSON.parse((<sse.IOnMessageEvent>event).data)

        let update = <IWorldState>data
	
        renderer.pushState(update)
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
