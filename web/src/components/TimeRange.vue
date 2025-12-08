<template>
  <el-date-picker
    v-model="timeRange"
    type="datetimerange"
    :picker-options="pickerOptions"
    :start-placeholder="$t('common.datePicker.startPlaceholder')"
    :end-placeholder="$t('common.datePicker.endPlaceholder')"
    value-format="yyyy-MM-dd HH:mm:ss"
    align="right"
  ></el-date-picker>
</template>
<script>
export default {
  data() {
    return {
      timeRange: [],
      pickerOptions: {
        disabledDate(time) {
          return time.getTime() > Date.now();
        },
        shortcuts: [
          {
            text: this.$t('common.datePicker.week'),
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 7);
              picker.$emit('pick', [start, end]);
            },
          },
          {
            text: this.$t('common.datePicker.month'),
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 30);
              picker.$emit('pick', [start, end]);
            },
          },
          {
            text: this.$t('common.datePicker.threeMonth'),
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 90);
              picker.$emit('pick', [start, end]);
            },
          },
        ],
      },
    };
  },
  methods: {
    setTimeRange(time) {
      this.timeRange = time;
    },
    clearTimeRange() {
      this.timeRange = [];
    },
    getTimeRange() {
      return this.timeRange;
    },
  },
};
</script>
