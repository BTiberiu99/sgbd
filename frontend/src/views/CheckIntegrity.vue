<template>
  <v-container fluid class="px-0">
    <v-row style="padding:25px;">
      <!-- Drop List with tables -->
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

      <!-- Offset -->
      <v-col cols="8" style="text-align:center;" />

      <!-- Current Tabel -->
      <v-col cols="8" offset="2" style="text-align:center;">
        <Table v-if="table" :table="table" style="margin:0 0 20px 0 ;" />

        <!-- Buttons -->

        <!-- NOT NULL -->
        <v-btn v-if="table && !table.HasOneNotNull" style="margin:0 40px 0 0;" @click="addNotNull(table)">
          Fix Not Null
        </v-btn>

        <!-- PRIMARY KEY -->
        <v-btn v-if="table && (!table.HasPrimaryKey || !table.HasCorrectPrimaryKey)" @click="addPrimaryKey(table)">
          Fix Primary Key
        </v-btn>
      </v-col>
    </v-row>

    <!-- Dialog with the user -->
    <v-dialog
      v-model="dialog"
      max-width="350"
      @click:outside="cancelAction"
    >
      <v-card v-if="table">
        <!-- Title -->
        <v-card-title class="headline">
          {{ title }}
        </v-card-title>

        <v-card-text :key="table.vueKey">
          <!-- Not NULL check -->

          <!-- NOT NULL CASE -->
          <template v-if="action === 'notnull'">
            <template v-if="badConstruction === 0">
              <h3>Columns</h3>

              <!-- All Columns that can have NOT NULL CONSTRAINT -->
              <template v-for="(column,index) in table.Columns">
                <v-checkbox v-if="!column.HasNotNull && !column.HasPrimaryKey && column.WithoutNULL"
                            :key="index"
                            v-model="columnsNotNull"
                            :label="`${column.Name} - ${column.Type}`"
                            :value="{index:index, column:column}"
                />
              </template>
            </template>

            <!-- Bad Construction -->
            <template v-else-if="badConstruction === 1">
              <p>
                Tabelul  {{ table.Name }} are intrari NULL in toate coloanele care nu au constrangere de cheie primara,
                astfel neputand sa se adauge o constrangere not null asupra niciunei coloane,
                completati toate intrarile NULL ale unei coloane ca sa puteti adauga o constrangere not null
              </p>
            </template>

            <template v-else-if="badConstruction === 2">
              <p>
                Tabelul  {{ table.Name }} nu are nicio coloana inafara de coloane cu constrangeri de cheie primara.
                Adaugati o noua coloana fara intrari NULL
              </p>
            </template>
          </template>

          <!-- PRIMARY KEY CASE -->
          <template v-if="action === 'primarykey'">
            <!-- Text -->
            <template v-if="!table.HasPrimaryKey">
              Nu exista nicio cheie primara in tabel , introduceti numele mai jos ca sa creati una surogat
            </template>
            <template v-else>
              Cheia curenta primara nu este formata corect, introduceti numele mai jos ca sa creati una surogat
            </template>

            <!-- Input Name Of Primary Key -->
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

        <!-- Actions -->
        <v-card-actions>
          <!-- Cancel -->
          <v-spacer />

          <v-btn
            color="red darken-1"
            text
            @click="cancelAction"
          >
            Cancel
          </v-btn>

          <!-- Continue -->
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
import { WAILSINIT, REFRESHTABLES } from '@/store/events'

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

      action: '',
      dialog: false,
      wait: false,
      continue: false,

      table: null,
      columnsNotNull: [],

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
      let countPrimaryKeys = 0
      if (this.table !== null) {
        this.table.Columns.map(column => {
          if (!column.HasPrimaryKey && column.WithoutNULL) {
            count++
          } else if (column.HasPrimaryKey) {
            countPrimaryKeys++
          }
        })
      }

      if (countPrimaryKeys === this.table.Columns.length) {
        return 2
      }

      return count < 1 ? 1 : 0
    },
    title () {
      switch (this.action) {
        case 'notnull':
          return 'Add NOT NULL'
        case 'primarykey':
          return 'Fix Primary Key'
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

    // cancel action taken by the user
    cancelAction () {
      this.dialog = false
      this.continue = false
      this.wait = false
    },

    // continue action taken by the user if is valid
    continueAction () {
      if (!this.validate[this.action]()) return
      this.dialog = false
      this.continue = true
      this.wait = false
    },

    // solve the primary key normalization
    async addPrimaryKey (table) {
      this.action = 'primarykey'
      this.wait = true
      this.dialog = true
      var rez
      await this.$sync(() => !this.wait)

      if (!this.continue) return

      if (!table.HasPrimaryKey) {
        rez = await this.$backend.AddPrimaryKey(table.Name, this.primaryKeyName)
      } else {
        rez = await this.$backend.FixPrimaryKey(table.Name, this.primaryKeyName)
      }

      if (rez) {
        getInstanceQueueMessage().addMessage(rez)

        if (rez.data) {
          this.$root.$emit(REFRESHTABLES, rez.data)
        }
      }
    },
    // solve the constraints of existence normalization
    async addNotNull (table) {
      this.action = 'notnull'
      this.wait = true
      this.dialog = true
      this.columnsNotNull = []

      await this.$sync(() => !this.wait)

      if (!this.continue) return

      await Promise.all(this.columnsNotNull.map(async (item) => {
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

      getInstanceQueueMessage().addMessage(await this.$backend.ResetTables())

      this.$root.$emit(REFRESHTABLES, true)
    },

    // check the primaryKeyName to not be taken
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
