<template>
  <v-container fluid class="px-0">
      <v-form
      ref="form"
      lazy-validation
    >
        <v-row>
          <v-col class="d-flex" cols="12" sm="4" style="padding:35px;">
             <v-text-field
              v-model="form.host"
              :rules="rules.host"
              label="Host"
              required
            ></v-text-field>
          </v-col>

          <v-col class="d-flex" cols="12" sm="4" style="padding:35px;">
               <v-text-field
              v-model="form.port"
              :rules="rules.port"
              label="Port"
              required
            ></v-text-field>
          </v-col>

          <v-col class="d-flex" cols="12" sm="4" style="padding:35px;">
             <v-text-field
              v-model="form.user"
              :rules="rules.user"
              autocomplete="off"
              label="User"
              required
            ></v-text-field>
          </v-col>

          <v-col class="d-flex" cols="12" sm="4" style="padding:35px;">
               <v-text-field
              v-model="form.password"
              :append-icon="showPass ? 'mdi-eye' : 'mdi-eye-off'"
              :type="showPass ? 'text' : 'password'"
              autocomplete="off"
              label="Password"
              @click:append="showPass = !showPass"
            ></v-text-field>
          </v-col>

          <v-col class="d-flex" cols="12" sm="4" style="padding:35px;">
               <v-text-field
              v-model="form.database"
              :rules="rules.database"
              label="Database"
              required
            ></v-text-field>
          </v-col>
        </v-row>
        <v-row>
          <v-col class="d-flex" cols="12" sm="4" :offset="3" style="padding:35px;">
          <v-btn
            :color="buttonColor"
            @click="createConnection"
            :loading="isLoading"
            :disbaled="isLoading"
          >
            Create Connection
          </v-btn>
          </v-col>
        </v-row>
      </v-form>
      
  </v-container>
</template>


<script>

function ValidateString (v) {
  return !!v || "Campul este necesar"
}

function ValidateNumber (v) {
  return (!isNaN(parseInt(v)) && parseInt(v) > 0) || "Trebuie sa introduceti un numar mai mare ca 0"
}
import { CREATECONNECTION, ENDCREATECONNECTION, WAILSINIT } from '@/store/events'
export default {
  name: 'Connection',
  components: {
  },
  props: {
  },

  data () {
    var model = {
      host: "",
      port: "",
      user: "",
      password: "",
      database: "",
    }
    return {
      valid: true,

      showPass: false,
      isLoading: false,

      buttonColor: "normal",

      model,
      form: {
        ...model,
      },

      rules: {
        host: [ValidateString],
        port: [ValidateString, ValidateNumber],
        user: [ValidateString],
        database: [ValidateString]
      }
    };
  },

  created () {
    this.$root.$on(WAILSINIT, this.resetLoadings)
    this.$root.$on(ENDCREATECONNECTION, this.finishCreateConnection)
  },

  methods: {
    resetLoadings () {
      this.isLoading = false
    },
    createConnection () {
      if (this.isLoading) return
      if (!this.$refs.form.validate()) return

      this.isLoading = true

      this.form.port = parseInt(this.form.port)

      this.$emit(CREATECONNECTION, this.form)

    },


    finishCreateConnection (r) {

      this.buttonColor = r.type

      if (r.type === "success") {
        this.form = { ...this.model }
      }

      this.isLoading = false

    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
