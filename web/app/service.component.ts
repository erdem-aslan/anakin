/**
 * Created by erdemaslan on 13/05/16.
 */
import {
    Component, OnInit, OnDestroy, ElementRef,
    Input, Output, EventEmitter
} from 'angular2/core';

import {Observable} from 'rxjs/Rx';
import {AnakinService} from "./anakin.service"
import {SlashIfMissing} from "./slashIfMissing";
import {Application} from "./application";
import {Service, Endpoint} from "./service";
import {Subscription} from "rxjs/Subscription";


@Component({
    selector: 'service',
    pipes: [SlashIfMissing],
    templateUrl: 'app/service.component.html'
})


export class ServiceComponent implements OnInit,OnDestroy {

    @Input() app:Application;
    @Input() service:Service;
    @Input() interval:number = 5;
    @Input() autoRefresh:boolean = false;
    @Input() editable:boolean = false;

    @Output('changed') mutationEmitter:EventEmitter<any>
        = new EventEmitter(true);
    @Output('changeError') errorEmitter:EventEmitter<any>
        = new EventEmitter(true);

    private timerSubscription:Subscription = null;

    private tempEditable:boolean = false;

    private bsSelectedIndex:number;

    // -- stub for creating new endpoints
    private newEndpoint:Endpoint = {
        id: '',
        host: '',
        port: '',
        scheme: 'http',
        state: 'active'
    };


    constructor(private _anakinService:AnakinService,
                private _dom:ElementRef) {

    }

    ngOnInit() {

        this.resolveBsIndex();
        this.fetchEndpoints();

        if (this.autoRefresh) {
            let timer = Observable.timer(this.interval * 1000, this.interval * 1000);
            this.timerSubscription = timer.subscribe(t=> {
                this.selfRefresh();
            });
        }
    }

    ngOnDestroy() {
        if (this.timerSubscription != null) {
            this.timerSubscription.unsubscribe();
        }
    }

    resolveBsIndex() {
        if (this.service.balanceStrategy == "round-robin") {
            this.bsSelectedIndex = 0;
        } else if (this.service.balanceStrategy == "source-hashing") {
            this.bsSelectedIndex = 1;
        }
    }

    selfRefresh() {
        this._anakinService.getService(this.app.id, this.service.id)
            .subscribe(
                service => this.service = service,
                error => this.errorEmitter.emit({value: error}),
                () => this.fetchEndpoints()
            )
    }

    fetchEndpoints() {
        this._anakinService.getEndpoints(this.app.id, this.service.id)
            .subscribe(
                endpoints => this.addEndpoints(endpoints),
                error => this.handleEndpointsError(error),
                () => this.endpointsFetched()
            )
    }

    private addEndpoints(endpoints:Endpoint[]) {

        if (this.service.realEndpoints == null) {
            this.service.realEndpoints = [];
        }
        for (var endpoint of endpoints) {

            this.service.realEndpoints.push(endpoint);
        }
    }

    private handleEndpointsError(error) {
        this.errorEmitter.emit({value: error});
    }

    private endpointsFetched() {
    }

    onCreateNewEndpoint(appId:string, service:Service) {
        this._anakinService.createEndpoint(appId, service.id, this.newEndpoint)
            .subscribe(
                endpoint => service.realEndpoints.push(endpoint),
                error => this.errorEmitter.emit({value: error}),
                () => console.log("Endpoint has been added")
            )
    }

    showRemoveSelfDialog() {
        this._dom.nativeElement.querySelector("#delete-self").open();
    }

    onRemoveSelf() {
        this._anakinService.deleteService(this.app.id, this.service.id)
            .subscribe(
                null,
                error => this.errorEmitter.emit({value:error}),
                () => this.mutationEmitter.emit({value:"deleted"})
            )
    }

    balanceStrategyChanged() {
        this.bsSelectedIndex =
            this._dom.nativeElement.querySelector("#balance-strategy").selected;
    }

    editingFinished() {

        if (this.bsSelectedIndex == 0) {
            this.service.balanceStrategy = "round-robin";
        } else if (this.bsSelectedIndex == 1) {
            this.service.balanceStrategy = "source-hashing";
        }

        let body = {
            id: this.service.id,
            serviceUrl:this.service.serviceUrl,
            balanceStrategy: this.service.balanceStrategy,
            nested: this.service.nested
        };


        this._anakinService.updateService(this.app.id, this.service.id, body
        ).subscribe(
            null,
            error => this.errorEmitter.emit({value: error}),
            () => this.updateFinished()
        );

        this.tempEditable = !this.tempEditable;

    }

    updateFinished() {
        console.log("Service has been updated");
        this.selfRefresh();
        this.mutationEmitter.emit(null)
    }

    nestedChanged() {
        this.service.nested = 
            this._dom.nativeElement.querySelector("#nested").checked;
    }

    onRemoveEndpoint(endpointId:string) {
        this._anakinService.deleteEndpoint(this.app.id, this.service.id, endpointId).subscribe(
            null,
            error => this.deleteError(error),
            () => this.deleteEndpointCompleted(this.service.id)
        )
    }

    private deleteError(error) {
        console.log(error);
        this.errorEmitter.emit({value: error})
    }

    private deleteEndpointCompleted(id:string) {
        this.service.realEndpoints = null;
        this.fetchEndpoints();
    }

    showRemoveEndpointDialog() {
        this._dom.nativeElement.querySelector("#delete-endpoint").open();
    }





}