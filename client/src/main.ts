/// <reference path="sse.d.ts" />

import { IWorldState } from "./protocol";
import { GameRenderer } from './renderer'
import { UserPanel } from './gameUI'
import * as $ from 'jquery'
import * as sweetalert from 'sweetalert';

var userPanel = new UserPanel("#user-id", "#user-resources")

$(function() {
    if (!userPanel.isUserAuthenticated()) {
      sweetalert({
          title: "<h1 id=\"title\"></h1>",
          text: "<a href=\"/login\" class=\"btn btn--action\">Login with Github</a>",
          html: true,
          showConfirmButton: false
      })
    } else {
      userPanel.setUserID()
    }
})

$(function() {
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
