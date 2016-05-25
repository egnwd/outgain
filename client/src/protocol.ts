export interface ICreature {
    id: number
    name: string
    color: string
    x: number
    y: number
}

export interface IWorldUpdate {
    time: number
    creatures: ICreature[]
    logEvents : string[]
}
