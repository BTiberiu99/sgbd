<template>
  <v-container fluid class="px-0">
    <v-row style="padding:25px;">
      <v-col cols="4" class="TablesList">
        <div style="width:200px;">
          Tables
          <v-select v-model="table"
                    :items="tables"
                    label="Table"
                    dense
                    return-object
                    persistent-hint
                    single-line
                    item-text="Name"
                    :hint="table && table.Hint"
          />
        </div>
      </v-col>
      <v-col cols="8" style="text-align:center;" />
      <v-col cols="8" offset="2" style="text-align:center;">
        <Table v-if="table" :table="table" style="margin:0 0 20px 0 ;" />
        <v-btn v-if="table && !table.HasOneNotNull" style="margin:0 40px 0 0;" @click="addNotNull(table)">
          Fix Not Null
        </v-btn>

        <v-btn v-if="table && (!table.HasPrimaryKey || !table.HasCorrectPrimaryKey)" @click="addPrimaryKey(table)">
          Fix Primary Key
        </v-btn>
      </v-col>
    </v-row>

    <v-dialog
      v-model="dialog"
      max-width="350"
    >
      <v-card v-if="table">
        <v-card-title class="headline">
          {{ title }}
        </v-card-title>

        <v-card-text :key="table.vueKey">
          <!-- Not NULL check -->
          <template v-if="action === 'notnull'">
            <template v-if="!badConstruction">
              <h3>Columns</h3>
              <template v-for="(column,index) in table.Columns">
                <v-checkbox v-if="!column.HasNotNull && !column.HasPrimaryKey && column.WithoutNULL"
                            :key="index"
                            v-model="columnsNotNull"
                            :label="`${column.Name} - ${column.Type}`"
                            :value="{index:index, column:column}"
                />
              </template>
            </template>
            <template v-else>
              <p>
                Tabelul  {{ table.Name }} are intrari NULL in toate coloanele ,
                astfel neputand sa se adauge o constrangere not null asupra niciunei coloane,
                stergeti toate intrarile NULL ale unei coloane ca sa puteti adauga o constrangere not null
              </p>
            </template>
          </template>

          <template v-if="action === 'primarykey'">
            <template v-if="!table.HasPrimaryKey">
              Nu exista nicio cheie primara in tabel , introduceti numele mai jos ca sa creati una surogat
            </template>
            <template v-else>
              Cheia curenta primara nu este formata corect, introduceti numele mai jos ca sa creati una surogat
            </template>
            <v-text-field
              v-model="primaryKeyName"
              label="Nume"
              placeholder=""
              @input="checkPrimaryKeyName"
            />
            <div v-if="isTakenName" style="color: #ff3c00ce;">
              Exista deja o coloana in tabel cu acest nume!
            </div>
          </template>
        </v-card-text>

        <v-card-actions>
          <v-spacer />

          <v-btn
            color="red darken-1"
            text
            @click="cancelAction"
          >
            Cancel
          </v-btn>

          <v-btn
            color="green darken-1"
            text
            @click="continueAction"
          >
            Continue
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script>
import { WAILSINIT } from '@/store/events'
import Column from '@/utils/Column'

import { getInstanceQueueMessage } from '@/utils/Queue.js'
export default {
  name: 'CheckIntegriy',
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
      table: null,
      dialog: false,
      columnsNotNull: [],
      continue: false,
      wait: false,
      action: '',
      primaryKeyName: '',
      isTakenName: false,
      validate: {
        notnull: () => true,
        primarykey: () => !this.isTakenName
      }
    }
  },

  computed: {
    badConstruction () {
      let count = 0

      if (this.table !== null) {
        this.table.Columns.map(column => {
          if (!column.HasPrimaryKey && column.WithoutNULL) {
            count++
          }
        })
      }

      return count < 1
    },
    title () {
      switch (this.action) {
        case 'notnull':
          return 'Add NOT NULL'
        default:
          return ''
      }
    }
  },

  created () {
    this.$root.$on(WAILSINIT, this.init)
  },
  methods: {
    init () {

    },
    getMessage () {
      // var self = this;

    },
    print () {
      console.log(arguments)
    },

    cancelAction () {
      this.dialog = false
      this.continue = false
      this.wait = false
    },
    continueAction () {
      if (!this.validate[this.action]()) return
      this.dialog = false
      this.continue = true
      this.wait = false
    },
    async addPrimaryKey (table) {
      this.action = 'primarykey'
      this.wait = true
      this.dialog = true
      var rez
      await this.$sync(() => !this.wait)
      if (!table.HasPrimaryKey) {
        rez = await this.$backend.AddPrimaryKey(table.Name, this.primaryKeyName)
      } else {
        rez = await this.$backend.FixPrimaryKey(table.Name, this.primaryKeyName)
      }

      if (rez) {
        getInstanceQueueMessage().addMessage(rez)

        table.Columns = rez.data.Columns.map(col => {
          return new Column(col)
        })
      }
    },
    async addNotNull (table) {
      this.action = 'notnull'
      this.wait = true
      this.dialog = true
      this.columnsNotNull = []

      await this.$sync(() => !this.wait)

      if (!this.continue) return

      var items = await Promise.all(this.columnsNotNull.map(async (item) => {
        var rez = item
        if (!item.column.HasNotNull) {
          const response = await this.$backend.AddNotNull(table.Name, item.column)
          if (response.type === 'success') {
            rez.modified = true
            rez.column = response.data
          }

          getInstanceQueueMessage().addMessage(response)
        }

        return rez
      }))

      items.map(item => {
        if (item.modified) {
          table.Columns[item.index] = new Column(item.column)
        }
      })
    },

    checkPrimaryKeyName () {
      let i
      this.isTakenName = false
      for (i in this.table.Columns) {
        if (this.table.Columns[i].Name === this.primaryKeyName) {
          this.isTakenName = true
          break
        }
      }
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss">
.TablesList {
  text-align: center;
  .v-text-field__details {
    .v-messages__message.message-transition-enter-to {
      color: #ff3c00ce;
    }
  }
}
</style>
