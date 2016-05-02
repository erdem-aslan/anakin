System.register(['angular2/core', "./anakin.service"], function(exports_1, context_1) {
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
    var core_1, anakin_service_1;
    var ConfigurationComponent;
    return {
        setters:[
            function (core_1_1) {
                core_1 = core_1_1;
            },
            function (anakin_service_1_1) {
                anakin_service_1 = anakin_service_1_1;
            }],
        execute: function() {
            ConfigurationComponent = (function () {
                function ConfigurationComponent(_anakinService) {
                    this._anakinService = _anakinService;
                }
                ConfigurationComponent.prototype.ngOnInit = function () {
                    console.log("Configuration component init");
                };
                ConfigurationComponent = __decorate([
                    core_1.Component({
                        selector: 'configuration',
                        templateUrl: 'app/configuration.component.html'
                    }), 
                    __metadata('design:paramtypes', [anakin_service_1.AnakinService])
                ], ConfigurationComponent);
                return ConfigurationComponent;
            }());
            exports_1("ConfigurationComponent", ConfigurationComponent);
        }
    }
});
//# sourceMappingURL=configuration.component.js.map