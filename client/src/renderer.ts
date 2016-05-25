import { IEntity, IWorldUpdate } from "./protocol";
import { lerp } from "./util";

class Entity {
    previous: IEntity
    current: IEntity

    constructor(state: IEntity) {
        this.previous = state
        this.current = state
    }

    pushState(state: IEntity, interpolation: number) {
        this.current.x = lerp(this.previous.x, this.current.x, interpolation)
        this.current.y = lerp(this.previous.y, this.current.y, interpolation)

        this.previous = this.current
        this.current = state
    }

    render(ctx: CanvasRenderingContext2D, scale: number, interpolation: number) {
        let x = lerp(this.previous.x, this.current.x, interpolation)
        let y = lerp(this.previous.y, this.current.y, interpolation)
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

        if (name != null) {
            ctx.scale(2, 2)
            ctx.textAlign = "center"
            ctx.textBaseline = 'middle'
            ctx.fillStyle = "black"
            ctx.fillText(name, 0, 0)
        }

        ctx.restore()
    }
}

export class GameRenderer {
    canvas: HTMLCanvasElement
    ctx: CanvasRenderingContext2D

    entities: { [key:number]: Entity } = {}

    previous: IWorldUpdate
    current: IWorldUpdate
    dt: number
    lastUpdate: number

    constructor(canvas: HTMLCanvasElement) {
        this.canvas = canvas
        this.ctx = canvas.getContext("2d")

        this.onResize()
    }

    onResize() {
        this.canvas.width = this.canvas.clientWidth
        this.canvas.height = this.canvas.clientHeight
    }

    interpolation() : number {
        if (this.previous != null) {
            let elapsed = Date.now() - this.lastUpdate
            return elapsed / this.dt
        } else {
            return 1
        }
    }

    render() {
        let interpolation = this.interpolation()

        let height = this.canvas.height
        let width = this.canvas.width

        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height)

        let scale = Math.min(width / 10, height / 10)

        this.drawGrid(10, 10, scale)

        for (let id in this.entities) {
            this.entities[id].render(this.ctx, scale, interpolation)
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

    onUpdate(update: IWorldUpdate) {
        let interpolation = this.interpolation()

        this.previous = this.current
        this.current = update

        for (let state of update.entities) {
            let entity = this.entities[state.id]
            if (typeof entity === "undefined") {
                this.entities[state.id] = new Entity(state)
            } else {
                entity.pushState(state, interpolation)
            }
        }

        this.lastUpdate = Date.now()
        if (this.previous != null) {
            this.dt = this.current.time - this.previous.time
        }
    }
}
