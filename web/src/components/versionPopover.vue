<template>
  <div>
    <el-popover
      ref="versionPopover"
      placement="bottom"
      width="400"
      trigger="hover"
      popper-class="version-popover"
      @hide="onPopoverHide"
    >
      <VersionTimeLine
        ref="versionTimeLine"
        :appId="appId"
        :appType="appType"
        @reloadData="reloadData"
      />

      <i slot="reference" :class="iconClass" :style="iconStyle" />
    </el-popover>
  </div>
</template>

<script>
import VersionTimeLine from '@/components/versionTimeLine';
export default {
  name: 'VersionPopover',
  components: {
    VersionTimeLine,
  },
  props: {
    appId: {
      type: String,
      required: true,
      default: () => '',
    },
    appType: {
      type: String,
      required: true,
      default: () => '',
    },
    iconStyle: {
      type: Object,
      default: () => ({
        margin: '13px 12px',
        fontSize: '30px',
        color: '#5983ff',
        cursor: 'pointer',
      }),
    },
    iconClass: {
      type: String,
      default: 'el-icon-time',
    },
  },
  methods: {
    reloadData() {
      this.$emit('reloadData');
    },
    onPopoverHide() {
      this.$nextTick(() => {
        const popover = document.querySelector('.version-popover');
        if (popover && popover.contains(document.activeElement)) {
          document.activeElement.blur();
        }
      });
    },
  },
};
</script>
