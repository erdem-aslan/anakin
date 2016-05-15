/**
 * Created by erdemaslan on 19/04/16.
 */
import {Component, ElementRef, ViewChild, OnInit} from 'angular2/core';
import {AnakinService} from "./anakin.service"
import {MapToIterable} from "./mapToIterable"
import {SlashIfMissing} from "./slashIfMissing";
import {SearchFilterByName} from "./searchFilterByName";
import {SearchFilterById} from "./searchFilterById";
import {Application} from "./application";
import {NgSwitch, NgSwitchWhen, FORM_DIRECTIVES} from "angular2/common";
import {ApplicationComponent} from "./application.component";

@Component({
    selector: 'configuration',
    templateUrl: 'app/configuration.component.html',
    pipes: [MapToIterable, SlashIfMissing, SearchFilterByName, SearchFilterById],
    directives: [NgSwitch, NgSwitchWhen, FORM_DIRECTIVES, ApplicationComponent]
})

export class ConfigurationComponent implements OnInit {

    // @ViewChild(AppsComponent) private appsComponent:AppsComponent;

    newApp:Application;
    apps:Application[];
    error:string;
    appsError:string;
    appChangeError:string;

    constructor(private _anakinService:AnakinService, private _dom:ElementRef) {
        this.newApp = {id: '', name: '', baseUrl: '', services: {}, state: '', error: {}}
    }

    ngOnInit() {
        this.getApps();
    }

    private getApps() {
        this._anakinService.getApplications().subscribe(
            apps => this.apps = apps,
            error => this.getApplicationsError(error),
            () => this.getApplicationsFinished()
        )
    }

    onCreateNewApplication() {
        var state = this._dom.nativeElement.querySelector("#selected-state").selected;

        console.log(state);

        switch (state) {
            case 0:
                this.newApp.state = "active";
                break;
            case 1:
                this.newApp.state = "passive";
                break;
            default:
                console.log(state + " not handled.");
                break;
        }

        this.newApp.id = '';

        if (this.newApp.name == '') {
            return;
        }

        if (this.newApp.baseUrl == '') {
            this.newApp.baseUrl = '/';
        }

        this._anakinService
            .createApplication(this.newApp)
            .subscribe(app => this.newApp = app,
                error => this.createApplicationError(error),
                () => this.createApplicationFinished()
            );
    }

    createApplicationFinished() {
        console.log(this.newApp);
        this.getApps();

    }

    createApplicationError(error) {
        this.error = error;
        this._dom.nativeElement.querySelector("#createAppError").open();
    }

    getApplicationsError(error) {
        this.appsError = error;
        this._dom.nativeElement.querySelector("#getAppsError").open();
    }


    getApplicationsFinished() {
        console.log(this.apps);
    }

    onAppChangeError(event) {
        console.log(event);
        this.appChangeError = event.value;

        this._dom.nativeElement.querySelector("#appChangeError").open();
    }

    onAppChanged(event) {
        console.log(event);
        this.getApps();
    }
    
    
}
