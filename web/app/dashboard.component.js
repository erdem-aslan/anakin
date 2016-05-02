System.register(['angular2/core', "./anakin.service", "./mapToIterable"], function(exports_1, context_1) {
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
    var core_1, anakin_service_1, mapToIterable_1;
    var DashboardComponent;
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
            }],
        execute: function() {
            DashboardComponent = (function () {
                function DashboardComponent(_dom, _anakinService) {
                    this._dom = _dom;
                    this._anakinService = _anakinService;
                    this.loadingApps = true;
                    this.loadingServices = true;
                    this.loadingEndpoints = true;
                    this.appsElevation = 1;
                }
                DashboardComponent.prototype.ngOnInit = function () {
                    console.log("DashboardComponent  init");
                    this.getApplications();
                };
                DashboardComponent.prototype.onHoverApplications = function () {
                    this.appsElevation = 5;
                };
                DashboardComponent.prototype.onLeaveApplications = function () {
                    this.appsElevation = 1;
                };
                DashboardComponent.prototype.getApplications = function () {
                    var _this = this;
                    this._anakinService.getApplications()
                        .subscribe(function (apps) { return _this.apps = apps; }, function (error) { return _this.errorString = error; }, function () { return _this.getApplicationsCompleted(); });
                };
                DashboardComponent.prototype.getApplicationsCompleted = function () {
                    console.log("getApps finished");
                    this.loadingApps = false;
                    if (this.errorString) {
                        console.error(this.errorString);
                    }
                };
                DashboardComponent = __decorate([
                    core_1.Component({
                        selector: 'dashboard',
                        templateUrl: 'app/dashboard.component.html',
                        pipes: [mapToIterable_1.MapToIterable]
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