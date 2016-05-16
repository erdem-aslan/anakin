/**
 * Created by erdemaslan on 18/04/16.
 */
import {Component, OnInit, ElementRef} from 'angular2/core';
import {ROUTER_DIRECTIVES, Router} from "angular2/router";
import {RouteConfig} from "angular2/router";
import {DashboardComponent} from "./dashboard.component";
import {ConfigurationComponent} from "./configuration.component";
import {StatisticsComponent} from "./statistics.component";
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
    {path: 'statistics', name: 'Statistics', component: StatisticsComponent}
])

export class AppComponent implements OnInit {


    constructor(private _router:Router,
                private _dom:ElementRef,
                private _anakinService:AnakinService) {
    }

    ngOnInit() {
    }

    onDashboardSelected() {
        this._router.navigate(['Dashboard']);
        console.log(this._router);
        this.toggleAnakinDrawer();
    }

    private toggleAnakinDrawer() {
        this._dom.nativeElement.querySelector("#anakin-drawer").togglePanel();
    }

    onConfigurationSelected() {
        this._router.navigate(['Configuration']);
        console.log(this._router);
        this.toggleAnakinDrawer();
    }

    onStatisticsSelected() {
        this._router.navigate(['Statistics']);
        console.log(this._router);
        this.toggleAnakinDrawer();
    }
}
