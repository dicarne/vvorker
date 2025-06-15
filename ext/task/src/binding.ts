export interface TaskBinding {
    client: () => Promise<TaskClient>;
}

export interface TaskClient {
    should_exit: () => Promise<boolean>;
    log: (text: string) => Promise<void>;
    complete: () => Promise<void>;
}