<template>
  <div>
    <h1>Create a Room</h1>
    <p>Enter a name for the room.</p>
    <input v-model="name" placeholder="Room Name" />
    <button @click="createRoom">Create Room</button>

    <div v-if="roomId">
      <p>Room created! Share this link with participants:</p>
      <input v-model="shareLink" readonly />
      <button @click="copyLink">Copy Link</button>
      <p v-if="linkCopied" style="color: green;">Link copied to clipboard!</p>
    </div>
  </div>
</template>

<script>
import { useRouter } from "vue-router";
import apiClient from "../api-client";

export default {
  data() {
    return {
      name: "",
      roomId: null,
      shareLink: "",
      linkCopied: false,
    };
  },
  methods: {
    async createRoom() {
      try {
        const response = await apiClient.post("/rooms", { name: this.name });
        this.roomId = response.data.room_id;

        this.shareLink = `${window.location.origin}/rooms/${this.roomId}/join`;
        this.linkCopied = false;
      } catch (error) {
        console.error("Error creating room:", error);
      }
    },
    copyLink() {
      navigator.clipboard.writeText(this.shareLink).then(() => {
        this.linkCopied = true;
      });
    },
  },
};
</script>
