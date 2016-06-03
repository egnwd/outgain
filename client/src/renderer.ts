import { IEntity, IWorldState } from "./protocol";
import { lerp } from "./util";

class Entity {
    previous: IEntity
    current: IEntity
    img: HTMLImageElement

    constructor(state: IEntity) {
        this.previous = state
        this.current = state
        if (this.current.sprite != null) {
            let img = new Image()
            img.src = this.current.sprite
            img.onload = () => {
                this.img = img
            }
        }
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
        let radius = lerp(this.previous.radius, this.current.radius, interpolation)
        let color = this.current.color
        let name = this.current.name
        let entityType = this.current.entityType

	

        ctx.save()
        ctx.translate(x * scale, y * scale)

        if (entityType == 2){
	    drawStar(ctx, 0, 0, 12, radius * scale, radius /2 * scale)
	} else if (this.img != null) {
            let bumper = 1.4651162791
            var size = radius * scale * 2 * bumper
            ctx.drawImage(this.img, -size / 2, -size / 2, size, size)
        } else {
            ctx.beginPath()
            ctx.arc(0, 0, radius * scale, 0, 2 * Math.PI, false)
            ctx.fillStyle = color
            ctx.fill()
            ctx.closePath()
        }

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

function drawStar(ctx, cx, cy, spikes, outerRadius, innerRadius) {
    var rot = Math.PI / 2 * 3
    var x = cx
    var y = cy
    var step = Math.PI / spikes

    ctx.strokeSyle = "#000"
    ctx.beginPath()
    ctx.moveTo(cx, cy - outerRadius)
    for (let i = 0; i < spikes; i++) {
        x = cx + Math.cos(rot) * outerRadius
        y = cy + Math.sin(rot) * outerRadius
        ctx.lineTo(x, y)
        rot += step

        x = cx + Math.cos(rot) * innerRadius
        y = cy + Math.sin(rot) * innerRadius
        ctx.lineTo(x, y)
        rot += step
    }
    ctx.lineTo(cx, cy - outerRadius)
    ctx.closePath()
    ctx.lineWidth=5
    ctx.strokeStyle='red'
    ctx.stroke()
    ctx.fillStyle='black'
    ctx.fill()
}

export class GameRenderer {
    canvas: HTMLCanvasElement
    ctx: CanvasRenderingContext2D

    entities: { [key:number]: Entity } = {}

    previousTime: number
    currentTime: number
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
        if (this.previousTime != null) {
            let elapsed = Date.now() - this.lastUpdate
            return elapsed / this.dt
        } else {
            return 1
        }
    }

    render() {
        this.ctx.save()

        let interpolation = this.interpolation()

        let height = this.canvas.height
        let width = this.canvas.width

        this.ctx.clearRect(0, 0, width, height)

        let scale = Math.min(width / 10, height / 10)

        let xOffset = (width - 10 * scale) / 2
        this.ctx.translate(xOffset, 0)

        this.drawGrid(10, 10, scale)

        for (let id in this.entities) {
            this.entities[id].render(this.ctx, scale, interpolation)
        }

        this.ctx.restore()
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

    pushState(state: IWorldState) {
        let interpolation = this.interpolation()

        this.previousTime = this.currentTime
        this.currentTime = state.time

        let entities : { [key: number]: Entity } = {}

        for (let entityState of state.entities) {
            let entity = this.entities[entityState.id]
            if (typeof entity === "undefined") {
                entities[entityState.id] = new Entity(entityState)
            } else {
                entity.pushState(entityState, interpolation)
                entities[entityState.id] = entity
            }
        }

        this.entities = entities

        this.lastUpdate = Date.now()
        if (this.previousTime != null) {
            this.dt = this.currentTime - this.previousTime
        }
    }
}
