/**
 * Created by erdemaslan on 10/05/16.
 */
import {Pipe, PipeTransform} from 'angular2/core';

@Pipe({
    name: 'searchFilterById'
})

export class SearchFilterById implements PipeTransform {
    transform(value:any, args:string[]):any {
        let filter = args[0].toLocaleLowerCase();
        return filter ? value.filter(entity=> entity.id.toLocaleLowerCase().indexOf(filter) != -1) : value;
    }
}
