<template>
  <div class="min-h-screen bg-gray-50 flex flex-col items-center justify-center px-4">
    <!-- Card Container -->
    <div class="bg-white shadow-lg rounded-lg p-6 w-full max-w-md">
      <!-- Title -->
      <h1 class="text-2xl font-bold text-gray-800 text-center mb-6">
        Create a Room
      </h1>

      <!-- Room Name Input -->
      <div class="mb-4">
        <label class="block text-sm font-medium text-gray-700 mb-2">
          Room Name
        </label>
        <input
          v-model="name"
          type="text"
          placeholder="Enter room name"
          class="w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
        />
      </div>

      <!-- Host Name Input -->
      <div class="mb-4">
        <label class="block text-sm font-medium text-gray-700 mb-2">
          Your Name
        </label>
        <input
          v-model="host_name"
          type="text"
          placeholder="Enter your name"
          class="w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
        />
      </div>

      <!-- Estimation Type Dropdown -->
      <div class="mb-6">
        <label class="block text-sm font-medium text-gray-700 mb-2">
          Estimation Type
        </label>
        <select
          v-model="estimationType"
          class="w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
        >
          <option value="hours">Hours</option>
          <option value="days">Days</option>
          <option value="story_points">Story Points</option>
        </select>
      </div>

      <!-- Create Room Button -->
      <button
        @click="createRoom"
        class="w-full bg-blue-600 text-white font-medium py-2 rounded-md shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition"
      >
        Create Room
      </button>

      <!-- Room Created Section -->
      <div v-if="roomId" class="mt-6">
        <p class="text-center text-sm font-medium text-green-600 mb-4">
          Room created successfully!
        </p>
        <label class="block text-sm font-medium text-gray-700 mb-2">
          Share this link with participants:
        </label>
        <div class="flex items-center gap-2">
          <input
            v-model="shareLink"
            readonly
            class="flex-1 px-4 py-2 border border-gray-300 rounded-md shadow-sm"
          />
          <button
            @click="copyLink"
            class="bg-gray-200 text-gray-600 px-3 py-2 rounded-md shadow-sm hover:bg-gray-300 transition"
          >
            Copy
          </button>
        </div>
        <p
          v-if="linkCopied"
          class="mt-2 text-sm text-green-500 font-medium text-center"
        >
          Link copied to clipboard!
        </p>

        <!-- Go to Room Button -->
        <button
          @click="goToRoom"
          class="w-full mt-4 bg-blue-600 text-white font-medium py-2 rounded-md shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition"
        >
          Go to Room
        </button>
      </div>
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
      estimationType: "hours",
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
        const response = await apiClient.post("/rooms", {
          name: this.name,
          host_name: this.host_name,
          estimation_type: this.estimationType,
        });
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
