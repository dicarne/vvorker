export interface WorkerItem {
  UID: string
  ExternalPath?: string
  HostName?: string
  NodeName?: string
  Port?: number
  Entry?: string
  Code: string
  Name: string
  Template: string
}

export interface WorkerItemProperties {
  item: WorkerItem
}

// name is generated on server side
// @ts-expect-error
export const DEFAUTL_WORKER_ITEM: WorkerItem = {
  UID: 'worker',
  Code: btoa(`export default {
  async fetch(req, env) {
    try {
		let resp = new Response("worker: " + req.url + " is online! -- " + new Date())
		return resp
	} catch(e) {
		return new Response(e.stack, { status: 500 })
	}
  }
};`),
  Template: `
{
    "name": "worker",
    "version": "1.0.0",
    "entensions": [],
    "services": [],
}
`
}

export interface WorkerEditorProperties {
  item: string
}

export interface VorkerSettingsProperties {
  WorkerURLSuffix: string
  Scheme: string
  EnableRegister: boolean
}
