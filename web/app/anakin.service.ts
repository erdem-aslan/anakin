/**
 * Created by erdemaslan on 19/04/16.
 */

import {Injectable} from 'angular2/core';
import {Http, Response, Headers} from "angular2/http";
import {Service, Endpoint} from "./service";
import {Application} from "./application";
import {Observable} from "rxjs/Observable";
import 'rxjs/Rx'
import {Instance} from "./instance";


@Injectable()
export class AnakinService {
    
    constructor(private http:Http) {
    }

    getAnakinInstances() {
        var headers:Headers = new Headers();
        headers.append("Accept", "application/json");

        return this.http.get("/anakin/v1/cluster", {
                headers: headers
            })
            .map(res => <Instance[]> res.json())
            .catch(this.handleError)
    }

    getApplications() {
        var headers:Headers = new Headers();
        headers.append("Accept", "application/json");

        return this.http.get("/anakin/v1/apps", {
                headers: headers
            })
            .map(res => <Application[]> res.json())
            .catch(this.handleError)
    }

    getApplication(appId:string) {
        var headers:Headers = new Headers();
        headers.append("Accept", "application/json");

        return this.http.get("/anakin/v1/apps/" + appId, {
                headers: headers
            })
            .map(res => <Application> res.json())
            .catch(this.handleError)
    }
    

    updateApplication(appId:string, body) {
        var headers:Headers = new Headers();
        headers.append("Accept", "application/json");
        headers.append("Content-Type", "application/json");

        return this.http.put("/anakin/v1/apps/" + appId, JSON.stringify(body), {
                headers: headers
            }).catch(this.handleError)
    }
    
    updateService(appId:string, serviceId:string,  body) {
        var headers:Headers = new Headers();
        headers.append("Accept", "application/json");
        headers.append("Content-Type", "application/json");

        return this.http.put("/anakin/v1/apps/" + appId + "/services/" + serviceId, JSON.stringify(body), {
                headers: headers
            }).catch(this.handleError)
    }

    createApplication(application:Application) {
        var headers:Headers = new Headers();
        headers.append("Accept", "application/json");
        headers.append("Content-Type", "application/json");

        return this.http.post("/anakin/v1/apps", JSON.stringify(application), {
            headers: headers
        }).map(res => <Application> res.json())
            .catch(this.handleError)
    }

    deleteApplication(id:string) {
        return this.http.delete("/anakin/v1/apps/" + id)
            .catch(this.handleError)
    }

    deleteService(appId:string, id:string) {
        return this.http.delete("/anakin/v1/apps/" + appId + "/services/" + id)
            .catch(this.handleError)
    }

    deleteEndpoint(appId:string, serviceId:string, id:string) {
        return this.http.delete("/anakin/v1/apps/" + appId + "/services/" + serviceId + "/endpoints/" + id)
            .catch(this.handleError)
    }


    getServices(applicationId:string) {
        return this.http.get("/anakin/v1/apps/" + applicationId + "/services")
            .map(res => <Service[]> res.json())
            .catch(this.handleError)
    }

    getService(applicationId:string, serviceId:string) {
        return this.http.get("/anakin/v1/apps/" + applicationId + "/services/" + serviceId)
            .map(res => <Service> res.json())
            .catch(this.handleError)
    }

    getEndpoints(applicationId:string, serviceId:string) {
        return this.http.get("/anakin/v1/apps/" + applicationId + "/services/" + serviceId + "/endpoints")
            .map(res => <Endpoint[]> res.json())
            .catch(this.handleError)

    }

    createService(applicationId:string, service:Service) {
        var headers:Headers = new Headers();
        headers.append("Accept", "application/json");
        headers.append("Content-Type", "application/json");

        return this.http.post("/anakin/v1/apps/" + applicationId + "/services", JSON.stringify(service), {
            headers:headers
        }).map(res => <Service> res.json())
            .catch(this.handleError)
    }

    createEndpoint(applicationId:string, serviceId:string, endpoint:Endpoint) {
        var headers:Headers = new Headers();
        headers.append("Accept", "application/json");
        headers.append("Content-Type", "application/json");

        return this.http.post("/anakin/v1/apps/" + applicationId + "/services/" + serviceId + "/endpoints", JSON.stringify(endpoint), {
            headers:headers
        }).map(res => <Endpoint> res.json())
            .catch(this.handleError)
    }


    private handleError(error:Response) {
        return Observable.throw(error.json().error || 'Server error')
    }
}
