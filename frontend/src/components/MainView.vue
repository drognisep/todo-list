<template>
  <div id="main-view">
    <router-view></router-view>
  </div>
  <Confirm/>
</template>

<script>
import {RouterView} from 'vue-router';
import Confirm from "./Confirm.vue";

export default {
  name: "MainView",
  components: {RouterView, Confirm},
}
</script>

<style scoped>
#main-view {
  position: absolute;
  top: var(--toolbar-height);
  left: 0;
  right: 0;
  bottom: 0;
  overflow: hidden;
}
</style>