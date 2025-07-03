export interface Node {
	Name: string
	UID: string
}

export interface PingMap {
	[key: string]: number
}

export interface GetNodeResponse {
	code: number
	msg: string
	data: {
		nodes: Node[]
		ping: PingMap
	}
}

export interface PingMapList {
	[key: string]: number[];
}