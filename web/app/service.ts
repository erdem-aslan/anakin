/**
 * Created by erdemaslan on 27/04/16.
 */

export class Service {
    
    constructor(private id:string, 
                private serviceUrl:string, 
                private endpoints:string[], 
                private balanceStrategy:string,
                private nested:boolean,
                private state:string) {}

}
