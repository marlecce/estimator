<template>
  <div>
    <h1>Create a Room</h1>
    <p>Enter a name for the room.</p>
    <input v-model="name" placeholder="Room Name" />
    <p>Enter your name.</p>
    <input v-model="host_name" placeholder="Host Name" />
    <button @click="createRoom">Create Room</button>

    <div v-if="roomId">
      <p>Room created! Share this link with participants:</p>
      <input v-model="shareLink" readonly />
      <button @click="copyLink">Copy Link</button>
      <p v-if="linkCopied" style="color: green;">Link copied to clipboard!</p>
      <button @click="goToRoom">Go to Room</button>
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
      hostId: "",
      host_name: "",
      shareLink: "",
      linkCopied: false,
    };
  },
  setup() {
    const router = useRouter();
    return { router };
  },
  methods: {
    async createRoom() {
      try {
        const response = await apiClient.post("/rooms", { name: this.name, host_name: this.host_name });
        this.roomId = response.data.room_id;
        this.hostId = response.data.host.id;

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
    goToRoom() {
      this.router.push(`/rooms/${this.roomId}?participantId=${this.hostId}`);
    },
  },
};
</script>
