/**
 * Created by erdemaslan on 27/04/16.
 */
export class Endpoint {

    constructor(private id:string,
                private host:string,
                private port:number, private scheme:string, private state:string) {
    }

}
