System.register([], function(exports_1, context_1) {
    "use strict";
    var __moduleName = context_1 && context_1.id;
    var Instance, InstanceStats;
    return {
        setters:[],
        execute: function() {
            /**
             * Created by erdemaslan on 06/05/16.
             */
            Instance = (function () {
                function Instance(id, version, adminPort, adminIp, proxyIp, proxyPort, started, state, stats) {
                    this.id = id;
                    this.version = version;
                    this.adminPort = adminPort;
                    this.adminIp = adminIp;
                    this.proxyIp = proxyIp;
                    this.proxyPort = proxyPort;
                    this.started = started;
                    this.state = state;
                    this.stats = stats;
                }
                return Instance;
            }());
            exports_1("Instance", Instance);
            InstanceStats = (function () {
                function InstanceStats(os, cpuCores, mem, rps) {
                    this.os = os;
                    this.cpuCores = cpuCores;
                    this.mem = mem;
                    this.rps = rps;
                }
                return InstanceStats;
            }());
            exports_1("InstanceStats", InstanceStats);
        }
    }
});
//# sourceMappingURL=instance.js.map