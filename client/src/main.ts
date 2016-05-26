/// <reference path="sse.d.ts" />

import { IWorldState } from "./protocol";
import { GameRenderer } from './renderer'
import * as $ from 'jquery'
import * as sweetalert from 'sweetalert';

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

class UserPanel {
    username: string;
    resources: number;
    private usernameEl: string;
    private resourcesEl: string;

    constructor(usernameEl: string, resourcesEl: string) {
      this.usernameEl = usernameEl
      this.resourcesEl = resourcesEl
    }

    public setUserID() {
      let el = $(this.usernameEl)
      var username = ""

      $.ajax({url: "/currentUser", success: function(result){
        username = result
        el.html(result)
      }});

      this.username = username
    }
}

$(function() {
    function isUserAuthenticated() {
      let cookie = document.cookie.match(/^(.*; )?session=[^;]+(.*)?$/)
        return cookie != null
    }

    if (!isUserAuthenticated()) {
      sweetalert({
          title: "<h1 id=\"title\"></h1>",
          text: "<a href=\"/login\" class=\"btn btn--action\">Login with Github</a>",
          html: true,
          showConfirmButton: false
      })
    } else {
      var userPanel = new UserPanel("#user-id", "#user-resources")
      userPanel.setUserID()
    }
})
