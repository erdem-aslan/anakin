/**
 * Created by erdemaslan on 27/04/16.
 */
export class Service {

    constructor(public id:string,
                public serviceUrl:string,
                public endpoints:string[],
                public realEndpoints:Endpoint[],
                public balanceStrategy:string,
                public nested:boolean,
                public state:string,
                public editable?:boolean, 
                public bsIndex?:number) {
    }

}

export class Endpoint {

    
    constructor(public id:string,
                public host:string,
                public port:string,
                public scheme:string,
                public state:string) {

    }
}

