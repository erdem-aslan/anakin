/**
 * Created by erdemaslan on 19/04/16.
 */
import {Component, OnInit, ElementRef} from 'angular2/core';
import {AnakinService} from "./anakin.service"
import {Application} from "./application";
import {MapToIterable} from "./mapToIterable"
import {SlashIfMissing} from "./slashIfMissing";

@Component({
    selector:'configuration',
    templateUrl:'app/configuration.component.html',
    pipes: [MapToIterable, SlashIfMissing],
})

export class ConfigurationComponent implements OnInit {

    constructor(private _anakinService:AnakinService) {

    }
    
    ngOnInit() {
        console.log("Configuration component init")
    }
}
