/**
 * Created by erdemaslan on 06/05/16.
 */
import {Pipe, PipeTransform} from 'angular2/core';
/**
 *
 * Prepends forward slash if missing
 */
@Pipe({ name: 'dateFormatter' })
export class DateFormatter implements PipeTransform {

    transform(value) {
        var valueDate:Date = new Date();
        valueDate.setTime(Date.parse(value));

        return valueDate.toLocaleDateString() + " " + valueDate.toLocaleTimeString();
    }
}
