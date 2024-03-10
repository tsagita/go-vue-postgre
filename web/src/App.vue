<template>
  <b-container>
    <b-row class="mt-3">
      <b-col lg="12" class="my-1">
        <b-form-group
          label-for="filter-input"
          label-cols-sm="3"
          label-align-sm="right"
          label-size="sm"
          class="mb-0"
        >
          <b-input-group size="sm">
            <b-form-input
              id="filter-input"
              v-model="filter"
              type="search"
              placeholder="Search by order number or product name"
            ></b-form-input>
          </b-input-group>
        </b-form-group>
      </b-col>
    </b-row>

    <b-row class="mt-3">
      <b-col lg="12">
        <label for="example-datepicker">Choose a date</label>
      </b-col>
    </b-row>

    <b-row class="mt-1">
      <b-col lg="3">
        <b-form-datepicker
          id="start-date"
          v-model="startDate"
          placeholder="Start date"
        ></b-form-datepicker>
      </b-col>
      <b-col lg="3">
        <b-form-datepicker
          id="end-date"
          v-model="endDate"
          placeholder="End date"
        ></b-form-datepicker>
      </b-col>
    </b-row>

    <b-table
      class="mt-5"
      head-variant="dark"
      :items="fetchData"
      :fields="fields"
      :current-page="currentPage"
      :per-page="perPage"
      :filter="filter"
      :startDate="startDate"
      :endDate="endDate"
      :filter-included-fields="filterOn"
      :sort-by.sync="sortBy"
      :sort-desc.sync="sortDesc"
      :sort-direction="sortDirection"
      label-sort-asc=""
      label-sort-desc=""
      stacked="md"
      show-empty
      small
      @filtered="onFiltered"
    >
    </b-table>
    
    <b-row class="mt-5 justify-content-center">
      <b-col sm="12" md="5">
        <b-pagination
          v-model="currentPage"
          :total-rows="totalRows"
          :per-page="perPage"
          align="fill"
          size="sm"
          class="my-0"
        ></b-pagination>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
import axios from "axios";

const host = "http://localhost:4000";

export default {
  data() {
    return {
      fields: [
        {
          key: "order_name",
          label: "Order Name",
        },
        {
          key: "company_name",
          label: "Customer Company",
        },
        {
          key: "customer_name",
          label: "Customer Name",
        },
        {
          key: "created_at",
          label: "Order Date",
          sortable: true
        },
        {
          key: "delivery_amount",
          label: "Delivery Amount",
        },
        {
          key: "total_amount",
          label: "Total Amount",
        },
      ],
      startDate: '',
      endDate: '',
      totalRows: 90,
      currentPage: 1,
      perPage: 5,
      sortBy: "",
      sortDesc: false,
      sortDirection: "asc",
      filter: null,
      filterOn: [],
    };
  },
  computed: {
    sortOptions() {
      return this.fields
        .filter((f) => f.sortable)
        .map((f) => {
          return { text: f.label, value: f.key };
        });
    },
  },
  methods: {
    fetchData(ctx) {
      console.log(this.endDate)
      const params = [
        { key: 'product', value: ctx.filter },
        { key: 'start_date', value: this.startDate },
        { key: 'end_date', value: this.endDate },
        { key: 'per_page', value: this.perPage },
        { key: 'page', value: ctx.currentPage },
      ]
        .filter(param => param.value && param.value !== undefined)
        .map(param => `${param.key}=${encodeURIComponent(param.value)}`)
        .join('&');

      const url = `${host}/api/order/list?${params}`;
      const promise = axios.get(url);

      return promise.then((data) => {
        return data.data.data || [];
      });
    },
    onFiltered(filteredItems) {
      this.totalRows = filteredItems.length;
      this.currentPage = 1;
    },
  },
};
</script>
