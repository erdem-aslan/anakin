/**
 * Created by erdemaslan on 18/04/16.
 */
import {Component, OnInit, ElementRef} from 'angular2/core';
import {AnakinService} from "./anakin.service"
import {Application} from "./application";
import {MapToIterable} from "./mapToIterable"
import {SlashIfMissing} from "./slashIfMissing";


@Component({
    selector: 'dashboard',
    templateUrl: 'app/dashboard.component.html',
    pipes: [MapToIterable, SlashIfMissing]

})

export class DashboardComponent implements OnInit {

    loadingApps:boolean = true;
    loadingServices:boolean = true;
    loadingEndpoints:boolean = true;

    appsElevation:number = 1;
    animatedShadow:boolean = true;

    apps:Application[];
    errorString:string;

    constructor(private _dom:ElementRef,
                private _anakinService:AnakinService) {

    }

    ngOnInit() {
        console.log("DashboardComponent  init");
        this.getApplications()
    }

    onHoverApplications() {
        this.appsElevation = 5;
    }

    onLeaveApplications() {
        this.appsElevation = 1;
    }

    getApplications() {
        this._anakinService.getApplications()
            .subscribe(
                apps => this.apps = apps,
                error => this.errorString = <any> error,
                () => this.getApplicationsCompleted()
            );
    }

    getApplicationsCompleted() {
        console.log("getApps finished");

        this.loadingApps = false;

        if (this.errorString) {
            console.error(this.errorString);
        } else {
            console.log(JSON.stringify(this.apps))
        }
    }

}