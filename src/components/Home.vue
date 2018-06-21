<template>
  <div class='hello'>
    <h1>{{ msg }}</h1>
    <h2 v-if="authenticated">
      You are logged in!
    </h2>
    <h2 v-if="!authenticated">
      You are NOT logged in! Please
      <a @click="auth.login()">Log In</a> to continue.
    </h2>
    <div>
      <button v-if="!authenticated" @click="auth.login">login</button>
      <button v-if="authenticated" @click="auth.logout">logout</button>
    </div>
    <button @click="apiPublic">public</button>
    <button @click="apiPrivate">private</button>
    <button @click="apiPrivateScoped">private-scoped</button>
  </div>
</template>

<script>
import axios from "axios";
export default {
  name: "HelloWorld",
  props: ["auth"],
  data() {
    return {
      msg: "Message will display here"
    };
  },
  computed: {
    authenticated() {
      return this.$store.getters.authenticated;
    }
  },
  methods: {
    apiPublic: async function() {
      try {
        let res = await axios.get("http://localhost:50000/public");
        this.msg = res.data;
      } catch (e) {
        this.msg = e.message;
      }
    },
    apiPrivate: async function() {
      try {
        let res = await axios.get("http://localhost:50000/private", {
          headers: {
            Authorization: `Bearer ${localStorage.getItem("access_token")}`
          }
        });
        this.msg = res.data;
      } catch (e) {
        this.msg = e.message;
      }
    },
    apiPrivateScoped: async function() {
      try {
        let res = await axios.get("http://localhost:50000/private-scoped", {
          headers: {
            Authorization: `Bearer ${localStorage.getItem("access_token")}`
          }
        });
        this.msg = res.data;
      } catch (e) {
        this.msg = e.message;
      }
    }
  }
};
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style scoped>
h1,
h2 {
  font-weight: normal;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
button {
  margin: 10px 0;
  padding: 10px;
}
</style>