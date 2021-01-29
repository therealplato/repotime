const headers = { "Content-Type": "application/json" };
const Timeline = {
  data() {
    return {
      username: "",
      reposAvailable: false,
      selectedRepo: null,
      repos: [],
    };
  },

  created() {
    fetch("/api/username", { headers })
      .then((response) => response.json())
      .then((data) => (this.username = data.login));

    fetch("/api/repositories", { headers })
      .then((response) => response.json())
      .then((data) => {
        this.reposAvailable = true;
        this.repos = data;
      });
  },

  methods: {
    onRepoSelected(event) {
      repoName = event.target.value;
      console.log(repoName);
      let selectedRepo = {};
      for (repo of this.repos) {
        if (repo.full_name === repoName) {
          selectedRepo = repo;
        }
      }
      fetch("/api/set-repository", {
        method: "POST",
        headers,
        body: JSON.stringify(selectedRepo),
      })
        .then((response) => {
          console.log("repository has been selected");
        })
        .catch((e) => {
          console.error(e);
        });
    },
  },
};

const app = Vue.createApp(Timeline);
app.mount("#repotime");
