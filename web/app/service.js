System.register([], function(exports_1, context_1) {
    "use strict";
    var __moduleName = context_1 && context_1.id;
    var Service, Endpoint;
    return {
        setters:[],
        execute: function() {
            /**
             * Created by erdemaslan on 27/04/16.
             */
            Service = (function () {
                function Service(id, serviceUrl, endpoints, realEndpoints, balanceStrategy, nested, state, editable, bsIndex) {
                    this.id = id;
                    this.serviceUrl = serviceUrl;
                    this.endpoints = endpoints;
                    this.realEndpoints = realEndpoints;
                    this.balanceStrategy = balanceStrategy;
                    this.nested = nested;
                    this.state = state;
                    this.editable = editable;
                    this.bsIndex = bsIndex;
                }
                return Service;
            }());
            exports_1("Service", Service);
            Endpoint = (function () {
                function Endpoint(id, host, port, scheme, state) {
                    this.id = id;
                    this.host = host;
                    this.port = port;
                    this.scheme = scheme;
                    this.state = state;
                }
                return Endpoint;
            }());
            exports_1("Endpoint", Endpoint);
        }
    }
});
//# sourceMappingURL=service.js.map