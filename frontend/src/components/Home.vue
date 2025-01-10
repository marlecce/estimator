<template>
  <div class="flex flex-col items-center justify-center min-h-screen bg-gray-100">
    <!-- Header -->
    <div class="text-center mb-8">
      <h1 class="text-3xl font-bold text-gray-800 mb-2">Benvenuto nel sistema di stima</h1>
      <p class="text-gray-600">Crea o entra in una stanza per iniziare a stimare i tuoi task!</p>
    </div>

    <!-- Actions Card -->
    <div class="bg-white shadow-md rounded-lg p-6 w-80">
      <button
        @click="goToCreateRoom"
        class="w-full bg-blue-600 text-white font-semibold py-2 rounded-lg mb-4 hover:bg-blue-700 transition"
      >
        Crea una nuova stanza
      </button>

      <div class="mb-4">
        <input
          v-model="roomId"
          placeholder="Inserisci ID stanza"
          class="w-full border rounded-lg p-2 text-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>
      <button
        @click="joinRoom"
        class="w-full bg-green-600 text-white font-semibold py-2 rounded-lg hover:bg-green-700 transition"
      >
        Entra nella stanza
      </button>
    </div>

    <!-- Error Message -->
    <p v-if="errorMessage" class="text-red-600 mt-4">{{ errorMessage }}</p>
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
/* Responsive layout and utility styles */
body {
  font-family: 'Inter', sans-serif;
}

button {
  transition: background-color 0.3s ease;
}

input:focus {
  outline: none;
}
</style>
