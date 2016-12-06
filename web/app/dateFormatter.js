System.register(['angular2/core'], function(exports_1, context_1) {
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
    var core_1;
    var DateFormatter;
    return {
        setters:[
            function (core_1_1) {
                core_1 = core_1_1;
            }],
        execute: function() {
            /**
             *
             * Prepends forward slash if missing
             */
            DateFormatter = (function () {
                function DateFormatter() {
                }
                DateFormatter.prototype.transform = function (value) {
                    var valueDate = new Date();
                    valueDate.setTime(Date.parse(value));
                    return valueDate.toLocaleDateString() + " " + valueDate.toLocaleTimeString();
                };
                DateFormatter = __decorate([
                    core_1.Pipe({ name: 'dateFormatter' }), 
                    __metadata('design:paramtypes', [])
                ], DateFormatter);
                return DateFormatter;
            }());
            exports_1("DateFormatter", DateFormatter);
        }
    }
});
//# sourceMappingURL=dateFormatter.js.map