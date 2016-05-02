/**
 * Created by erdemaslan on 19/04/16.
 */
System.register(['angular2/core', "angular2/http", "rxjs/Observable", 'rxjs/Rx'], function(exports_1, context_1) {
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
    var core_1, http_1, Observable_1;
    var AnakinService;
    return {
        setters:[
            function (core_1_1) {
                core_1 = core_1_1;
            },
            function (http_1_1) {
                http_1 = http_1_1;
            },
            function (Observable_1_1) {
                Observable_1 = Observable_1_1;
            },
            function (_1) {}],
        execute: function() {
            AnakinService = (function () {
                function AnakinService(http) {
                    this.http = http;
                    this.version = "1.0";
                }
                AnakinService.prototype.getAnakinVersion = function () {
                    return this.version;
                };
                AnakinService.prototype.getAnakinInstances = function () {
                };
                AnakinService.prototype.getApplications = function () {
                    return this.http.get("/anakin/v1/apps")
                        .map(function (res) { return res.json(); })
                        .catch(this.handleError);
                };
                AnakinService.prototype.getServices = function (applicationId) {
                };
                AnakinService.prototype.getEndpoints = function (serviceId) {
                };
                AnakinService.prototype.handleError = function (error) {
                    console.error(error);
                    return Observable_1.Observable.throw(error.json().error || 'Server error');
                };
                AnakinService = __decorate([
                    core_1.Injectable(), 
                    __metadata('design:paramtypes', [http_1.Http])
                ], AnakinService);
                return AnakinService;
            }());
            exports_1("AnakinService", AnakinService);
        }
    }
});
//# sourceMappingURL=anakin.service.js.map