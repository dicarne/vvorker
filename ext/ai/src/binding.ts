export interface AIBinding {
    ask(messages: any[]): Promise<string>;
    chat(messages: any[]): Promise<any>;
    chatStream(messages: any[]): Promise<{
        next(): Promise<IteratorYieldResult<any> | IteratorReturnResult<any> | {
            done: boolean;
            err: any;
        }>;
    }>;
}