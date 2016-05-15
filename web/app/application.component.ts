/**
 * Created by erdemaslan on 11/05/16.
 */
import {Component, OnInit, ElementRef, Input, Output, EventEmitter} from 'angular2/core';
import {AnakinService} from "./anakin.service"
import {MapToIterable} from "./mapToIterable"
import {SlashIfMissing} from "./slashIfMissing";
import {DateFormatter} from "./dateFormatter";
import {Application} from "./application";
import {Service, Endpoint} from "./service";
import {ServiceComponent} from "./service.component";


@Component({
        selector: 'application',
        pipes: [MapToIterable, SlashIfMissing, DateFormatter],
        templateUrl: "app/application.component.html",
        directives:[ServiceComponent]
    }
)

export class ApplicationComponent implements OnInit {

    @Input('editable') editable:boolean = false;
    @Input('app') app:Application;

    @Output('changeError') errorEmitter:EventEmitter<any> = new EventEmitter(true);
    @Output('changed') mutationEmitter:EventEmitter<any> = new EventEmitter(true);

    elevation:number = 4;

    private services:Service[];

    private newService:Service = new Service(
        '','',null,null,'',true,'active'
    );

    private stateSelectedIndex:number;
    private bs:number = 0;




    constructor(private _anakinService:AnakinService,
                private _dom:ElementRef) {
    }

    ngOnInit() {

        console.log("Fetching services...");

        this.initializeSelected();
        this.getServices();

    }

    private initializeSelected() {
        if (this.app.state == "active") {
            this.stateSelectedIndex = 0;
        } else if (this.app.state == "passive") {
            this.stateSelectedIndex = 1;
        }

    }

    private getServices() {
        this._anakinService.getServices(this.app.id).subscribe(
            services => this.services = services,
            error => this.handleServicesError(error),
            () => this.servicesFetched()
        )

    }

    private handleServicesError(error) {
        this.errorEmitter.emit({value: error});
    }

    private servicesFetched() {
        // fetch endpoints
    }

    onRemoveSelf() {
        this._anakinService.deleteApplication(this.app.id)
            .subscribe(
                null,
                error => this.deleteError(error),
                () => this.selfDeleteCompleted(this.app.id)
            );
    }

    updateError(error) {
        console.log(error.json());
        this.errorEmitter.emit({value: error})
    }

    updateCompleted() {
        // this.refresh(this.app.id)
    }

    private deleteError(error) {
        console.log(error);
        this.errorEmitter.emit({value: error})
    }

    private selfDeleteCompleted(id:string) {
        console.log("Self destruct completed, so long world...");
        this.mutationEmitter.emit({value: id})
    }

    updateSelf() {
        var currentSelection:number = this._dom.nativeElement.querySelector("#selected-state").selected;

        switch (currentSelection) {
            case 0:
                this.app.state = "active";
                break;
            case 1:
                this.app.state = "passive";
                break;
        }

        let updateApp = {id: this.app.id, baseUrl:this.app.baseUrl, state: this.app.state};

        this._anakinService.updateApplication(this.app.id, updateApp)
            .subscribe(
                null,
                error => this.updateError(error),
                () => this.updateCompleted()
            );


    }


    showRemoveSelfDialog() {
        this._dom.nativeElement.querySelector("#delete-self").open();
    }

    servicesPresent() {
        return !(this.services == null || this.services.length == 0);
    }

    serviceChangeError(event) {
        this.errorEmitter.emit(event);
    }
    
    serviceChanged(event, service) {
        if (event != null) {
            if (event.value = "deleted") {
                console.log("Service has been deleted: " + service);
                this.getServices();
            }
        }
    }

    selectionChanged() {
        this.stateSelectedIndex = this._dom.nativeElement.querySelector("#selected-state").selected;
    }

    private onCreateNewService() {

        switch (this.bs) {
            case 0:
                this.newService.balanceStrategy = "round-robin";
                break;
            case 1:
                this.newService.balanceStrategy = "source-hashing";
                break;
            default:
                console.log(this.bs);
                break;
        }

        this.newService.id = '';

        if (this.newService.serviceUrl == '') {
            this.newService.serviceUrl = '/';
        }

        this._anakinService.createService(this.app.id, this.newService)
            .subscribe(service => this.newService = service,
                error => this.errorEmitter.emit({value: error}),
                () => this.getServices()
            );

    }

    

}
