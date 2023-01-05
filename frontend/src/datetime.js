export function coerceDate(d) {
    if (typeof d === 'string') {
        let s = Date.parse(d);
        if (isNaN(s)) {
            let n = Number.parseInt(d)
            if (isNaN(n)) {
                return new Date(0);
            }
            return new Date(n);
        }
        return new Date(s);
    }
    return new Date(d);
}

function duration(start, end) {
    if (!end) {
        return [0, 0, 0];
    }
    start = coerceDate(start);
    end = coerceDate(end);

    let seconds = Math.floor((end - start) / 1000);
    let minutes = 0;
    while (seconds >= 60) {
        minutes++;
        seconds -= 60;
    }
    let hours = 0;
    while (minutes >= 60) {
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

const weekdayMap = {
    0: 'Sunday',
    1: 'Monday',
    2: 'Tuesday',
    3: 'Wednesday',
    4: 'Thursday',
    5: 'Friday',
    6: 'Saturday',
}

export function weekday(date) {
    date = coerceDate(date);
    return weekdayMap[date.getDay()];
}

const second = 1000;
const minute = second * 60;
const hour = minute * 60;

function dateMatches(base, other) {
    if (base.getFullYear() === other.getFullYear()) {
        if (base.getMonth() === other.getMonth()) {
            if (base.getDate() === other.getDate()) {
                return true;
            }
        }
    }
    return false;
}

export function weekdaySemantic(date) {
    date = coerceDate(date);
    let today = new Date();
    let yesterday = new Date(today - 24 * hour);
    let tomorrow = new Date(today - -24 * hour);
    if (dateMatches(yesterday, date)) {
        return 'Yesterday';
    }
    if (dateMatches(today, date)) {
        return 'Today';
    }
    if (dateMatches(tomorrow, date)) {
        return 'Tomorrow';
    }
    return weekday(date);
}

export function formatClockTime(date) {
    date = coerceDate(date);
    let hours = date.getHours();
    let minutes = date.getMinutes();
    let seconds = date.getSeconds();
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
