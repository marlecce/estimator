<template>
  <div class="flex flex-col items-center justify-center min-h-screen bg-gradient-to-r from-blue-50 via-white to-blue-50">
    <!-- Header -->
    <div class="text-center mb-10">
      <h1 class="text-4xl font-extrabold text-gray-800 mb-3">Benvenuto</h1>
      <p class="text-lg text-gray-600">Crea o entra in una stanza per stimare i tuoi task in team!</p>
    </div>

    <!-- Actions Container -->
    <div class="bg-white shadow-xl rounded-2xl p-8 w-full max-w-md space-y-10">
      <!-- Create Room Section -->
      <div class="text-center">
        <h2 class="text-2xl font-bold text-gray-700 mb-4">Crea una nuova stanza</h2>
        <button
          @click="goToCreateRoom"
          class="w-full bg-blue-500 text-white font-bold py-3 rounded-xl hover:bg-blue-600 transition-transform transform hover:scale-105"
        >
          🛠 Crea una nuova stanza
        </button>
      </div>

      <!-- Join Room Section -->
      <div class="border-t pt-6">
        <h2 class="text-2xl font-bold text-gray-700 mb-4 text-center">Entra in una stanza esistente</h2>
        <div class="relative mb-6">
          <input
            v-model="roomId"
            placeholder="Inserisci ID stanza"
            class="w-full border border-gray-300 rounded-xl p-3 pl-10 text-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-400"
          />
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="absolute left-3 top-3.5 w-5 h-5 text-gray-400"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M5 13l4 4L19 7"
            />
          </svg>
        </div>
        <button
          @click="joinRoom"
          class="w-full bg-green-500 text-white font-bold py-3 rounded-xl hover:bg-green-600 transition-transform transform hover:scale-105"
        >
          🚪 Entra nella stanza
        </button>
      </div>
    </div>

    <!-- Error Message -->
    <p v-if="errorMessage" class="text-red-500 font-medium mt-6">{{ errorMessage }}</p>
  </div>
</template>

<script>
export default {
  data() {
    return {
      roomId: '',
      errorMessage: '',
    };
  },
  methods: {
    goToCreateRoom() {
      this.$router.push({ name: 'create-room' });
    },
    joinRoom() {
      if (this.roomId.trim() === '') {
        this.errorMessage = 'Inserisci un ID stanza valido!';
      } else {
        this.errorMessage = '';
        this.$router.push({ name: 'join-room', params: { roomId: this.roomId } });
      }
    },
  },
};
</script>

<style scoped>
body {
  font-family: 'Inter', sans-serif;
}

button {
  transition: background-color 0.3s ease, transform 0.2s ease;
}

input:focus {
  outline: none;
}
</style>
