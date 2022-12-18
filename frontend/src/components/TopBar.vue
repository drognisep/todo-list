<template>
  <div id="topbar">
    <div class="nav">
      <img alt="tasks" id="brand-img" src="../assets/images/gear.svg"/>
      <h3 id="brand-text">Tasks</h3>
      <RouterLink to="/" class="link">Dashboard</RouterLink>
      <RouterLink to="/allTasks" class="link">Tasks</RouterLink>
    </div>
    <div class="actions">
      <span class="material-icons export" @click="exportModel" title="Export">save_as</span>
      <span class="material-icons import" @click="showImportDialog = true" title="Import">download</span>
    </div>
  </div>
  <ImportDialog
      v-show="showImportDialog"
      @dialogClosed="closeImport"
      @doImport="importModel"
  />
</template>

<script>
import {RouterLink} from 'vue-router';
import {Import, Export} from "../wailsjs/go/main/TaskController.js";
import ImportDialog from "./ImportDialog.vue";

export default {
  name: "TopBar",
  components: {RouterLink, ImportDialog},
  data() {
    return {
      showImportDialog: false,
    };
  },
  methods: {
    exportModel() {
      showProgress("Exporting model...");
      Export()
          .then(snapshotFile => {
            if (snapshotFile !== "") {
              confirmDialog("Export complete!", `Your data has been exported to ${snapshotFile}`, () => {});
            }
          })
          .catch(console.error)
          .then(() => {
            closeProgress();
          })
    },
    importModel(data) {
      showProgress("Importing snapshot...");
      Import(data.strategy)
          .then(() => {
            confirmDialog("Import successful", "Click OK to reload the app", () => {
              window.location.reload();
            })
          })
          .catch(console.error)
          .then(() => {
            closeProgress();
          })
    },
    closeImport() {
      this.showImportDialog = false;
    },
  },
}
</script>

<style scoped>
#topbar {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: var(--toolbar-height);
  background-color: var(--bg-color-light);
  box-shadow: black 0 0 5px;
  padding: 0 16px;
  z-index: var(--z-toolbar);

  display: flex;
  justify-content: space-between;

  --text-shadow: black 0 0 8px;
  --border-right: 2px solid var(--overlay-color);
}

#topbar > .nav {
  display: flex;
  height: 100%;
}

.nav > * {
  display: inline-block;
  max-height: var(--toolbar-height);
  line-height: var(--toolbar-height);
  margin: 0;
  padding: 0;
}

.nav > #brand-img {
  height: calc(var(--toolbar-height) - 12px);
  margin: 6px 16px 0 0;
  float: left;
}

.nav > #brand-text {
  font-weight: bold;
  height: var(--toolbar-height);
  line-height: var(--toolbar-height);
  padding: 0 16px 0 0;
  cursor: default;
  text-shadow: var(--text-shadow);
  border-right: var(--border-right);
}

.nav > .link {
  color: var(--fg-color);
  text-decoration: none;
  position: relative;
  padding: 0 8px;
  margin: 0;
  cursor: pointer;
  transition: all 200ms linear;
  border-right: var(--border-right);
}

.nav > .link:hover, .nav > .link.router-link-active {
  text-shadow: var(--text-shadow);
  color: var(--fg-highlight);
  background-color: var(--bg-highlight);
}

div.actions {
  text-align: right;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: end;
}

div.actions .export, div.actions .import {
  cursor: pointer;
  text-shadow: 0 0 4px black;
  font-size: 1.8em;
  border-radius: 4px;
  padding: 4px;
}

div.actions .export:hover, div.actions .import:hover {
  background-color: var(--bg-highlight);
}
</style>