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
    var SearchFilterByName;
    return {
        setters:[
            function (core_1_1) {
                core_1 = core_1_1;
            }],
        execute: function() {
            SearchFilterByName = (function () {
                function SearchFilterByName() {
                }
                SearchFilterByName.prototype.transform = function (value, args) {
                    if (args[0] == 'undefined') {
                        return value;
                    }
                    var filter = args[0].toLocaleLowerCase();
                    console.log(filter);
                    return filter ? value.filter(function (entity) { return entity.name.toLocaleLowerCase().indexOf(filter) != -1; }) : value;
                };
                SearchFilterByName = __decorate([
                    core_1.Pipe({
                        name: 'searchFilterByName'
                    }), 
                    __metadata('design:paramtypes', [])
                ], SearchFilterByName);
                return SearchFilterByName;
            }());
            exports_1("SearchFilterByName", SearchFilterByName);
        }
    }
});
//# sourceMappingURL=searchFilterByName.js.map