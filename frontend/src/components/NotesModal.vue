<template>
<Modal v-show="show" @dialogClose="emitClosed">
  <div class="header">
    <h3>{{taskName}}</h3>
    <button @click="toggleShowNew">{{showNew ? 'Cancel' : 'Add New Note'}}</button>
  </div>
  <div v-if="showNew">
    <p>New Note</p>
    <textarea v-model="newNote"></textarea>
    <button @click="addNewNote">Add</button>
  </div>
  <div v-if="!notes || notes.length === 0">
    <h4>No notes for this task</h4>
  </div>
  <div v-else>
    <div class="note" v-for="note in notes" :key="note.id">
      <NoteView :note="note"></NoteView>
    </div>
  </div>
</Modal>
</template>

<script>
import Modal from "./Modal.vue";
import {AddNote, GetTaskNotes} from "../wailsjs/go/main/ModelController.js";
import {EventsOff, EventsOn} from "../wailsjs/runtime/runtime.js";
import NoteView from "./NoteView.vue";

export default {
  name: "NotesModal",
  components: {NoteView, Modal},
  emits: ["dialogClosed"],
  props: {
    taskId: Number,
    taskName: String,
    show: Boolean,
  },
  data() {
    return {
      notes: [],
      showNew: false,
      newNote: "",
    }
  },
  methods: {
    emitClosed() {
      this.$emit("dialogClosed");
      this.updateEvents();
    },
    updateEvents() {
      GetTaskNotes(this.taskId)
          .then(notes => {
            if (!notes) {
              this.notes = [];
            }
            this.notes = notes;
          })
          .catch(console.errorEvent)
          .then(() => {
            this.newNote = '';
            this.showNew = false;
          });
    },
    toggleShowNew() {
      this.showNew = !this.showNew;
      if (!this.showNew) {
        this.newNote = '';
      }
    },
    addNewNote() {
      AddNote(this.taskId, this.newNote)
          .catch(console.errorEvent);
    },
  },
  watch: {
    taskID: {
      handler() {
        this.updateEvents();
      },
      immediate: true,
    }
  },
  created() {
    EventsOn("notesUpdated", this.updateEvents);
  },
  destroyed() {
    EventsOff("notesUpdated");
  },
}
</script>

<style scoped>
.notes-container .header {
  flex: var(--toolbar-height);
}
</style>