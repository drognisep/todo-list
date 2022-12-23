import {InfoEvent, WarnEvent, ErrorEvent, DebugEvent, SetDebug} from "./wailsjs/go/eventlog/EventLog.js";

const debugEnabled = false;

export default {
    setup() {
        console.logEvent = (msg, args) => {
            if (!args) {
                args = [];
            }
            InfoEvent(msg, args)
                .catch(console.error)
        };
        console.warnEvent = (msg, args) => {
            if (!args) {
                args = [];
            }
            WarnEvent(msg, args)
                .catch(console.error)
        };
        console.errorEvent = (msg, args) => {
            if (!args) {
                args = [];
            }
            ErrorEvent(msg, args)
                .catch(console.error)
        };
        console.debugEvent = (msg, args) => {
            if (!args) {
                args = [];
            }
            DebugEvent(msg, args)
                .catch(console.error)
        };
        SetDebug(debugEnabled)
            .catch(console.errorEvent)
    }
}
