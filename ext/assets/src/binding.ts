export interface AssetsBinding {
    fetch(request: Request): Promise<Response>;
}