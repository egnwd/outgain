/// <reference path="sse.d.ts" />
/// <reference path="extra.d.ts" />

import { IWorldState, ILogEvent } from "./protocol";
import { GameRenderer } from './renderer'
import { UserPanel, Timer, GameLog } from './gameUI'
import Editor from './editor'
import GameLeaderboard from './live-leaderboard'
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
  let re = /([0-9]+)\/$/g
  return re.exec(url)[1]
}

$(function() {
    var userPanel = new UserPanel("#user-id", "#user-gains-text")
    let timer = new Timer("#elapsed")
    let leaderboard = new GameLeaderboard("#game-leaderboard .table")
    let gameLog = new GameLog("game-log")


    let lobbyName = document.getElementById("lobby-name")
    $.ajax({ url: window.location.pathname + "/name", }).done((lobbyTitle) => {
      lobbyName.innerHTML = lobbyTitle
    })

    let idField = document.getElementById("id-field")
    let canvas = <HTMLCanvasElement> document.getElementById("game-view")
    let roundNamePopup = document.getElementById("round-name-popup")
    let roundNameSidebar = document.getElementById("round-name")

    let renderer = new GameRenderer(canvas, userPanel.username)

    let lobbyId = getLobbyId()
    idField.setAttribute("href", "/lobbies/" + lobbyId + "/summary")

    let source = new EventSource("/updates/" + lobbyId)

    source.addEventListener("state", function(event) {
        let data = JSON.parse((<sse.IOnMessageEvent>event).data)
        let update = <IWorldState>data

        roundNamePopup.style.display = "none"
        timer.pushState(update.progress, update.time)
        renderer.pushState(update)
    })

    source.addEventListener("round", function(event) {
      let message = JSON.parse((<sse.IOnMessageEvent>event).data)
      roundNamePopup.innerHTML = message
      roundNameSidebar.innerHTML = message
      roundNamePopup.style.display = "block"
      timer.reset()
    })

    source.addEventListener("gameover", function(event) {
      let message = "Game Over"
      roundNamePopup.innerHTML = message
      roundNameSidebar.innerHTML = message
      roundNamePopup.style.display = "block"
      setTimeout(() => {
        window.location.href = window.location.pathname + "/summary"
      }, 1500);
    })

    source.addEventListener("log", function(lEvent) {
	      let data = JSON.parse((<sse.IOnMessageEvent>lEvent).data)
        let logEvent = <ILogEvent>data

        gameLog.update(logEvent)

        if (userPanel.username == logEvent.protagName) {
          userPanel.updateScore(logEvent.gains)
        }

        leaderboard.refresh()
    })

    window.addEventListener("resize", () => renderer.onResize())
    window.requestAnimationFrame(function draw() {
        renderer.render()
        timer.render()
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
