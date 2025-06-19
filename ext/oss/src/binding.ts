export interface OSSBinding {
    listBuckets(): Promise<any>;
    listObjects(bucket: string): Promise<any>;
    downloadFile(fileName: string): Promise<Uint8Array>;
    uploadFile(data: Uint8Array, fileName: string): Promise<any>;
    uploadStreamFile(stream: ReadableStream<Uint8Array>, fileName: string): Promise<any>;
    downloadStreamFile(fileName: string): Promise<ReadableStream<Uint8Array>>;
    deleteObject(fileName: string): Promise<any>;
}