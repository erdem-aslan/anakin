/**
 * Created by erdemaslan on 10/05/16.
 */
import {Pipe, PipeTransform} from 'angular2/core';

@Pipe({
    name: 'searchFilterByName'
})

export class SearchFilterByName implements PipeTransform {
    transform(value:any, args:string[]):any {

        if (args[0] == 'undefined') {
            return value;
        }

        let filter = args[0].toLocaleLowerCase();
        console.log(filter);
        return filter ? value.filter(entity=> entity.name.toLocaleLowerCase().indexOf(filter) != -1) : value;
    }
}
