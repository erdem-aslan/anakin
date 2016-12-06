System.register(['angular2/core', 'rxjs/Rx', "./anakin.service", "./slashIfMissing", "./service"], function(exports_1, context_1) {
    "use strict";
    var __moduleName = context_1 && context_1.id;
    var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
        var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
        if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
        else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
        return c > 3 && r && Object.defineProperty(target, key, r), r;
    };
    var __metadata = (this && this.__metadata) || function (k, v) {
        if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
    };
    var core_1, Rx_1, anakin_service_1, slashIfMissing_1, service_1;
    var ServiceComponent;
    return {
        setters:[
            function (core_1_1) {
                core_1 = core_1_1;
            },
            function (Rx_1_1) {
                Rx_1 = Rx_1_1;
            },
            function (anakin_service_1_1) {
                anakin_service_1 = anakin_service_1_1;
            },
            function (slashIfMissing_1_1) {
                slashIfMissing_1 = slashIfMissing_1_1;
            },
            function (service_1_1) {
                service_1 = service_1_1;
            }],
        execute: function() {
            ServiceComponent = (function () {
                function ServiceComponent(_anakinService, _dom) {
                    this._anakinService = _anakinService;
                    this._dom = _dom;
                    this.interval = 5;
                    this.autoRefresh = false;
                    this.editable = false;
                    this.mutationEmitter = new core_1.EventEmitter(true);
                    this.errorEmitter = new core_1.EventEmitter(true);
                    this.timerSubscription = null;
                    this.tempEditable = false;
                    // -- stub for creating new endpoints
                    this.newEndpoint = {
                        id: '',
                        host: '',
                        port: '',
                        scheme: 'http',
                        state: 'active'
                    };
                }
                ServiceComponent.prototype.ngOnInit = function () {
                    var _this = this;
                    this.resolveBsIndex();
                    this.fetchEndpoints();
                    if (this.autoRefresh) {
                        var timer = Rx_1.Observable.timer(this.interval * 1000, this.interval * 1000);
                        this.timerSubscription = timer.subscribe(function (t) {
                            _this.selfRefresh();
                        });
                    }
                };
                ServiceComponent.prototype.ngOnDestroy = function () {
                    if (this.timerSubscription != null) {
                        this.timerSubscription.unsubscribe();
                    }
                };
                ServiceComponent.prototype.resolveBsIndex = function () {
                    if (this.service.balanceStrategy == "round-robin") {
                        this.bsSelectedIndex = 0;
                    }
                    else if (this.service.balanceStrategy == "source-hashing") {
                        this.bsSelectedIndex = 1;
                    }
                };
                ServiceComponent.prototype.selfRefresh = function () {
                    var _this = this;
                    this._anakinService.getService(this.app.id, this.service.id)
                        .subscribe(function (service) { return _this.service = service; }, function (error) { return _this.errorEmitter.emit({ value: error }); }, function () { return _this.fetchEndpoints(); });
                };
                ServiceComponent.prototype.fetchEndpoints = function () {
                    var _this = this;
                    this._anakinService.getEndpoints(this.app.id, this.service.id)
                        .subscribe(function (endpoints) { return _this.addEndpoints(endpoints); }, function (error) { return _this.handleEndpointsError(error); }, function () { return _this.endpointsFetched(); });
                };
                ServiceComponent.prototype.addEndpoints = function (endpoints) {
                    if (this.service.realEndpoints == null) {
                        this.service.realEndpoints = [];
                    }
                    for (var _i = 0, endpoints_1 = endpoints; _i < endpoints_1.length; _i++) {
                        var endpoint = endpoints_1[_i];
                        this.service.realEndpoints.push(endpoint);
                    }
                };
                ServiceComponent.prototype.handleEndpointsError = function (error) {
                    this.errorEmitter.emit({ value: error });
                };
                ServiceComponent.prototype.endpointsFetched = function () {
                };
                ServiceComponent.prototype.onCreateNewEndpoint = function (appId, service) {
                    var _this = this;
                    this._anakinService.createEndpoint(appId, service.id, this.newEndpoint)
                        .subscribe(function (endpoint) { return service.realEndpoints.push(endpoint); }, function (error) { return _this.errorEmitter.emit({ value: error }); }, function () { return console.log("Endpoint has been added"); });
                };
                ServiceComponent.prototype.showRemoveSelfDialog = function () {
                    this._dom.nativeElement.querySelector("#delete-self").open();
                };
                ServiceComponent.prototype.onRemoveSelf = function () {
                    var _this = this;
                    this._anakinService.deleteService(this.app.id, this.service.id)
                        .subscribe(null, function (error) { return _this.errorEmitter.emit({ value: error }); }, function () { return _this.mutationEmitter.emit({ value: "deleted" }); });
                };
                ServiceComponent.prototype.balanceStrategyChanged = function () {
                    this.bsSelectedIndex =
                        this._dom.nativeElement.querySelector("#balance-strategy").selected;
                };
                ServiceComponent.prototype.editingFinished = function () {
                    var _this = this;
                    if (this.bsSelectedIndex == 0) {
                        this.service.balanceStrategy = "round-robin";
                    }
                    else if (this.bsSelectedIndex == 1) {
                        this.service.balanceStrategy = "source-hashing";
                    }
                    var body = {
                        id: this.service.id,
                        serviceUrl: this.service.serviceUrl,
                        balanceStrategy: this.service.balanceStrategy,
                        nested: this.service.nested
                    };
                    this._anakinService.updateService(this.app.id, this.service.id, body).subscribe(null, function (error) { return _this.errorEmitter.emit({ value: error }); }, function () { return _this.updateFinished(); });
                    this.tempEditable = !this.tempEditable;
                };
                ServiceComponent.prototype.updateFinished = function () {
                    console.log("Service has been updated");
                    this.selfRefresh();
                    this.mutationEmitter.emit(null);
                };
                ServiceComponent.prototype.nestedChanged = function () {
                    this.service.nested =
                        this._dom.nativeElement.querySelector("#nested").checked;
                };
                ServiceComponent.prototype.onRemoveEndpoint = function (endpointId) {
                    var _this = this;
                    this._anakinService.deleteEndpoint(this.app.id, this.service.id, endpointId).subscribe(null, function (error) { return _this.deleteError(error); }, function () { return _this.deleteEndpointCompleted(_this.service.id); });
                };
                ServiceComponent.prototype.deleteError = function (error) {
                    console.log(error);
                    this.errorEmitter.emit({ value: error });
                };
                ServiceComponent.prototype.deleteEndpointCompleted = function (id) {
                    this.service.realEndpoints = null;
                    this.fetchEndpoints();
                };
                ServiceComponent.prototype.showRemoveEndpointDialog = function () {
                    this._dom.nativeElement.querySelector("#delete-endpoint").open();
                };
                __decorate([
                    core_1.Input(), 
                    __metadata('design:type', Object)
                ], ServiceComponent.prototype, "app", void 0);
                __decorate([
                    core_1.Input(), 
                    __metadata('design:type', service_1.Service)
                ], ServiceComponent.prototype, "service", void 0);
                __decorate([
                    core_1.Input(), 
                    __metadata('design:type', Number)
                ], ServiceComponent.prototype, "interval", void 0);
                __decorate([
                    core_1.Input(), 
                    __metadata('design:type', Boolean)
                ], ServiceComponent.prototype, "autoRefresh", void 0);
                __decorate([
                    core_1.Input(), 
                    __metadata('design:type', Boolean)
                ], ServiceComponent.prototype, "editable", void 0);
                __decorate([
                    core_1.Output('changed'), 
                    __metadata('design:type', core_1.EventEmitter)
                ], ServiceComponent.prototype, "mutationEmitter", void 0);
                __decorate([
                    core_1.Output('changeError'), 
                    __metadata('design:type', core_1.EventEmitter)
                ], ServiceComponent.prototype, "errorEmitter", void 0);
                ServiceComponent = __decorate([
                    core_1.Component({
                        selector: 'service',
                        pipes: [slashIfMissing_1.SlashIfMissing],
                        templateUrl: 'app/service.component.html'
                    }), 
                    __metadata('design:paramtypes', [anakin_service_1.AnakinService, core_1.ElementRef])
                ], ServiceComponent);
                return ServiceComponent;
            }());
            exports_1("ServiceComponent", ServiceComponent);
        }
    }
});
//# sourceMappingURL=service.component.js.map