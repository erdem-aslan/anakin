/**
 * Created by erdemaslan on 06/05/16.
 */
export class Instance {

    constructor(private id:string,
                private version:string,
                private adminPort:string,
                private adminIp:string,
                private proxyIp:string,
                private proxyPort:string,
                private started:Date,
                private state:string,
                private stats:InstanceStats) {
    }
}

export class InstanceStats {

    constructor(private os:string,
                private cpuCores:number,
                private mem:string,
                private rps:number) {
    }
}
