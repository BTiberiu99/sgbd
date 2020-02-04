<template>
  <v-container fluid class="px-0">
    <v-row>
      <!-- Run Migrations -->
      <v-col cols="4" style="text-align:left;padding:25px;">
        <v-btn style="margin-left:25px;" @click="runMigrations">
          Run Migrations
        </v-btn>
      </v-col>

      <!-- Run User SQL-->
      <v-col cols="6" offset="2" style="text-align:left;padding:25px;">
        <v-btn :disabled="isLoading" :loading="isLoading" style="margin-right:50px;" @click="runSql">
          Run
        </v-btn>
      </v-col>

      <!-- SQL Input -->
      <v-col cols="12" style="text-align:center;padding:25px;">
        <v-textarea
          v-model="sqlRun"
          :clearable="true"
          solo
          label="SQL"
          height="300"
        />
      </v-col>

      <!-- Data from SELECT  -->
      <v-col v-if="data" cols="12" style="text-align:center;padding:25px;">
        <!-- Pagination -->
        <v-btn v-if="page > 1" style="margin-right:25px;" :disabled="isLoading" :loading="isLoading" @click="back">
          Back
        </v-btn>
        <v-btn v-if="data.Rows.length >= limit" style="margin-right:25px;" :disabled="isLoading" :loading="isLoading" @click="next">
          Next
        </v-btn>

        <span :key="page" class="text-left">
          Items from  {{ (page - 1) * limit }} to {{ page * limit - (limit - data.Rows.length) }}
        </span>

        <!-- Data -->
        <v-simple-table style="margin-top:25px;" dark>
          <template v-slot:default>
            <!-- Head of tabel -->
            <thead>
              <tr>
                <th v-for="(item,index) in data.Columns" :key="index" class="text-center">
                  {{ item | capitalize }}
                </th>
              </tr>
            </thead>

            <!-- Body of table -->
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
import { WAILSINIT, REFRESHTABLES } from '@/store/events'

const regexLimit = new RegExp('LIMIT [0-9]*', 'g')
const regexOffset = new RegExp('OFFSET [0-9]*', 'g')
const regexLimitSimple = new RegExp('LIMIT', 'g')
const regexOffsetSimple = new RegExp('OFFSET', 'g')
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
      // Data
      data: null,
      limit: 15,
      page: 1,
      keepSql: '',

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
    // run sql
    async run (sql) {
      if (this.isLoading) return
      this.isLoading = true

      var sqlRun = sql.replace(/\s\s+/g, ' ')

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

    // run user input sql
    async runSql () {
      var sqlRun = this.sqlRun
      if (sqlRun.toUpperCase().indexOf(this.$static.SQL.SELECT) !== -1) {
        this.page = 1
        this.keepSql = sqlRun
        sqlRun = this.paginate(sqlRun)
      }

      this.run(sqlRun)
    },

    // go to the next page
    next () {
      if (this.isLoading) return

      this.page = this.page + 1

      this.run(this.paginate(this.keepSql))
    },
    // go back one page
    back () {
      if (this.isLoading) return
      this.page = this.page - 1

      this.run(this.paginate(this.keepSql))
    },

    // paginate the query
    paginate (run) {
      return run.toUpperCase()
        .replace(regexLimit, '')
        .replace(regexOffset, '')
        .replace(regexLimitSimple, '')
        .replace(regexOffsetSimple, '') + ` LIMIT ${this.limit} OFFSET ${(this.page - 1) * this.limit}`
    },
    async runMigrations () {
      await this.$migrations.run()
      this.$emit(REFRESHTABLES)
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
