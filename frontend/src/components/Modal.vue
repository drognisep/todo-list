<template>
  <overlay @overlayClicked="dialogClose">
    <div class="dialog" @click.stop.prevent>
      <h1 class="dialog-header">{{ this.$props.title }}</h1>
      <div class="dialog-content">
        <slot></slot>
      </div>
      <span class="material-icons dialog-x" @click.stop.prevent="dialogClose">close</span>
    </div>
  </overlay>
</template>

<script>
import Overlay from "./Overlay.vue";

export default {
  name: "Modal",
  components: {Overlay},
  props: {
    "title": String,
  },
  emits: ["dialogClose"],
  methods: {
    dialogClose() {
      this.$emit("dialogClose")
    },
  },
}
</script>

<style scoped>
.dialog {
  display: flex;
  flex-direction: column;
  background-color: var(--bg-color);
  box-shadow: 0 0 8px black;
  position: fixed;
  height: calc(66vh - var(--toolbar-height));
  top: calc(16vh + var(--toolbar-height));
  border-radius: 16px;
  padding: 16px;
  z-index: var(--z-dialog);
}

.dialog .dialog-header {
  margin: 0;
}

.dialog .dialog-content {
  overflow-y: auto;
  flex: auto;
}

.dialog .dialog-x {
  color: var(--fg-color);
  position: absolute;
  top: 12px;
  right: 12px;
  cursor: pointer;
}

/* Dialog media queries */

@media (min-width: 0px) and (max-width: 449px) {
  .dialog {
    width: 90vw;
    left: 5vw;
  }
}

@media (min-width: 450px) and (max-width: 984px) {
  .dialog {
    width: 66vw;
    left: 16vw;
  }
}

@media (min-width: 985px) {
  .dialog {
    width: 650px;
    left: calc((100vw - 640px) / 2)
  }
}
</style>