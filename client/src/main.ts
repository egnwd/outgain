/// <reference path="sse.d.ts" />

import * as protocol from "./protocol";
import * as util from "./util";
import * as $ from 'jquery';
import * as sweetalert from 'sweetalert';

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
        let radius = this.current.radius
        let color = this.current.color
        let name = this.current.name

        ctx.save()
        ctx.translate(x * scale, y * scale)

        ctx.beginPath()
        ctx.arc(0, 0, radius * scale, 0, 2 * Math.PI, false)
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

class Resource {
    state: protocol.IResource

    constructor(state: protocol.IResource) {
        this.state = state
    }

    render(ctx: CanvasRenderingContext2D, scale: number) {
        ctx.save()
        ctx.translate(this.state.x * scale, this.state.y * scale)

        ctx.beginPath()
        ctx.arc(0, 0, this.state.radius * scale, 0, 2 * Math.PI, false)
        ctx.fillStyle = this.state.color
        ctx.fill()
        ctx.closePath()

        ctx.restore()
    }
}

class GameRenderer {
    canvas: HTMLCanvasElement
    ctx: CanvasRenderingContext2D
    creatures: { [key:number]: Creature } = {}
    resources: Resource[] = []

    constructor(canvas: HTMLCanvasElement) {
        this.canvas = canvas
        this.ctx = canvas.getContext("2d")
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

        for (let resource of this.resources) {
            resource.render(this.ctx, scale)
        }
        for (let id in this.creatures) {
            this.creatures[id].render(this.ctx, scale, interpolation)
        }
    }

    drawGrid(xsize, ysize, scale: number) {
        for (let x = 0; x <= xsize; x += 1) {
            this.ctx.beginPath()
            this.ctx.moveTo(x * scale, 0)
            this.ctx.lineTo(x * scale, ysize * scale)
            this.ctx.stroke()
        }

        for (let y = 0; y <= ysize; y += 1) {
            this.ctx.beginPath()
            this.ctx.moveTo(0, y * scale)
            this.ctx.lineTo(xsize * scale, y * scale)
            this.ctx.stroke()
        }
    }
}

$(function() {
    let canvas = <HTMLCanvasElement> document.getElementById("game-view")
    let renderer = new GameRenderer(canvas)

    let gameLog = document.getElementById("game-log")

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
            let update = gameLog.scrollHeight - gameLog.clientHeight <= gameLog.scrollTop + 1
            gameLog.innerHTML = gameLog.innerHTML + logEvent
            if (update) {
              gameLog.scrollTop = gameLog.scrollHeight - gameLog.clientHeight;
            }
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

        renderer.resources = update.resources.map(function(r) { return new Resource(r) })

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
      let cookie = document.cookie.match(/^(.*;)? session=[^;]+(.*)?$/)
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
