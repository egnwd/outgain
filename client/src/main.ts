/// <reference path="jquery.d.ts" />
/// <reference path="sse.d.ts" />

import { IWorldUpdate } from "./protocol";
import { GameRenderer } from './renderer'
import * as $ from 'jquery'


$(function() {
    let canvas = <HTMLCanvasElement> document.getElementById("game-view")
    let renderer = new GameRenderer(canvas)

    let source = new EventSource("/updates")

    source.onmessage = function(event) {
        let data = JSON.parse(event.data)
        renderer.onUpdate(<IWorldUpdate>data)
    }
    window.addEventListener("resize", () => renderer.onResize())
    window.requestAnimationFrame(function draw() {
        renderer.render()
        window.requestAnimationFrame(draw)
    })

})
