/// <reference path="jquery.d.ts" />
/// <reference path="sse.d.ts" />

import * as protocol from "./protocol";
import * as util from "./util";
import * as $ from 'jquery';

class Creature {
    previous: protocol.ICreature
    current: protocol.ICreature

    constructor(state: protocol.ICreature) {
        this.previous = state
        this.current = state
    }

    pushState(state: protocol.ICreature, interpolation: number) {
        this.current.x = util.lerp(this.previous.x, this.current.x, interpolation)
        this.current.y = util.lerp(this.previous.y, this.current.y, interpolation)

        this.previous = this.current
        this.current = state
    }

    render(ctx: CanvasRenderingContext2D, scale: number, interpolation: number) {
        let x = util.lerp(this.previous.x, this.current.x, interpolation)
        let y = util.lerp(this.previous.y, this.current.y, interpolation)
        let color = this.current.color
        let name = this.current.name

        ctx.save()
        ctx.translate(x * scale, y * scale)

        ctx.beginPath()
        ctx.arc(0, 0, scale / 2, 0, 2 * Math.PI, false)
        ctx.fillStyle = color
        ctx.fill()
        ctx.closePath()

        ctx.scale(2, 2)
        ctx.textAlign = "center"
        ctx.textBaseline = 'middle'
        ctx.fillStyle = "black"
        ctx.fillText(name, 0, 0)

        ctx.restore()
    }
}

class GameRenderer {
    canvas: HTMLCanvasElement
    ctx: CanvasRenderingContext2D
    creatures: { [key:number]: Creature }

    constructor(canvas: HTMLCanvasElement) {
        this.canvas = canvas
        this.ctx = canvas.getContext("2d")
        this.creatures = {}
    }

   onResize() {
        this.canvas.width = this.canvas.clientWidth
        this.canvas.height = this.canvas.clientHeight
    }

    render(interpolation: number) {
        let height = this.canvas.height
        let width = this.canvas.width

        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height)

        let scale = Math.min(width / 10, height / 10)

        this.drawGrid(10, 10, scale)

        for (let id in this.creatures) {
            let creature = this.creatures[id]
            creature.render(this.ctx, scale, interpolation)
        }
    }

    drawGrid(xsize, ysize, scale: number) {
        for (let x = 0; x < (xsize + 1) * scale; x += scale) {
            this.ctx.beginPath()
            this.ctx.moveTo(x, 0)
            this.ctx.lineTo(x, ysize * scale)
            this.ctx.stroke()
        }

        for (let y = 0; y < (ysize + 1) * scale; y += scale) {
            this.ctx.beginPath()
            this.ctx.moveTo(0, y)
            this.ctx.lineTo(xsize * scale, y)
            this.ctx.stroke()
        }
    }
}

$(function() {
    let canvas = <HTMLCanvasElement> document.getElementById("game-view")
    
    let gameLog = document.getElementById("game-log")


    let renderer = new GameRenderer(canvas)

    $(window).resize(function() {
        renderer.onResize()
    })

    renderer.onResize()

    let previous = null
    let current = null
    let dt = 0
    let lastUpdate = null

    let source = new EventSource("/updates")

    source.onmessage = function(event) {
        let data = JSON.parse(event.data)
        let update = <protocol.IWorldUpdate>data

	for (let logEvent of update.logEvents) {
	    gameLog.innerHTML = gameLog.innerHTML + logEvent
	}
        let interpolation = 1
        if (current != null && previous != null) {
            let elapsed = (Date.now() - lastUpdate)
            interpolation = elapsed / dt
        }

        previous = current
        current = update

        for (let state of update.creatures) {
            let creature = renderer.creatures[state.id]
            if (typeof creature === "undefined") {
                renderer.creatures[state.id] = new Creature(state)
            } else {
                creature.pushState(state, interpolation)
            }
        }

        lastUpdate = Date.now()
        if (previous != null) {
            dt = current.time - previous.time
        }
    }

    window.requestAnimationFrame(function draw() {
        if (previous == null) {
            renderer.render(1)
        } else {
            let elapsed = Date.now() - lastUpdate
            let interpolation = elapsed / dt

            renderer.render(interpolation)
        }

        window.requestAnimationFrame(draw)
    })
})

