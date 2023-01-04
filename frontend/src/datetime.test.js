import {test, expect} from "vitest";
import {durationGo, durationClock} from "./datetime.js";

const second = 1000;
const minute = second * 60;
const hour = minute * 60;

test('Duration clock should round', () => {
    let start = Date.now();
    let end = start + 900;
    expect(durationClock(start, end)).eq("00:00:00");
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
