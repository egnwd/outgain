// server/protocol/protocol.go should be kept in sync with this

export interface IEntity {
    id: number
    name?: string
    color: string
    x: number
    y: number
    radius: number
}

export interface ILogEvent {
    logType: number
    protagID: number
    antagID: number
}

export interface IWorldState {
    time: number
    entities: IEntity[]
}
