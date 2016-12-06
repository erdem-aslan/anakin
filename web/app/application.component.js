System.register(['angular2/core', "./anakin.service", "./mapToIterable", "./slashIfMissing", "./dateFormatter", "./service", "./service.component"], function(exports_1, context_1) {
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
    var core_1, anakin_service_1, mapToIterable_1, slashIfMissing_1, dateFormatter_1, service_1, service_component_1;
    var ApplicationComponent;
    return {
        setters:[
            function (core_1_1) {
                core_1 = core_1_1;
            },
            function (anakin_service_1_1) {
                anakin_service_1 = anakin_service_1_1;
            },
            function (mapToIterable_1_1) {
                mapToIterable_1 = mapToIterable_1_1;
            },
            function (slashIfMissing_1_1) {
                slashIfMissing_1 = slashIfMissing_1_1;
            },
            function (dateFormatter_1_1) {
                dateFormatter_1 = dateFormatter_1_1;
            },
            function (service_1_1) {
                service_1 = service_1_1;
            },
            function (service_component_1_1) {
                service_component_1 = service_component_1_1;
            }],
        execute: function() {
            ApplicationComponent = (function () {
                function ApplicationComponent(_anakinService, _dom) {
                    this._anakinService = _anakinService;
                    this._dom = _dom;
                    this.editable = false;
                    this.errorEmitter = new core_1.EventEmitter(true);
                    this.mutationEmitter = new core_1.EventEmitter(true);
                    this.elevation = 4;
                    this.newService = new service_1.Service('', '', null, null, '', true, 'active');
                    this.bs = 0;
                }
                ApplicationComponent.prototype.ngOnInit = function () {
                    console.log("Fetching services...");
                    this.initializeSelected();
                    this.getServices();
                };
                ApplicationComponent.prototype.initializeSelected = function () {
                    if (this.app.state == "active") {
                        this.stateSelectedIndex = 0;
                    }
                    else if (this.app.state == "passive") {
                        this.stateSelectedIndex = 1;
                    }
                };
                ApplicationComponent.prototype.getServices = function () {
                    var _this = this;
                    this._anakinService.getServices(this.app.id).subscribe(function (services) { return _this.services = services; }, function (error) { return _this.handleServicesError(error); }, function () { return _this.servicesFetched(); });
                };
                ApplicationComponent.prototype.handleServicesError = function (error) {
                    this.errorEmitter.emit({ value: error });
                };
                ApplicationComponent.prototype.servicesFetched = function () {
                    // fetch endpoints
                };
                ApplicationComponent.prototype.onRemoveSelf = function () {
                    var _this = this;
                    this._anakinService.deleteApplication(this.app.id)
                        .subscribe(null, function (error) { return _this.deleteError(error); }, function () { return _this.selfDeleteCompleted(_this.app.id); });
                };
                ApplicationComponent.prototype.updateError = function (error) {
                    console.log(error.json());
                    this.errorEmitter.emit({ value: error });
                };
                ApplicationComponent.prototype.updateCompleted = function () {
                    // this.refresh(this.app.id)
                };
                ApplicationComponent.prototype.deleteError = function (error) {
                    console.log(error);
                    this.errorEmitter.emit({ value: error });
                };
                ApplicationComponent.prototype.selfDeleteCompleted = function (id) {
                    console.log("Self destruct completed, so long world...");
                    this.mutationEmitter.emit({ value: id });
                };
                ApplicationComponent.prototype.updateSelf = function () {
                    var _this = this;
                    var currentSelection = this._dom.nativeElement.querySelector("#selected-state").selected;
                    switch (currentSelection) {
                        case 0:
                            this.app.state = "active";
                            break;
                        case 1:
                            this.app.state = "passive";
                            break;
                    }
                    var updateApp = { id: this.app.id, baseUrl: this.app.baseUrl, state: this.app.state };
                    this._anakinService.updateApplication(this.app.id, updateApp)
                        .subscribe(null, function (error) { return _this.updateError(error); }, function () { return _this.updateCompleted(); });
                };
                ApplicationComponent.prototype.showRemoveSelfDialog = function () {
                    this._dom.nativeElement.querySelector("#delete-self").open();
                };
                ApplicationComponent.prototype.servicesPresent = function () {
                    return !(this.services == null || this.services.length == 0);
                };
                ApplicationComponent.prototype.serviceChangeError = function (event) {
                    this.errorEmitter.emit(event);
                };
                ApplicationComponent.prototype.serviceChanged = function (event, service) {
                    if (event != null) {
                        if (event.value = "deleted") {
                            console.log("Service has been deleted: " + service);
                            this.getServices();
                        }
                    }
                };
                ApplicationComponent.prototype.selectionChanged = function () {
                    this.stateSelectedIndex = this._dom.nativeElement.querySelector("#selected-state").selected;
                };
                ApplicationComponent.prototype.onCreateNewService = function () {
                    var _this = this;
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
                        .subscribe(function (service) { return _this.newService = service; }, function (error) { return _this.errorEmitter.emit({ value: error }); }, function () { return _this.getServices(); });
                };
                __decorate([
                    core_1.Input('editable'), 
                    __metadata('design:type', Boolean)
                ], ApplicationComponent.prototype, "editable", void 0);
                __decorate([
                    core_1.Input('app'), 
                    __metadata('design:type', Object)
                ], ApplicationComponent.prototype, "app", void 0);
                __decorate([
                    core_1.Output('changeError'), 
                    __metadata('design:type', core_1.EventEmitter)
                ], ApplicationComponent.prototype, "errorEmitter", void 0);
                __decorate([
                    core_1.Output('changed'), 
                    __metadata('design:type', core_1.EventEmitter)
                ], ApplicationComponent.prototype, "mutationEmitter", void 0);
                ApplicationComponent = __decorate([
                    core_1.Component({
                        selector: 'application',
                        pipes: [mapToIterable_1.MapToIterable, slashIfMissing_1.SlashIfMissing, dateFormatter_1.DateFormatter],
                        templateUrl: "app/application.component.html",
                        directives: [service_component_1.ServiceComponent]
                    }), 
                    __metadata('design:paramtypes', [anakin_service_1.AnakinService, core_1.ElementRef])
                ], ApplicationComponent);
                return ApplicationComponent;
            }());
            exports_1("ApplicationComponent", ApplicationComponent);
        }
    }
});
//# sourceMappingURL=application.component.js.map