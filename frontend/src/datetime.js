function coerceDate(d) {
    if (typeof d === 'string') {
        let s = Date.parse(d);
        if (isNaN(s)) {
            let n = Number.parseInt(d)
            if (isNaN(n)) {
                return 0;
            }
            return n;
        }
        return s;
    }
    return d;
}

function duration(start, end) {
    start = coerceDate(start);
    end = coerceDate(end);

    let seconds = Math.floor((end - start) / 1000);
    let minutes = 0;
    while(seconds >= 60) {
        minutes++;
        seconds -= 60;
    }
    let hours = 0;
    while(minutes >= 60) {
        hours++;
        minutes -= 60;
    }

    return [hours, minutes, seconds];
}

export function durationClock(start, end) {
    let [hours, minutes, seconds] = duration(start, end);
    if (hours < 10) {
        hours = '0' + hours;
    }
    if (minutes < 10) {
        minutes = '0' + minutes;
    }
    if (seconds < 10) {
        seconds = '0' + seconds;
    }
    return `${hours}:${minutes}:${seconds}`;
}

export function durationGo(start, end) {
    let [hours, minutes, seconds] = duration(start, end);
    let str = '';
    if (hours > 0) {
        str += `${hours}h`;
    }
    if (minutes > 0) {
        str += `${minutes}m`;
    }
    if (seconds > 0) {
        str += `${seconds}s`;
    }

    if (hours === 0 && minutes === 0 && seconds === 0) {
        str = '0s';
    }
    return str;
}
