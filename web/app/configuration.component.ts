/**
 * Created by erdemaslan on 19/04/16.
 */
import {Component, OnInit} from 'angular2/core';
import {AnakinService} from "./anakin.service";

@Component({
    selector:'configuration',
    templateUrl:'app/configuration.component.html'
})

export class ConfigurationComponent implements OnInit {

    constructor(private _anakinService:AnakinService) {

    }
    
    ngOnInit() {
        console.log("Configuration component init")
    }
}
