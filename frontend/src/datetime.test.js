import {expect, test} from "vitest";
import {durationClock, durationGo, formatClockTime, weekday, weekdaySemantic} from "./datetime.js";

const second = 1000;
const minute = second * 60;
const hour = minute * 60;

test('Duration clock should round', () => {
    let start = Date.now();
    let end = start + 400;
    expect(durationClock(start, end)).eq("00:00:00");
    end = start + 500;
    expect(durationClock(start, end)).eq("00:00:01");
})

test('Duration is able to handle a missing end date by reporting 0 duration', () => {
    let start = Date.now();
    expect(durationClock(start, null)).eq("00:00:00");
    expect(durationClock(start, undefined)).eq("00:00:00");
    expect(durationClock(start, 0)).eq("00:00:00");
    expect(durationClock(start, '')).eq("00:00:00");
})

test('Duration clock represents each unit correctly', () => {
    let start = Date.now();
    let end = start + 59 * second + minute + hour;
    expect(durationClock(start, end)).eq("01:01:59");
})

test('Duration human should default to 0s', () => {
    let now = Date.now();
    expect(durationGo(now, now)).eq('0s');
})

test('Duration human should represent each unit correctly', () => {
    let start = Date.now();
    let end = start + 59 * second + minute + hour;
    expect(durationGo(start, end)).eq("1h1m59s");
})

test('Duration should convert strings to Dates', () => {
    let now = Date.now();
    let later = now + 5 * minute;
    let nowStr = now.toString();
    let laterStr = later.toString();
    expect(durationGo(nowStr, laterStr)).eq('5m');
})

test('Weekday should identify what day of the week the provided date falls on', () => {
    function testDate() {
        return new Date(2023, 0, 3);
    }

    expect(weekday(testDate() - 48 * hour)).eq('Sunday');
    expect(weekday(testDate() - 24 * hour)).eq('Monday');
    expect(weekday(testDate())).eq('Tuesday');
    // Why subtracting a negative instead of addition? Because JavaScript!
    // https://stackoverflow.com/questions/26322237/subtraction-working-in-js-but-addition-is-not-working-in-js
    expect(weekday(testDate() - -24 * hour)).eq('Wednesday');
    expect(weekday(testDate() - -48 * hour)).eq('Thursday');
    expect(weekday(testDate() - -72 * hour)).eq('Friday');
    expect(weekday(testDate() - -96 * hour)).eq('Saturday');
})

test('Weekday semantic should identify yesterday, today, and tomorrow', () => {
    function testDate() {
        return new Date();
    }

    expect(weekdaySemantic(testDate() - 24 * hour)).eq('Yesterday');
    expect(weekdaySemantic(testDate())).eq('Today');
    // Why subtracting a negative instead of addition? Because JavaScript!
    // https://stackoverflow.com/questions/26322237/subtraction-working-in-js-but-addition-is-not-working-in-js
    expect(weekdaySemantic(testDate() - -24 * hour)).eq('Tomorrow');
})

test('Format clock time should work as expected', () => {
    let date = new Date(2023, 0, 3, 1, 2, 3, 100);
    expect(formatClockTime(date)).eq('01:02:03');
})
