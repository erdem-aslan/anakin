System.register(['angular2/core', "angular2/router", "./dashboard.component", "./configuration.component", "./statistics.component", "./monitoring.component", "./anakin.service"], function(exports_1, context_1) {
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
    var core_1, router_1, router_2, dashboard_component_1, configuration_component_1, statistics_component_1, monitoring_component_1, anakin_service_1;
    var AppComponent, AnakinInstance, Endpoint;
    return {
        setters:[
            function (core_1_1) {
                core_1 = core_1_1;
            },
            function (router_1_1) {
                router_1 = router_1_1;
                router_2 = router_1_1;
            },
            function (dashboard_component_1_1) {
                dashboard_component_1 = dashboard_component_1_1;
            },
            function (configuration_component_1_1) {
                configuration_component_1 = configuration_component_1_1;
            },
            function (statistics_component_1_1) {
                statistics_component_1 = statistics_component_1_1;
            },
            function (monitoring_component_1_1) {
                monitoring_component_1 = monitoring_component_1_1;
            },
            function (anakin_service_1_1) {
                anakin_service_1 = anakin_service_1_1;
            }],
        execute: function() {
            AppComponent = (function () {
                function AppComponent(_router, _dom, _anakinService) {
                    this._router = _router;
                    this._dom = _dom;
                    this._anakinService = _anakinService;
                }
                AppComponent.prototype.ngOnInit = function () {
                    console.log("App Component init");
                    this.version = this._anakinService.getAnakinVersion();
                };
                AppComponent.prototype.onDashboardSelected = function () {
                    this._router.navigate(['Dashboard']);
                    this.toggleAnakinDrawer();
                };
                AppComponent.prototype.toggleAnakinDrawer = function () {
                    this._dom.nativeElement.querySelector("#anakin-drawer").togglePanel();
                };
                AppComponent.prototype.onConfigurationSelected = function () {
                    this._router.navigate(['Configuration']);
                    this.toggleAnakinDrawer();
                };
                AppComponent.prototype.onStatisticsSelected = function () {
                    this._router.navigate(['Statistics']);
                    this.toggleAnakinDrawer();
                };
                AppComponent.prototype.onMonitoringSelected = function () {
                    this._router.navigate(['Monitoring']);
                    this.toggleAnakinDrawer();
                };
                AppComponent = __decorate([
                    core_1.Component({
                        selector: 'app',
                        templateUrl: 'app/app.component.html',
                        directives: [router_1.ROUTER_DIRECTIVES],
                        providers: [anakin_service_1.AnakinService]
                    }),
                    router_2.RouteConfig([
                        { path: 'dashboard', name: 'Dashboard', component: dashboard_component_1.DashboardComponent, useAsDefault: true },
                        { path: 'configuration', name: 'Configuration', component: configuration_component_1.ConfigurationComponent },
                        { path: 'statistics', name: 'Statistics', component: statistics_component_1.StatisticsComponent },
                        { path: 'monitoring', name: 'Monitoring', component: monitoring_component_1.MonitoringComponent }
                    ]), 
                    __metadata('design:paramtypes', [router_1.Router, core_1.ElementRef, anakin_service_1.AnakinService])
                ], AppComponent);
                return AppComponent;
            }());
            exports_1("AppComponent", AppComponent);
            AnakinInstance = (function () {
                function AnakinInstance() {
                }
                return AnakinInstance;
            }());
            exports_1("AnakinInstance", AnakinInstance);
            Endpoint = (function () {
                function Endpoint() {
                }
                return Endpoint;
            }());
            exports_1("Endpoint", Endpoint);
        }
    }
});
//# sourceMappingURL=app.component.js.map