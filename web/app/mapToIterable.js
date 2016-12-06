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
    var MapToIterable;
    return {
        setters:[
            function (core_1_1) {
                core_1 = core_1_1;
            }],
        execute: function() {
            /**
             * Map to Iteratble Pipe
             *
             * It accepts Objects and [Maps](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Map)
             *
             * Example:
             *
             *  <div *ngFor="#keyValuePair of someObject | mapToIterable">
             *    key {{keyValuePair.key}} and value {{keyValuePair.value}}
             *  </div>
             *
             */
            MapToIterable = (function () {
                function MapToIterable() {
                }
                MapToIterable.prototype.transform = function (value) {
                    var result = [];
                    if (value.entries) {
                        for (var _i = 0, _a = value.entries(); _i < _a.length; _i++) {
                            var _b = _a[_i], key = _b[0], value = _b[1];
                            result.push({ key: key, value: value });
                        }
                    }
                    else {
                        for (var key_1 in value) {
                            result.push({ key: key_1, value: value[key_1] });
                        }
                    }
                    return result;
                };
                MapToIterable = __decorate([
                    core_1.Pipe({ name: 'mapToIterable' }), 
                    __metadata('design:paramtypes', [])
                ], MapToIterable);
                return MapToIterable;
            }());
            exports_1("MapToIterable", MapToIterable);
        }
    }
});
//# sourceMappingURL=mapToIterable.js.map