<template>
  <div :key="table.keyVue">
    <!-- Title -->
    <h3 class="Table__name">
      Table : {{ table.Name }}

      <v-tooltip v-if="!table.IsSafe" top>
        <template v-slot:activator="{ on }">
          <v-icon class="danger" v-on="on">
            mdi-close-octagon
          </v-icon>
        </template>
        <span class="danger">{{ table.Hint }}</span>
      </v-tooltip>
    </h3>
    <v-expansion-panels :key="table.keyVue">
      <!-- Columns -->
      <v-expansion-panel
        v-for="(column,indexColumn) in table.Columns"
        :key="indexColumn +'_'+ column.keyVue"
        class="Column"
        :class="[{'Column--active':column.Constraints.length>0}]"
      >
        <!-- Name && Type -->
        <v-expansion-panel-header class="Column__header">
          <div>
            <span>
              Column name: <b>{{ column.Name }}</b>
            </span>
            |
            <span>
              Column type:<b> {{ column.Type }}</b>
            </span>
          </div>
        </v-expansion-panel-header>

        <!-- Constraints -->
        <v-expansion-panel-content v-if="column.Constraints.length>0">
          <h4>Constraints</h4>
          <v-simple-table dark>
            <template v-slot:default>
              <thead>
                <tr>
                  <th class="text-left">
                    Name
                  </th>
                  <th class="text-left">
                    Type
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(item,index) in column.Constraints" :key="index+'_'+ item.keyVue">
                  <td>{{ item.Name }}</td>
                  <td>{{ item.Type }}</td>
                </tr>
              </tbody>
            </template>
          </v-simple-table>
        </v-expansion-panel-content>
      </v-expansion-panel>
    </v-expansion-panels>
  </div>
</template>

<script>
import Table from '@/utils/Table'
export default {
  name: 'Tables',

  props: {
    table: {
      type: Table,
      required: true
    }

  }

}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss" >
.Table {
  padding: 20px;

  &__name {
    margin: 0 0 10px 0;
  }
}
.Column {
  pointer-events: none;
  .v-icon {
    display: none;
  }
  &--active {
    pointer-events: auto;
    .v-icon {
      display: flex;
    }
  }

  &__header {
    div {
      // width: 200%;
    }
    span {
      display: inline-block;
      margin: 0 10px 0 10px;
      // width: calc(30% - 10px);
    }
  }
}
</style>
