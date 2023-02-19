<template>
  <div class="note-header">
    <p>{{dateTimeString}}</p>
    <span class="material-icons delete" @click="deleteNote">delete</span>
  </div>
  <textarea v-if="editing" @blur="stopEdit" v-model="note.text">{{ note.text }}</textarea>
  <p class="static-note" v-else @click="doEdit">{{note.text}}</p>
</template>

<script>
import {DeleteNote, UpdateNote} from "../wailsjs/go/main/ModelController.js";

export default {
  name: "NoteView",
  props: {
    note: Object,
  },
  data() {
    return {
      editing: false,
    }
  },
  methods: {
    doEdit() {
      this.editing = true;
    },
    stopEdit() {
      this.editing = false;
      if (this.note.text.length === 0) {
        UpdateNote(this.note.id, this.note.text)
            .catch(console.errorEvent);
      }
    },
    deleteNote() {
      DeleteNote(this.note.id)
          .catch(console.errorEvent);
    },
  },
  computed: {
    dateTimeString() {
      let date = new Date(this.note.created);
      return `${date.toLocaleDateString()} ${date.toLocaleTimeString()}`;
    }
  },
}
</script>

<style scoped>
.delete:hover {
  color: var(--fg-danger);
}
.note-header {
  display: flex;
  flex-direction: row;
  align-items: center;
  border-top: 1px solid var(--bg-color-light);
}
.note-header > p {
  font-weight: bold;
  margin: 8px 0 4px 0;
  padding: 0;
}
.note-header > .delete {
  padding-left: 8px;
  cursor: pointer;
}
.static-note {
  font-weight: bold;
  cursor: pointer;
  margin: 0 0 8px 0;
}
.static-note:hover {
  color: var(--fg-highlight);
}
</style>