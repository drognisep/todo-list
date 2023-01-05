import {InfoEvent, WarnEvent, ErrorEvent, DebugEvent, SetDebug} from "./wailsjs/go/eventlog/EventLog.js";

const debugEnabled = false;

export default {
    setup() {
        console.logEvent = function (msg, ...args) {
            if (typeof msg !== 'string') {
                msg = `${msg}`;
            }
            InfoEvent(msg, args)
                .catch(console.error)
        };
        console.warnEvent = function (msg, ...args) {
            console.trace();
            if (typeof msg !== 'string') {
                msg = `${msg}`;
            }
            WarnEvent(msg, args)
                .catch(console.error)
        };
        console.errorEvent = function (msg, ...args) {
            console.trace();
            if (typeof msg !== 'string') {
                msg = `${msg}`;
            }
            ErrorEvent(msg, args)
                .catch(console.error)
        };
        console.debugEvent = function (msg, ...args) {
            if (typeof msg !== 'string') {
                msg = `${msg}`;
            }
            DebugEvent(msg, args)
                .catch(console.error)
        };
        SetDebug(debugEnabled)
            .catch(console.errorEvent)
    }
}
