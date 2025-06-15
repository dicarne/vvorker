export interface OSSBinding {
    listBuckets(): Promise<any>;
    listObjects(bucket: string): Promise<any>;
    downloadFile(fileName: string): Promise<Uint8Array>;
    uploadFile(data: Uint8Array, fileName: string): Promise<any>;
    deleteObject(fileName: string): Promise<any>;
}