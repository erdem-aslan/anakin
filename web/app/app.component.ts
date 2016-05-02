/**
 * Created by erdemaslan on 18/04/16.
 */
import {Component, OnInit, ElementRef} from 'angular2/core';
import {ROUTER_DIRECTIVES, Router} from "angular2/router";
import {RouteConfig} from "angular2/router";
import {DashboardComponent} from "./dashboard.component";
import {ConfigurationComponent} from "./configuration.component";
import {StatisticsComponent} from "./statistics.component";
import {MonitoringComponent} from "./monitoring.component";
import {AnakinService} from "./anakin.service";

@Component({
    selector: 'app',
    templateUrl: 'app/app.component.html',
    directives: [ROUTER_DIRECTIVES],
    providers:[AnakinService]
})

@RouteConfig([
    {path: 'dashboard', name: 'Dashboard', component: DashboardComponent, useAsDefault: true},
    {path: 'configuration', name: 'Configuration', component: ConfigurationComponent},
    {path: 'statistics', name: 'Statistics', component: StatisticsComponent},
    {path: 'monitoring', name: 'Monitoring', component: MonitoringComponent}
])

export class AppComponent implements OnInit {

    version:string;


    constructor(private _router:Router,
                private _dom:ElementRef,
                private _anakinService:AnakinService) {
    }

    ngOnInit() {
        console.log("App Component init");
        this.version = this._anakinService.getAnakinVersion();
    }

    onDashboardSelected() {
        this._router.navigate(['Dashboard']);
        this.toggleAnakinDrawer();
    }

    private toggleAnakinDrawer() {
        this._dom.nativeElement.querySelector("#anakin-drawer").togglePanel();
    }

    onConfigurationSelected() {
        this._router.navigate(['Configuration']);
        this.toggleAnakinDrawer();
    }

    onStatisticsSelected() {
        this._router.navigate(['Statistics']);
        this.toggleAnakinDrawer();
    }

    onMonitoringSelected() {
        this._router.navigate(['Monitoring']);
        this.toggleAnakinDrawer();
    }

}

export class AnakinInstance {
    id:string;
    uptime:number;
    host:string;
    clustered:boolean;
    clusterId:string;
}


export class Endpoint {
    id:string;
    host:string;
    port:number;
    scheme:string;
    state:string;
}
