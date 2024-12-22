<template>
  <div id="app">
    <header>
      <h1>Estimation App</h1>
    </header>
    <router-view />
  </div>
</template>

<script>
export default {
  name: "App",
};
</script>

<style>
#app {
  text-align: center;
}

header {
  background-color: #282c34;
  padding: 20px;
  color: white;
}

h1 {
  margin: 0;
}
</style>
