<script setup lang="ts">
  defineProps<{ msg: string }>()
</script>
<template>
    <div class="home">
      <h1>Welcome to the Estimation App</h1>
      <div class="actions">
        <router-link to="/rooms/create">
          <button>Create a Room</button>
        </router-link>
        <router-link to="/rooms/join">
          <button>Join a Room</button>
        </router-link>
      </div>
    </div>
</template>

<style scoped>
.home {
  text-align: center;
  margin-top: 50px;
}

.actions button {
  margin: 10px;
  padding: 10px 20px;
  font-size: 16px;
  cursor: pointer;
}
</style>
  