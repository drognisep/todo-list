<template>
  <div id="footer">
    <div class="left">
      <p v-if="isTracking" v-text="taskDetailsString"></p>
    </div>
    <div class="right">
      <span v-show="hasWarnings" title="Warnings" class="warn-count badge"
            @click="showLogs = true">{{ warnings }}</span>
      <span v-show="hasErrors" title="Errors" class="error-count badge" @click="showLogs = true">{{ errors }}</span>
      <span v-show="noIssues" title="All good!" class="material-icons happy-indicator badge">done</span>
    </div>
  </div>
  <Modal
      v-show="showLogs"
      title="Logs"
      :close-handler="hideLogs"
  >
    <div class="log-view">
      <textarea readonly :value="logString"/>
    </div>
  </Modal>
</template>

<script>
import {EventsOff, EventsOn} from "../wailsjs/runtime/runtime.js";
import {LogEventName} from "../wailsjs/go/eventlog/EventLog.js";
import {GetTrackedTaskDetails} from "../wailsjs/go/main/ModelController.js";
import Modal from "./Modal.vue";
import {durationClock} from "../datetime.js";

export default {
  name: "Footer",
  components: {Modal},
  data() {
    return {
      errors: 0,
      warnings: 0,
      logs: [],
      showLogs: false,
      trackedDetails: null,
      secondsTracked: 0,
      secondsTicker: null,
    };
  },
  methods: {
    formatMessage(evt) {
      return `${evt.message}: ${JSON.stringify(evt.values)}`
    },
    onLogEventReceived(evt) {
      switch (evt.level) {
        case "DEBUG":
          console.debug(`[DEBUG] ${this.formatMessage(evt)}`);
          this.logs.push(evt);
          break;
        case "INFO":
          console.log(`[INFO] ${this.formatMessage(evt)}`);
          this.logs.push(evt);
          break;
        case "WARN":
          let warning = `[WARN] ${this.formatMessage(evt)}`;
          console.warn(warning);
          this.logs.push(evt);
          this.warnings++;
          break;
        case "ERROR":
          let error = `[ERROR] ${this.formatMessage(evt)}`;
          console.error(error);
          this.logs.push(evt);
          this.errors++;
          break;
      }
    },
    hideLogs() {
      this.showLogs = false;
    },
    onTaskStarted() {
      GetTrackedTaskDetails()
          .then(details => {
            if (details != null) {
              this.trackedDetails = details;
              this.secondsTicker = setInterval(() => {
                this.secondsTracked = Math.floor((Date.now() - new Date(this.trackedDetails.entry.start)) / 1000);
              }, 1000);
            }
          })
          .catch(console.errorEvent);
    },
    onTaskStopped() {
      clearInterval(this.secondsTicker);
      this.trackedDetails = null;
      this.secondsTracked = 0;
    },
  },
  computed: {
    logString() {
      let out = "";
      this.logs.forEach(log => {
        out += `${log.time} [${log.level}] ${log.message}: ${JSON.stringify(log.values)}` + "\n";
      })
      return out;
    },
    noIssues() {
      return this.errors === 0 && this.warnings === 0;
    },
    hasWarnings() {
      return this.warnings > 0;
    },
    hasErrors() {
      return this.errors > 0;
    },
    isTracking() {
      return this.trackedDetails != null;
    },
    timeString() {
      let now = Date.now();
      let start = now - this.secondsTracked * 1000
      return durationClock(start, now);
    },
    taskDetailsString() {
      if (this.trackedDetails == null) {
        return "";
      }
      return `[${this.timeString}] ${this.trackedDetails.task.name}`;
    },
  },
  created() {
    LogEventName()
        .then(name => {
          EventsOn(name, this.onLogEventReceived)
        })
    GetTrackedTaskDetails()
        .then(details => {
          if (details != null) {
            this.trackedDetails = details;
            this.secondsTracked = Math.floor((Date.now() - new Date(this.trackedDetails.entry.start)) / 1000);
            this.secondsTicker = setInterval(() => {
              this.secondsTracked = Math.floor((Date.now() - new Date(this.trackedDetails.entry.start)) / 1000);
            }, 1000);
          }
        })
        .catch(console.errorEvent);
    EventsOn("taskStarted", this.onTaskStarted)
    EventsOn("taskStopped", this.onTaskStopped)
  },
  destroyed() {
    LogEventName()
        .then(name => {
          EventsOff(name)
        })
    this.onTaskStopped();
  }
}
</script>

<style scoped>
#footer {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background-color: var(--bg-color-light);
  height: var(--footer-height);
  z-index: var(--z-toolbar);
  box-shadow: 0 -1px 5px black;
  display: flex;
  justify-content: space-between;
}

#footer .left, #footer .right {
  padding: 0 16px;
  display: flex;
  height: 100%;
  align-items: center;
}

#footer .left > *, #footer .right > * {
  margin-left: 8px;
}

#footer .right .error-count {
  background-color: var(--bg-danger);
  cursor: pointer;
}

#footer .right .warn-count {
  background-color: var(--bg-warn);
  cursor: pointer;
}

#footer .right .happy-indicator {
  background-color: var(--fg-happy);
}

.badge {
  padding: 2px 4px;
  color: var(--fg-color);
  text-shadow: 0 0 8px black;
  border-radius: 4px;
  font-weight: bold;
}

.log-view {
  width: 100%;
  height: 100%;
  overflow: hidden;
}

.log-view textarea {
  width: 100%;
  height: 100%;
  resize: none;
}
</style>