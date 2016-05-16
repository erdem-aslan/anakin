System.register(['angular2/core', "./anakin.service", "./mapToIterable", "./slashIfMissing", "./searchFilterByName", "angular2/common", "./application.component"], function(exports_1, context_1) {
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
    var core_1, anakin_service_1, mapToIterable_1, slashIfMissing_1, searchFilterByName_1, common_1, application_component_1;
    var ConfigurationComponent;
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
            function (searchFilterByName_1_1) {
                searchFilterByName_1 = searchFilterByName_1_1;
            },
            function (common_1_1) {
                common_1 = common_1_1;
            },
            function (application_component_1_1) {
                application_component_1 = application_component_1_1;
            }],
        execute: function() {
            ConfigurationComponent = (function () {
                function ConfigurationComponent(_anakinService, _dom) {
                    this._anakinService = _anakinService;
                    this._dom = _dom;
                    this.filter = '';
                    this.newApp = { id: '', name: '', baseUrl: '', services: {}, state: '', error: {} };
                }
                ConfigurationComponent.prototype.ngOnInit = function () {
                    this.getApps();
                };
                ConfigurationComponent.prototype.getApps = function () {
                    var _this = this;
                    this._anakinService.getApplications().subscribe(function (apps) { return _this.apps = apps; }, function (error) { return _this.getApplicationsError(error); }, function () { return _this.getApplicationsFinished(); });
                };
                ConfigurationComponent.prototype.onCreateNewApplication = function () {
                    var _this = this;
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
                        .subscribe(function (app) { return _this.newApp = app; }, function (error) { return _this.createApplicationError(error); }, function () { return _this.createApplicationFinished(); });
                };
                ConfigurationComponent.prototype.createApplicationFinished = function () {
                    console.log(this.newApp);
                    this.getApps();
                };
                ConfigurationComponent.prototype.createApplicationError = function (error) {
                    this.error = error;
                    this._dom.nativeElement.querySelector("#createAppError").open();
                };
                ConfigurationComponent.prototype.getApplicationsError = function (error) {
                    this.appsError = error;
                    this._dom.nativeElement.querySelector("#getAppsError").open();
                };
                ConfigurationComponent.prototype.getApplicationsFinished = function () {
                    console.log(this.apps);
                };
                ConfigurationComponent.prototype.onAppChangeError = function (event) {
                    console.log(event);
                    this.appChangeError = event.value;
                    this._dom.nativeElement.querySelector("#appChangeError").open();
                };
                ConfigurationComponent.prototype.onAppChanged = function (event) {
                    console.log(event);
                    this.getApps();
                };
                ConfigurationComponent = __decorate([
                    core_1.Component({
                        selector: 'configuration',
                        templateUrl: 'app/configuration.component.html',
                        pipes: [mapToIterable_1.MapToIterable, slashIfMissing_1.SlashIfMissing, searchFilterByName_1.SearchFilterByName],
                        directives: [common_1.NgSwitch, common_1.NgSwitchWhen, common_1.FORM_DIRECTIVES, application_component_1.ApplicationComponent]
                    }), 
                    __metadata('design:paramtypes', [anakin_service_1.AnakinService, core_1.ElementRef])
                ], ConfigurationComponent);
                return ConfigurationComponent;
            }());
            exports_1("ConfigurationComponent", ConfigurationComponent);
        }
    }
});
//# sourceMappingURL=configuration.component.js.map