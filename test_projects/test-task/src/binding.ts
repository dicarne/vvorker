/// <reference path="../worker-configuration.d.ts" />

export interface EnvBinding {

	task: TaskBinding

	vars: {
      MODE: string;
    }

}
export interface TaskBinding {
    should_exit: () => Promise<boolean>;
    log: (text: string) => Promise<void>;
    complete: () => Promise<void>;
    create: (name: string, trace_id?: string) => Promise<string | undefined>;
}

