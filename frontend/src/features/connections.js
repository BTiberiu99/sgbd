import { getInstanceQueueMessage } from "@/utils/Queue.js";
import { ENDCREATECONNECTION } from "@/store/events";
export const AppConnections = function () {
  return {
    data () {
      return {
        connections: [],
        connectionActive: "",
        connectionIsChanging: false,
        isDeletingConnection: false,
        delete: false,
      };
    },
    methods: {
      takeConnections () {
        this.connectionIsChangingRun(async () => {
          const rez = await this.$backend.GetConnections();
          this.connections = rez;
        });
      },
      switchConnection (conn) {
        this.connectionIsChangingRun(async () => {
          const rez = await this.$backend.SwitchConnection(conn);

          this.connectionActive = rez.Index;
        });
      },

      continueDelete () {

        this.delete = true
        this.isDeletingConnection = false
      },

      cancelDelete () {
        this.delete = false
        this.isDeletingConnection = false

      },
      async deleteConnection (conn) {

        this.delete = false
        this.isDeletingConnection = true

        await this.$sync(() => { return !this.isDeletingConnection })

        if (!this.delete) return

        this.connectionIsChangingRun(async () => {
          const rez = await this.$backend.RemoveConnection(conn);

          this.connections = rez.Data;

          getInstanceQueueMessage().addMessage(rez);
        });
      },

      createConnection (conn) {
        var rez
        this.connectionIsChangingRun(async () => {
          rez = await this.$backend.CreateConnection(conn);

          if (rez.data) {
            var add = true;
            this.connections.map(conn => {
              if (conn.Index == rez.data.Index) {
                add = false;
              }
            });
            if (add) {
              this.connections.push(rez.data);
            }

            this.connectionActive = rez.data.Index;
          }

          getInstanceQueueMessage().addMessage(rez);


        }, () => {
          this.$root.$emit(ENDCREATECONNECTION, rez);
        });
      },

      async connectionIsChangingRun (call, final) {

        if (typeof call !== 'function') {
          throw new Error('First parameter must be an callback')
        }

        this.connectionIsChanging = true;
        var start = Date.now();
        try {
          await call();
        } catch (e) {
          // console.log(e)
        }
        // the event you'd like to time goes here:
        var end = Date.now();
        var elapsed = end - start; // time in milliseconds
        setTimeout(() => {
          this.connectionIsChanging = false
          if (typeof final === 'function') {
            final()
          }

        }, 1000 - elapsed)

      }
    }
  };
};
