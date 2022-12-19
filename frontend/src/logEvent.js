import {InfoEvent, WarnEvent, ErrorEvent, DebugEvent} from "./wailsjs/go/eventlog/EventLog.js";

export default {
    setup() {
        console.logEvent = (msg) => {
            InfoEvent(msg, [])
                .catch(console.error)
        };
        console.warnEvent = (msg) => {
            WarnEvent(msg, [])
                .catch(console.error)
        };
        console.errorEvent = (msg) => {
            ErrorEvent(msg, [])
                .catch(console.error)
        };
        console.debugEvent = (msg) => {
            DebugEvent(msg, [])
                .catch(console.error)
        };
    }
}
