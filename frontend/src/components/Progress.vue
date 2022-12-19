<template>
  <div v-show="showProgress" class="outer" @click.stop.prevent>
    <div class="inner">
      <div class="progress">
        <h3>{{ message }}</h3>
        <img alt="loading" class="loading" src="../assets/images/gear.svg"/>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "Progress",
  data() {
    return {
      showProgress: false,
      message: '',
      eventHandler: (evt) => {
        this.handleEvent(evt)
      },
      closeEventHandler: (evt) => {
        this.handleCloseEvent(evt);
      }
    };
  },
  methods: {
    handleEvent(evt) {
      this.message = evt.detail.message;
      this.showProgress = true;
    },
    handleCloseEvent() {
      this.showProgress = false;
    },
  },
  created() {
    window.addEventListener('showProgressDialog', this.eventHandler);
    window.addEventListener('closeProgressDialog', this.closeEventHandler);
  },
  destroyed() {
    window.removeEventListener('showProgressDialog', this.eventHandler);
    window.removeEventListener('closeProgressDialog', this.closeEventHandler);
  }
}
</script>

<style scoped>
.outer {
  position: fixed;
  z-index: var(--z-toolbar);
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: var(--overlay-color);
}

.inner {
  position: absolute;
  left: 25vw;
  width: 50vw;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.progress {
  text-align: center;
  display: inline-block;
  background-color: var(--bg-color);
  border-radius: 8px;
  padding: 16px;
}

.loading {
  text-align: center;
  animation: rotation 1s infinite linear;
}

@keyframes rotation {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(359deg);
  }
}
</style>