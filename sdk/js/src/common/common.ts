
let mode = "production"

export function init(env: any) {
    if (env.vars && env.vars.MODE) {
        mode = env.vars.MODE
    } else {
        mode = (import.meta as any).env.MODE
    }
}

export function isDev() {
    return mode === "development"
}


export function config() {
    let url = (import.meta as any).env.VITE_VVORKER_BASE_URL
    // Remove trailing slash if exists
    if (url.endsWith('/')) {
        url = url.slice(0, -1);
    }
    const token = (import.meta as any).env.VITE_VVORKER_TOKEN
    return {
        url,
        token
    }
}