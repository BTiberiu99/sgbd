<template>
  <v-container fluid class="px-0">
    <v-row>
      <v-col cols="12" style="text-align:center;padding:25px;">
        <v-btn :disabled="isLoading" :loading="isLoading" @click="run()">
          Run
        </v-btn>
      </v-col>

      <v-col cols="12" style="text-align:center;padding:25px;">
        <v-textarea
          v-model="sqlRun"
          :clearable="true"
          solo
          label="SQL"
          height="300"
        />
      </v-col>

      <v-col v-if="data" cols="12" style="text-align:center;padding:25px;">
        <v-btn v-if="page > 1" style="margin-right:25px;" :disabled="isLoading" :loading="isLoading" @click="back">
          Back
        </v-btn>
        <v-btn v-if="data.Rows.length >= limit" style="margin-right:25px;" :disabled="isLoading" :loading="isLoading" @click="next">
          Next
        </v-btn>

        <span :key="page" class="text-left">
          Items from  {{ (page - 1) * limit }} to {{ page * limit - (limit - data.Rows.length) }}
        </span>
        <v-simple-table style="margin-top:25px;" dark>
          <template v-slot:default>
            <thead>
              <tr>
                <th v-for="(item,index) in data.Columns" :key="index" class="text-center">
                  {{ item | capitalize }}
                </th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(item,index) in data.Rows" :key="index">
                <td v-for="(item2,index2) in item" :key="index2">
                  {{ item2 }}
                </td>
              </tr>
            </tbody>
          </template>
        </v-simple-table>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { getInstanceQueueMessage } from '@/utils/Queue.js'
import { WAILSINIT } from '@/store/events'
export default {
  name: 'Run',
  filters: {
    capitalize: function (value) {
      if (!value) return ''
      value = value.toString().replace('_', ' ')
      return value.charAt(0).toUpperCase() + value.slice(1)
    }
  },
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
      message: '',
      data: null,
      limit: 15,
      page: 1,
      keep: '',
      isLoading: false
    }
  },

  created () {
    this.$root.$on(WAILSINIT, this.init)
  },
  methods: {
    init () {
      this.isLoading = false
    },
    async run (sql) {
      if (this.isLoading) return
      this.isLoading = true

      var sqlRun = sql || this.sqlRun.replace(/\s\s+/g, ' ')

      if (!sql && sqlRun.toUpperCase().indexOf(this.$static.SQL.SELECT) !== -1) {
        this.page = 1
        this.keep = sqlRun
        sqlRun = this.paginate(sqlRun)
      }

      const rez = await this.$backend.Run({
        run: sqlRun
      })

      getInstanceQueueMessage().addMessage(rez)

      if (rez.data && rez.data.Rows) {
        this.data = rez.data
      } else {
        this.data = null
      }

      this.isLoading = false
    },

    async next () {
      if (this.isLoading) return

      this.page = this.page + 1

      await this.run(this.paginate(this.keep))
    },
    async back () {
      if (this.isLoading) return
      this.page = this.page - 1

      await this.run(this.paginate(this.keep))
    },

    paginate (run) {
      var regexLimit = new RegExp('LIMIT [0-9]*', 'g')
      var regexOffset = new RegExp('OFFSET [0-9]*', 'g')
      var regexLimitSimple = new RegExp('LIMIT', 'g')
      var regexOffsetSimple = new RegExp('OFFSET', 'g')
      return run.toUpperCase().replace(regexLimit, '').replace(regexOffset, '').replace(regexLimitSimple, '').replace(regexOffsetSimple, '') + ` LIMIT ${this.limit} OFFSET ${(this.page - 1) * this.limit}`
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
