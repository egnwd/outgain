/// <reference path="sse.d.ts" />
/// <reference path="extra.d.ts" />

import { IWorldState, ILogEvent } from "./protocol";
import { GameRenderer } from './renderer'
import { UserPanel, Timer, GameLog } from './gameUI'
import Editor from './editor'
import * as $ from 'jquery'

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
    var userPanel = new UserPanel("#user-id", "user-gains-text")
    let timer = new Timer(15000, "#elapsed")

    let idField = document.getElementById("id-field")
    let gameLog = new GameLog("game-log")
    let canvas = <HTMLCanvasElement> document.getElementById("game-view")
    let roundName = document.getElementById("round-name")

    let renderer = new GameRenderer(canvas, userPanel.username)

    let lobbyId = getLobbyId()
    idField.setAttribute("value", lobbyId)

    let source = new EventSource("/updates/" + lobbyId)

    source.addEventListener("state", function(event) {
        let data = JSON.parse((<sse.IOnMessageEvent>event).data)
        let update = <IWorldState>data

        if (roundName.style.display == "block") {
            roundName.style.display = "none"
            console.log(update.progress)
            timer.start(update.progress)
        }

        renderer.pushState(update)
    })

    source.addEventListener("round", function(event) {
      roundName.innerHTML = JSON.parse((<sse.IOnMessageEvent>event).data)
      roundName.style.display = "block"
      timer.reset()
    })

    source.addEventListener("gameover", function(event) {
      window.location.href = "/lobbies"
    })

    source.addEventListener("log", function(lEvent) {
	      let data = JSON.parse((<sse.IOnMessageEvent>lEvent).data)
        let logEvent = <ILogEvent>data

        gameLog.update(logEvent)

        if (userPanel.username == logEvent.protagName) {
          userPanel.updateScore(logEvent.gains)
        }
    })

    window.addEventListener("resize", () => renderer.onResize())
    window.requestAnimationFrame(function draw() {
        renderer.render()
        window.requestAnimationFrame(draw)
    })

    $.ajax({
        url: "/token",
    }).done((token) => {
        let editor = new Editor(lobbyId, token);
        $('#edit-button').click(function() {
            editor.open()
        })
    })
})
