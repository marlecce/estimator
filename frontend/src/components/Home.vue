<template>
  <div>
    <h1>Benvenuto nel sistema di stima</h1>
    <button @click="goToCreateRoom">Crea una nuova stanza</button>
    <input v-model="roomId" placeholder="Inserisci ID stanza" />
    <button @click="joinRoom">Entra nella stanza</button>
  </div>
</template>

<script>
export default {
  data() {
    return {
      roomId: '',
    };
  },
  methods: {
    goToCreateRoom() {
      this.$router.push({ name: 'create-room' });
    },
    joinRoom() {
      if (this.roomId) {
        this.$router.push({ name: 'join-room', params: { roomId: this.roomId } });
      } else {
        alert('Inserisci un ID stanza valido!');
      }
    }
  }
};
</script>
