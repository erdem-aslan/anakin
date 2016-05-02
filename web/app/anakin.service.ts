/**
 * Created by erdemaslan on 19/04/16.
 */

import {Injectable, OnInit} from 'angular2/core';
import {Http, Response} from "angular2/http";
import {Service} from "./service";
import {Application} from "./application";
import {Endpoint} from "./endpoint";
import {Observable} from "rxjs/Observable";
import 'rxjs/Rx'


@Injectable()
export class AnakinService {

    private version:string;

    constructor(private http:Http) {
        this.version = "1.0"
    }

    getAnakinVersion() {
        return this.version;
    }

    getAnakinInstances() {
    }

    getApplications() {
        return this.http.get("/anakin/v1/apps")
            .map(res => <Application[]> res.json())
            .catch(this.handleError)
    }

    getServices(applicationId:string) {
    }


    getEndpoints(serviceId:string) {

    }

    private handleError(error:Response) {
        console.error(error);
        return Observable.throw(error.json().error || 'Server error')
    }
}
