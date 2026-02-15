/// <reference path="../worker-configuration.d.ts" />

export interface EnvBinding {

	task: TaskBinding

	vars: {}

}
export interface TaskBinding {
    should_exit: () => Promise<boolean>;
    log: (text: string) => Promise<void>;
    complete: () => Promise<void>;
    create: (trace_id?: string) => Promise<string | undefined>;
}

