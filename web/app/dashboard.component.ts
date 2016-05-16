/**
 * Created by erdemaslan on 18/04/16.
 */
import {Component, OnInit, ElementRef, OnDestroy} from 'angular2/core';
import {NumberPipe, DecimalPipe} from 'angular2/common';
import {Observable} from 'rxjs/Rx';
import {AnakinService} from "./anakin.service"
import {ApplicationComponent} from "./application.component";
import {Application} from "./application";
import {MapToIterable} from "./mapToIterable"
import {SlashIfMissing} from "./slashIfMissing";
import {Instance} from "./instance";
import {DateFormatter} from "./dateFormatter";
import {SearchFilterByName} from "./searchFilterByName";
import {SearchFilterById} from "./searchFilterById";
import {Subscription} from "rxjs/Subscription";


@Component({
    selector: 'dashboard',
    templateUrl: 'app/dashboard.component.html',
    directives:[ApplicationComponent],
    pipes: [MapToIterable,
        SlashIfMissing,
        DateFormatter,
        SearchFilterByName,
        SearchFilterById, DecimalPipe]
})

export class DashboardComponent implements OnInit, OnDestroy {

    timerSubscription:Subscription = null;

    loadingInstances:boolean = true;
    loadingApps:boolean = true;


    instancesElevation:number = 5;
    animatedShadow:boolean = true;

    instances:Instance[];
    instancesError:string;
    appsError:string;

    apps:Application[];

    selectedTab:number = 0;
    selectedApp:number = 0;


    constructor(private _dom:ElementRef,
                private _anakinService:AnakinService) {

    }

    ngOnInit() {
        console.log("DashboardComponent  init");

        this.getApps();
        this.getInstances();

        let timer = Observable.timer(1000, 1000);
        this.timerSubscription = timer.subscribe(t=> {
            this.getInstances();
        });

    }

    ngOnDestroy() {

        if (this.timerSubscription != null) {
            this.timerSubscription.unsubscribe();
        }

    }



    getInstances() {
        this._anakinService.getAnakinInstances()
            .subscribe(
                instances => this.instances = instances,
                error => this.instancesError = <any> error,
                () => this.getInstancesCompleted()
            );
    }

    getInstancesCompleted() {
        this.loadingInstances = false;

        if (this.instancesError) {
            console.error(this.instancesError);
        }
    }

    getApps() {
        this._anakinService.getApplications().subscribe(
            apps => this.apps = apps,
            error => this.appsError = <any> error,
            () => this.getAppsCompleted()
        );
    }

    getAppsCompleted() {
        this.loadingApps = false;
    }

    onDashboardTabSelected(event) {
        this.selectedTab = this._dom.nativeElement.querySelector("#dashboard-tabs").selected;
    }

}