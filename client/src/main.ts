/// <reference path="sse.d.ts" />

import { IWorldState } from "./protocol";
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

      	for (let logEvent of update.logEvents) {
            let update = gameLog.scrollHeight - gameLog.clientHeight <= gameLog.scrollTop + 1
            gameLog.innerHTML = gameLog.innerHTML + logEvent
            if (update) {
              gameLog.scrollTop = gameLog.scrollHeight - gameLog.clientHeight;
            }
      	}

        renderer.pushState(update)
    })

    window.addEventListener("resize", () => renderer.onResize())
    window.requestAnimationFrame(function draw() {
        renderer.render()
        window.requestAnimationFrame(draw)
    })
})
