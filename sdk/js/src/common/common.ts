export function isDev() {
    return (import.meta as any).env.VITE_VVORKER_DEBUG === "true"
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