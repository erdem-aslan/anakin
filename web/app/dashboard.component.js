System.register(['angular2/core', 'angular2/common', 'rxjs/Rx', "./anakin.service", "./application.component", "./mapToIterable", "./slashIfMissing", "./dateFormatter", "./searchFilterByName", "./searchFilterById"], function(exports_1, context_1) {
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
    var core_1, common_1, Rx_1, anakin_service_1, application_component_1, mapToIterable_1, slashIfMissing_1, dateFormatter_1, searchFilterByName_1, searchFilterById_1;
    var DashboardComponent;
    return {
        setters:[
            function (core_1_1) {
                core_1 = core_1_1;
            },
            function (common_1_1) {
                common_1 = common_1_1;
            },
            function (Rx_1_1) {
                Rx_1 = Rx_1_1;
            },
            function (anakin_service_1_1) {
                anakin_service_1 = anakin_service_1_1;
            },
            function (application_component_1_1) {
                application_component_1 = application_component_1_1;
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
            function (searchFilterByName_1_1) {
                searchFilterByName_1 = searchFilterByName_1_1;
            },
            function (searchFilterById_1_1) {
                searchFilterById_1 = searchFilterById_1_1;
            }],
        execute: function() {
            DashboardComponent = (function () {
                function DashboardComponent(_dom, _anakinService) {
                    this._dom = _dom;
                    this._anakinService = _anakinService;
                    this.timerSubscription = null;
                    this.loadingInstances = true;
                    this.loadingApps = true;
                    this.instancesElevation = 5;
                    this.animatedShadow = true;
                    this.selectedTab = 0;
                    this.selectedApp = 0;
                }
                DashboardComponent.prototype.ngOnInit = function () {
                    var _this = this;
                    console.log("DashboardComponent  init");
                    this.getApps();
                    this.getInstances();
                    var timer = Rx_1.Observable.timer(1000, 1000);
                    this.timerSubscription = timer.subscribe(function (t) {
                        _this.getInstances();
                    });
                };
                DashboardComponent.prototype.ngOnDestroy = function () {
                    if (this.timerSubscription != null) {
                        this.timerSubscription.unsubscribe();
                    }
                };
                DashboardComponent.prototype.getInstances = function () {
                    var _this = this;
                    this._anakinService.getAnakinInstances()
                        .subscribe(function (instances) { return _this.instances = instances; }, function (error) { return _this.instancesError = error; }, function () { return _this.getInstancesCompleted(); });
                };
                DashboardComponent.prototype.getInstancesCompleted = function () {
                    this.loadingInstances = false;
                    if (this.instancesError) {
                        console.error(this.instancesError);
                    }
                };
                DashboardComponent.prototype.getApps = function () {
                    var _this = this;
                    this._anakinService.getApplications().subscribe(function (apps) { return _this.apps = apps; }, function (error) { return _this.appsError = error; }, function () { return _this.getAppsCompleted(); });
                };
                DashboardComponent.prototype.getAppsCompleted = function () {
                    this.loadingApps = false;
                };
                DashboardComponent.prototype.onDashboardTabSelected = function (event) {
                    this.selectedTab = this._dom.nativeElement.querySelector("#dashboard-tabs").selected;
                };
                DashboardComponent = __decorate([
                    core_1.Component({
                        selector: 'dashboard',
                        templateUrl: 'app/dashboard.component.html',
                        directives: [application_component_1.ApplicationComponent],
                        pipes: [mapToIterable_1.MapToIterable,
                            slashIfMissing_1.SlashIfMissing,
                            dateFormatter_1.DateFormatter,
                            searchFilterByName_1.SearchFilterByName,
                            searchFilterById_1.SearchFilterById, common_1.DecimalPipe]
                    }), 
                    __metadata('design:paramtypes', [core_1.ElementRef, anakin_service_1.AnakinService])
                ], DashboardComponent);
                return DashboardComponent;
            }());
            exports_1("DashboardComponent", DashboardComponent);
        }
    }
});
//# sourceMappingURL=dashboard.component.js.map