const headers = { "Content-Type": "application/json" };
const combineCommitsAndIssues = (commits, issues) => {
  let combined = [];
  for (commit of commits) {
    combined.push({
      time: commit.commit.committer.date,
      type: "commit",
      message: commit.commit.message,
      url: commit.html_url,
    });
  }
  for (issue of issues) {
    let item = {};
    if (issue.issue.pull_request && issue.issue.pull_request.url) {
      if (issue.event == "closed") {
        item.type = "pr_closed";
        item.url = issue.issue.pull_request.html_url;
      } else if (issue.event == "merged") {
        item.type = "pr_merged";
        item.url = issue.issue.pull_request.html_url;
      } else {
        continue;
      }
    } else {
      if (issue.event == "closed") {
        item.type = "issue_closed";
        item.url = issue.issue.html_url;
      } else {
        continue;
      }
    }
    item.message = issue.issue.title;
    item.time = issue.created_at;
    combined.push(item);
  }
  combined.sort((a, b) => {
    return new Date(a.time) - new Date(b.time);
  });
  return combined;
};
const Timeline = {
  data() {
    return {
      username: "",
      reposAvailable: false,
      commitsAvailable: false,
      issuesAvailable: false,
      timelineDataAvailable: false,
      selectedRepo: null,
      repos: [],
      commits: [],
      issues: [],
    };
  },
  computed: {
    timelineData() {
      return combineCommitsAndIssues(this.commits, this.issues);
    },
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
      this.timelineDataAvailable = false;
      repoName = event.target.value;
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

      fetch("/api/commits", { headers })
        .then((response) => response.json())
        .then((data) => {
          this.commitsAvailable = true;
          this.commits = data;
        });

      fetch("/api/issues", { headers })
        .then((response) => response.json())
        .then((data) => {
          this.issuesAvailable = true;
          this.issues = data;
        });
      this.timelineDataAvailable = true;
    },
  },
};

const app = Vue.createApp(Timeline);
app.mount("#repotime");
