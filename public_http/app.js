const Timeline = {
  data() {
    return {
      username: "@therealplato",
    };
  },

  created() {
    // GET request using fetch with set headers
    const headers = { "Content-Type": "application/json" };
    fetch("/api/username", { headers })
      .then((response) => response.json())
      .then((data) => (this.username = data.username));
  },
};

const app = Vue.createApp(Timeline);
app.mount("#repotime");
