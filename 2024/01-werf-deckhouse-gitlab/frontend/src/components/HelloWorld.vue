<template>
  <div class="container mt-5">
    <form class="row g-3">
      <div class="col-auto">
        <div class="input-group mb-3">
          <span class="input-group-text" id="name">Name</span>
          <input type="text" class="form-control" placeholder="Name" name="name" v-model="name">
        </div>
      </div>
      <div class="col-auto">
        <div class="input-group mb-3">
          <span class="input-group-text" id="message">Message</span>
          <input type="text" class="form-control" placeholder="Message" name="message" v-model="message">
        </div>
      </div>
      <div class="col-auto">
        <button v-on:click="send" type="submit" class="btn btn-primary mb-3">Send</button>
      </div>
    </form>
    <div class="container mt-5">
      <h2 class="text-center">Messages from users</h2>
      <button v-on:click="get" type="button" class="btn btn-outline-success mb-2">Проверить записи</button>

      <div class="alert alert-warning" role="alert" v-if="error">
        Nothing to show!
      </div>

      <table class="table" v-else>
        <thead>
        <th>Name</th>
        <th>Message</th>
        </thead>
        <tbody>
        <tr v-for="talker in talkers.Messages" :key="talker.Name">
          <td>{{ talker.Name }}</td>
          <td>{{ talker.Message }}</td>
        </tr>
        </tbody>
      </table>
    </div>
  </div>

</template>

<script>
import axios from 'axios';
export default {
  name: 'HelloWorld',
  props: {
    msg: String
  },
  data() {
    return {
      name: '',
      message: '',
      talkers: [],
      error: true
    }
  },
  methods: {
    send: function()
    {
      axios.get('/api/remember', {
        params: {
          name: this.name,
          message: this.message
        }
      }).catch(error => alert(error));
    },
    get: function () {
      axios.get("/api/say")
          .then(response => {
            this.talkers = response.data;
            this.error = false
        })
    }
  }
}
</script>

<style></style>
