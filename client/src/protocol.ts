// server/protocol/protocol.go should be kept in sync with this

export interface ICreature {
    id: number
    name: string
    color: string
    x: number
    y: number
    radius: number
}

export interface IResource {
    color: string
    x: number
    y: number
    radius: number
}

export interface IWorldUpdate {
    time: number
    creatures: ICreature[]
    resources: IResource[]
}
