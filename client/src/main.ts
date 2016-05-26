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
      this.username = this.getUserID()
    }

    public setUserID() {
      $(this.usernameEl).html(this.username)
    }

    private getUserID() {
      var request = new XMLHttpRequest();
      request.open('GET', '/currentUser', false);
      request.send(null);

      let username = ""

      if (request.status == 200) {
        console.log("User: " + request.responseText);
        username = request.responseText
      }
      if (request.status == 401) {
        console.log("Not logged in");
        username = ""
      }

      return username
    }

    public isUserAuthenticated() {
      return this.username != ""
    }
}

$(function() {
    var userPanel = new UserPanel("#user-id", "#user-resources")

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
