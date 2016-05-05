/**
 * Created by erdemaslan on 27/04/16.
 */
import {Pipe, PipeTransform} from 'angular2/core';
/**
 * 
 * Prepends forward slash if missing
 */
@Pipe({ name: 'slashIfMissing' })
export class SlashIfMissing implements PipeTransform {

    transform(value) {

        if (!value.startsWith("/")) {
            return "/" + value;
        }
        
        return value
    }
}
