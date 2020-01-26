<template>
  <v-container fluid class="px-0">
    <v-row>
      <v-col cols="12" style="text-align:center;padding:25px;">
        <v-btn @click="run">
          Run
        </v-btn>
      </v-col>

      <v-col cols="12" style="text-align:center;padding:25px;">
        <v-textarea
          :clearable="true"
          solo
          label="SQL"
          height="300"

          :value="sqlRun"
        />
      </v-col>

      <v-col v-if="message" cols="12" style="text-align:center;padding:25px;">
        <v-textarea
          :clearable="true"
          height="300"
          solo
          label="Message"
          :value="message"
        />
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { getInstanceQueueMessage } from '@/utils/Queue.js'
import { WAILSINIT } from '@/store/events'
export default {
  name: 'Run',
  props: {
    tables: {
      type: Array,
      required: true
    },
    isLoadingTables: {
      type: Boolean,
      required: true
    }
  },
  data () {
    return {
      sqlRun: '',
      message: ''
    }
  },

  created () {
    this.$root.$on(WAILSINIT, this.init)
  },
  methods: {
    init () {

    },
    async run () {
      const rez = await this.$backend.Run({
        run: this.sqlRun
      })
      getInstanceQueueMessage().addMessage(rez)
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1 {
  margin-top: 2em;
  position: relative;
  min-height: 5rem;
  width: 100%;
}
a:hover {
  font-size: 1.7em;
  border-color: blue;
  background-color: blue;
  color: white;
  border: 3px solid white;
  border-radius: 10px;
  padding: 9px;
  cursor: pointer;
  transition: 500ms;
}
a {
  font-size: 1.7em;
  border-color: white;
  background-color: #121212;
  color: white;
  border: 3px solid white;
  border-radius: 10px;
  padding: 9px;
  cursor: pointer;
}
</style>
