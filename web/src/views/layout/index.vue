<template>
  <div class="layout full-menu" :style="`background: ${bgColor}`">
    <el-container class="outer-container">
      <div class="left-nav" v-if="isShowNav">
        <div class="left-nav-container">
          <!--不展示平台的图标-->
          <div style="padding: 0 15px">
            <div
              style="padding: 10px 0 14px; border-bottom: 1px solid #d9d9d9"
              v-if="homeLogoPath"
            >
              <img
                style="width: 50px; margin-left: -5px"
                :src="basePath + '/user/api' + homeLogoPath"
              />
            </div>
          </div>
          <div class="left-nav-content-wrap">
            <div class="left-nav-content">
              <div
                :class="[
                  'nav-item',
                  { 'is-active': currentNavMenu.key === item.key },
                ]"
                v-for="(item, index) in navList"
                :key="item.key + index"
                @click="clickNavMenu(item)"
                v-if="checkPerm(item.perm)"
              >
                <div v-if="item.key !== 'line'">
                  <div class="left-nav-img-wrap">
                    <img
                      class="left-menu-width left-nav-img"
                      :src="
                        currentNavMenu.key === item.key
                          ? item.imgActive
                          : item.img
                      "
                      alt=""
                    />
                  </div>
                  <div class="nav-menu-name">{{ item.name }}</div>
                </div>
                <div v-if="item.key === 'line'">
                  <div
                    style="padding: 0 18px; height: 0.5px; background: #d9d9d9"
                  ></div>
                </div>
              </div>
            </div>
          </div>
          <!--取消整体的新建展示-->
          <!--<div style="padding: 0 15px">
            <div style="padding: 14px 0 10px; border-top: 1px solid #D9D9D9">
              <img class="total-create" src="@/assets/imgs/totalCreate.png" alt="" @click="showCreateTotalDialog">
              <CreateTotalDialog ref="createTotalDialog" />
            </div>
          </div>-->
          <div class="nav-bottom">
            <!--隐藏文档下载菜单-->
            <!--<div>
              <img class="left-menu-width" src="@/assets/imgs/doc.png" alt="" @click="showDocDownloadDialog" />
              <DocDownloadDialog ref="docDownloadDialog" />
            </div>-->
            <AboutDialog ref="aboutDialog" />
            <div style="margin-top: 15px">
              <el-popover placement="right" width="220" trigger="click">
                <div
                  style="margin-bottom: 6px"
                  class="menu--popover-item"
                  :title="getCurrentOrgName()"
                >
                  <el-select
                    v-model="org.orgId"
                    :placeholder="$t('header.org.placeholder')"
                    filterable
                    class="menu__org_select"
                    v-if="orgList && orgList.length"
                    @change="changeOrg"
                  >
                    <el-option
                      v-for="(item, index) in orgList"
                      :command="index"
                      :key="item.id + index"
                      :label="item.name"
                      :value="item.id"
                    />
                  </el-select>
                </div>
                <div
                  :class="[
                    'menu--popover-wrap',
                    { 'wrap-last': popoverList.length === index + 1 },
                  ]"
                  v-for="(it, index) in popoverList"
                  :key="'popoverList' + index"
                >
                  <div
                    v-if="checkPerm(item.perm)"
                    v-for="item in it"
                    :key="item.name"
                    class="menu--popover-item"
                    @click="menuClick(item)"
                  >
                    <img
                      class="menu--popover-item-img"
                      :src="item.img"
                      alt=""
                    />
                    <el-tooltip
                      v-if="item.isTip"
                      effect="dark"
                      :content="item.tipContent"
                      placement="top-start"
                    >
                      <span
                        style="display: inline-block; width: 150px"
                        class="menu--popover-item-name"
                      >
                        {{ item.name }}
                      </span>
                    </el-tooltip>
                    <span v-if="!item.isTip" class="menu--popover-item-name">
                      {{ item.name }}
                    </span>
                    <img
                      v-if="item.icon"
                      class="menu--popover-item-icon"
                      :src="item.icon"
                      alt=""
                    />
                    <span
                      v-if="item.version"
                      class="menu--popover-item-version"
                    >
                      {{ version || '' }}
                    </span>
                  </div>
                </div>
                <div slot="reference">
                  <img
                    class="left-menu-width"
                    src="@/assets/imgs/account.png"
                    alt=""
                  />
                </div>
              </el-popover>
            </div>
          </div>
        </div>
      </div>
      <div class="left-page-container" v-if="isShowNav"></div>
      <!-- 容器 -->
      <el-container :class="['inner-container']">
        <!--取消整体的菜单展示 isShowMenu 一直为 false-->
        <el-aside
          v-if="isShowMenu && menuList && menuList.length"
          class="full-menu-aside"
        >
          <el-menu
            :default-openeds="defaultOpeneds"
            :default-active="activeIndex"
            :key="menuKey"
            :class="[{ 'el-menu-hasOrg': currentNavMenu.key === 'workspace' }]"
          >
            <!--组织切换-->
            <div
              class="header__org_container"
              v-if="currentNavMenu.key === 'workspace'"
            >
              <div class="header__org_wrapper">
                <img class="head-icon" src="@/assets/imgs/head.png" alt="" />
                <el-select
                  v-model="org.orgId"
                  :placeholder="$t('header.org.placeholder')"
                  filterable
                  class="header__org_select"
                  v-if="orgList && orgList.length"
                  @change="changeOrg"
                >
                  <el-option
                    v-for="(item, index) in orgList"
                    :command="index"
                    :key="item.id + index"
                    :class="org.orgId === item.id ? 'header__org_active' : ''"
                    :label="item.name"
                    :value="item.id"
                  />
                </el-select>
              </div>
            </div>
            <!--菜单渲染-->
            <div v-for="(n, i) in menuList" :key="`${i}ml`">
              <!--有下一级-->
              <el-submenu
                v-if="n.children && checkPerm(n.perm)"
                :index="n.index"
                :class="['edit-popover']"
              >
                <template slot="title">
                  <img
                    class="menu-icon"
                    :src="activeIndex.includes(n.index) ? n.imgActive : n.img"
                    alt=""
                  />
                  <span class="menu-withIcon-title">{{ n.name }}</span>
                </template>
                <div
                  v-for="(m, j) in n.children"
                  v-if="checkPerm(m.perm)"
                  :key="`${j}cl`"
                >
                  <el-submenu
                    v-if="m.children"
                    :index="m.index"
                    :class="['menu-indent', 'edit-popover']"
                  >
                    <template slot="title">{{ m.name }}</template>
                    <div
                      v-for="(p, k) in m.children"
                      :key="`${k}pl`"
                      v-if="checkPerm(p.perm)"
                    >
                      <el-submenu
                        v-if="p.children"
                        :index="p.index"
                        :class="['menu-indent-sub', 'edit-popover']"
                      >
                        <template slot="title">{{ p.name }}</template>
                        <el-menu-item
                          v-for="(item, index) in p.children"
                          :key="`${index}itemEl`"
                          :index="item.index"
                          v-if="checkPerm(item.perm)"
                          @click="menuClick(item)"
                          :class="[
                            'edit-popover',
                            { 'is-active': activeIndex === item.index },
                          ]"
                        >
                          {{ item.name }}
                        </el-menu-item>
                      </el-submenu>
                      <el-menu-item
                        v-else
                        :index="p.index"
                        @click="menuClick(p)"
                        :class="[
                          'edit-popover',
                          { 'is-active': activeIndex === p.index },
                        ]"
                      >
                        {{ p.name }}
                      </el-menu-item>
                    </div>
                  </el-submenu>
                  <el-menu-item
                    v-else
                    :index="m.index"
                    @click="menuClick(m)"
                    :class="[
                      'menu-indent-item',
                      'edit-popover',
                      { 'is-active': activeIndex === m.index },
                    ]"
                  >
                    {{ m.name }}
                  </el-menu-item>
                </div>
              </el-submenu>
              <!--没有下一级-->
              <el-menu-item
                :index="n.index"
                v-if="!n.children && checkPerm(n.perm)"
                @click="menuClick(n)"
                :class="[
                  'edit-popover',
                  { 'is-active': activeIndex === n.index },
                ]"
              >
                <img
                  class="menu-icon"
                  :src="activeIndex === n.index ? n.imgActive : n.img"
                  alt=""
                />
                <span class="menu-withIcon-title">{{ n.name }}</span>
              </el-menu-item>
            </div>
          </el-menu>
        </el-aside>
        <!-- 右侧内容 -->
        <el-main>
          <div class="page-container">
            <div class="right-page-content">
              <router-view></router-view>
              <div id="container" class="qk-container"></div>
            </div>
          </div>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script>
// import { start } from 'qiankun'
import { mapActions, mapGetters } from 'vuex';
import { checkPerm, PERMS } from '@/router/permission';
import { menuList } from './menu';
import { changeLang } from '@/api/user';
import {
  fetchPermFirPath,
  fetchCurrentPathIndex,
  replaceIcon,
  replaceTitle,
  redirectUserInfoPage,
} from '@/utils/util';
import ChangeLang from '@/components/changeLang.vue';
import DocDownloadDialog from '@/components/docDownloadDialog.vue';
import CreateTotalDialog from '@/components/createTotalDialog.vue';
import AboutDialog from '@/components/aboutDialog.vue';
import { DOC_FIRST_KEY } from '@/views/docCenter/constants';

export default {
  name: 'Layout',
  components: { ChangeLang, DocDownloadDialog, CreateTotalDialog, AboutDialog },
  data() {
    return {
      basePath: this.$basePath,
      homeLogoPath: '',
      bgColor: '',
      version: '',
      defaultOpeneds: [],
      orgList: [],
      org: { orgId: '' },
      navList: menuList,
      currentNavMenu: {},
      menuList: [],
      menuKey: 'menu_key',
      activeIndex: '0',
      isShowMenu: false,
      isShowNav: true,
      popoverList: [
        [
          {
            name: this.$t('menu.account'),
            path: '/userInfo',
            img: require('@/assets/imgs/user_icon.svg'),
          },
          {
            name: this.$t('menu.setting'),
            path: '/permission',
            img: require('@/assets/imgs/setting_icon.svg'),
            isTip: true,
            tipContent: this.$t('menu.settingTip'),
            perm: PERMS.PERMISSION,
          },
        ],
        [
          {
            name: this.$t('menu.helpDoc'),
            img: require('@/assets/imgs/helpDoc_icon.svg'),
            icon: require('@/assets/imgs/link_icon.png'),
            redirect: () => {
              // window.open('https://github.com/UnicomAI/wanwu/tree/main/docs/manual')
              window.open(
                window.location.origin +
                  `${this.$basePath}/aibase/docCenter/pages/${DOC_FIRST_KEY}`,
              );
            },
          },
          {
            name: 'Github',
            img: require('@/assets/imgs/github_icon.svg'),
            icon: require('@/assets/imgs/link_icon.png'),
            redirect: () => {
              window.open('https://github.com/UnicomAI/wanwu');
            },
          },
          {
            name: this.$t('menu.about'),
            img: require('@/assets/imgs/about_icon.svg'),
            version: 'version',
            redirect: () => {
              // 不展示关于弹窗
              // this.showAboutDialog()
            },
          },
        ],
        [
          {
            name: this.$t('header.logout'),
            img: require('@/assets/imgs/logout_icon.svg'),
            redirect: () => {
              this.logout();
            },
          },
        ],
      ],
    };
  },
  watch: {
    $route: {
      handler(val) {
        // this.justifyIsShowMenu(val.path)
        this.justifyIsShowNav(val.path);
        this.getMenuList(val.path);
        this.redirectUserInfo();
      },
      // 深度观察监听
      deep: true,
    },
    orgInfo: {
      handler(val) {
        this.orgList = val.orgs || [];
      },
      deep: true,
    },
    commonInfo: {
      handler(val) {
        const { home = {}, tab = {}, about = {} } = val.data || {};
        this.homeLogoPath = home.logo ? home.logo.path : '';
        this.bgColor =
          home.backgroundColor ||
          'linear-gradient(1deg, #FFFFFF 42%, #FFFFFF 42%, #EBEDFE 98%, #EEF0FF 98%)';
        this.version = about.version || '1.0';
        replaceIcon(tab.logo ? tab.logo.path : '');
        replaceTitle(tab.title);
      },
      deep: true,
    },
    permission: {
      handler(val) {
        // 如果没修改过密码，重新向到修改密码
        this.redirectUserInfo();
      },
      deep: true,
    },
  },
  computed: {
    ...mapGetters('user', ['orgInfo', 'userInfo', 'commonInfo', 'permission']),
  },
  async created() {
    // 判断是否展示左侧菜单
    this.justifyIsShowNav(this.$route.path);
    // this.justifyIsShowMenu(this.$route.path)

    // 设置语言
    // await this.setLanguage()

    // 获取菜单
    this.getCurrentMenu();

    // 只有登陆状态下才查询接口，否则会一直刷新
    if (localStorage.getItem('access_cert')) this.getPermissionInfo();

    // 设置组织列表以及当前的组织
    this.orgList = this.orgInfo.orgs || [];
    this.org.orgId = this.userInfo.orgId;

    // 获取平台名称以及 logo 等信息
    this.getCommonInfo();
  },
  /* 保证容器 DIV 在 qiankun start 时一定存在 */
  mounted() {
    /* start() */
  },
  methods: {
    ...mapActions('user', ['LoginOut', 'getPermissionInfo', 'getCommonInfo']),
    checkPerm,
    logout() {
      window.localStorage.removeItem('access_cert');
      window.location.href =
        window.location.origin + this.$basePath + '/aibase/login';
    },
    getCurrentOrgName() {
      const currentOrg =
        this.orgList.filter(item => item.id === this.org.orgId)[0] || {};
      return currentOrg.name;
    },
    redirectUserInfo() {
      redirectUserInfoPage(this.permission.isUpdatePassword, () => {
        return null;
      });
    },
    justifyDocPages(val) {
      const path = `${this.$basePath}/aibase` + val;
      return val && path.includes(`${this.$basePath}/aibase/docCenter/pages`);
    },
    justifyIsShowNav(path) {
      const notShowArr = ['/userInfo', '/permission', '/workflow'];
      let isShowNav = true;
      if (this.justifyDocPages(path)) {
        isShowNav = false;
      } else {
        for (let item of notShowArr) {
          if (item === path) {
            isShowNav = false;
            break;
          }
        }
      }
      this.isShowNav = isShowNav;
    },
    justifyIsShowMenu(path) {
      const notShowArr = ['/workflow', '/agent/test', '/rag/test', '/explore'];
      let isShowMenu = true;
      for (let item of notShowArr) {
        if (item === path) {
          isShowMenu = false;
          break;
        }
      }
      this.isShowMenu = isShowMenu;
    },
    /*showCreateTotalDialog() {
      this.$refs.createTotalDialog.openDialog()
    },*/
    showDocDownloadDialog() {
      this.$refs.docDownloadDialog.openDialog();
    },
    showAboutDialog() {
      this.$refs.aboutDialog.openDialog();
    },
    clickNavMenu(item) {
      this.currentNavMenu = item || {};
      const menus = item.children || [];
      this.menuList = menus;

      if (menus && menus.length) {
        // 切换 nav 菜单跳转有权限的第一个
        const { path } = fetchPermFirPath(menus);
        this.$router.push({ path });
        this.changeMenuIndex(fetchCurrentPathIndex(path, menus));
      } else {
        this.$router.push({ path: item.path });
      }
    },
    async setLanguage() {
      const langCode = localStorage.getItem('locale');
      // 主要解决本地和线上两个 localStorage 语言不同问题，使用用户本地缓存的语言
      if (langCode) await changeLang({ language: langCode });
    },
    menuClick(item) {
      if (item.redirect) {
        item.redirect();
      } else {
        if (item.path) this.$router.push({ path: item.path });
      }
    },
    getCurrentMenu() {
      const { path } = this.$route || {};
      // 获取当前菜单
      this.getMenuList(path);
    },
    getCurrentNav(path) {
      // 获取一级路由, 如果是 appSpace 获取两级
      const pathArray = path.split('/') || [];
      const firstLevelPath =
        pathArray[1] === 'appSpace'
          ? `/${pathArray[1] || ''}/${pathArray[2] || ''}`
          : `/${pathArray[1] || ''}`;

      const currentNav = menuList.find(item =>
        JSON.stringify(item).includes(firstLevelPath),
      );
      return currentNav || {};
    },
    getMenuList(path) {
      const currentNavMenu = this.getCurrentNav(path);
      this.currentNavMenu = currentNavMenu;
      // 获取当前菜单列表
      const menus = currentNavMenu.children || [];
      if (!menus.length) return;

      this.menuList = menus;
      this.defaultOpeneds = menus.map(item => item.index);

      // 给当前 activeIndex 赋值
      this.changeMenuIndex(fetchCurrentPathIndex(path, menus));
    },
    changeMenuIndex(index) {
      this.activeIndex = index;
    },
    async changeOrg(orgId) {
      this.$store.state.user.userInfo.orgId = orgId;
      // 切换组织更新权限，跳转有权限的页面；如果是用模型跳转用模型，其他跳转模型开发平台
      await this.getPermissionInfo();

      // 更新 storage 用户信息中组织 id
      const info = JSON.parse(localStorage.getItem('access_cert'));
      info.user.userInfo.orgId = orgId;
      localStorage.setItem('access_cert', JSON.stringify(info));

      const { path } = fetchPermFirPath();
      // 如果当前页面 path 与第一个有权限的 path 相同，需要刷新页面以确保数据为新切换组织的
      if (path === this.$route.path) {
        location.reload();
        return;
      }
      // 切换组织, 根据当前路径有权限的第一个路径找到对应的 menu
      this.getMenuList(path);
      this.menuClick({ path });
    },
  },
};
</script>

<style lang="scss" scoped>
.disabled {
  cursor: not-allowed !important;
}
.full-menu.layout {
  height: 100%;
  /*background: linear-gradient(1deg, #FFFFFF 42%, #FFFFFF 42%, #EBEDFE 98%, #EEF0FF 98%);*/
  /*min-height: 660px;*/
  .outer-container {
    height: 100%;
    .left-page-container {
      //position: relative;
      width: 80px;
      height: 100%;
    }
    .left-nav {
      width: 75px;
      text-align: center;
      padding: 0.5% 0 8px 0;
      position: fixed;
      height: calc(100% - 16px);
      overflow: auto;
      background: #f7f7fc;
      border-radius: 8px;
      box-shadow: 0 1px 4px 0 rgba(0, 0, 0, 0.15);
      margin: 8px 6px;
      z-index: 20;
      .left-nav-container {
        position: relative;
        min-height: 650px;
        height: 100%;
      }
      .left-nav-content-wrap {
        /*display: flex;
        flex-direction: column;
        justify-content: center;*/
      }
      .left-nav-content {
        padding: 6px 5px;
        display: flex;
        flex-direction: column;
        height: auto;
        align-items: center;
        justify-content: space-around;
      }
      .total-create {
        width: 24px;
        cursor: pointer;
      }
      .left-menu-width {
        width: 20px;
        height: 20px;
        object-fit: contain;
      }
      .nav-item {
        /*margin: calc((100vh - 560px) / 28) 0;*/
        margin: 0.6vh 0;
        color: $inactivate_color;
        font-weight: bold;
        cursor: pointer;
        border-radius: 8px;
        .nav-menu-name {
          font-size: 11px;
          margin-top: 3px;
        }
      }
      //.nav-item:hover,
      .nav-item.is-active {
        color: $color;
        .left-nav-img {
          width: 100%;
          height: 100%;
          padding: 8px;
        }
        .left-nav-img-wrap {
          width: 36px;
          height: 36px;
          display: inline-block;
          border-radius: 50%;
          background: #fff;
          box-shadow: 0 2px 8px 0 rgba(0, 0, 0, 0.15);
        }
      }
      .nav-bottom {
        position: absolute;
        bottom: 0;
        width: 100%;
        text-align: center;
        padding-bottom: 10px;
        img {
          cursor: pointer;
        }
      }
    }
    /*element ui 样式重写*/
    .inner-container {
      width: calc(100% - 80px);
      height: 100%;
      border-radius: 10px;
      // border: 1px solid #e6e6e6;
      /*box-shadow: 0px 1px 4px 0px rgba(0, 0, 0, 0.15);*/
      .el-aside.full-menu-aside {
        height: 100%;
        width: 220px !important;
        background-color: rgba(255, 255, 255, 0);
        border-radius: 10px 0 0 10px;
        overflow-y: auto;
        overflow-x: auto;
        .el-menu {
          min-height: 100%;
          width: auto;
          overflow-x: auto;
          overflow-y: hidden;
          .menu-indent /deep/ .el-submenu__title,
          .menu-indent-item {
            padding-left: 49px !important;
          }
          .menu-indent-sub /deep/ .el-submenu__title {
            padding-left: 60px !important;
          }
          .menu-icon {
            width: 16px;
            margin-right: 10px;
          }
          .menu-withIcon-title {
            display: inline-block;
          }
        }
      }
      .el-main {
        position: relative;
        margin: 0 !important;
        padding: 0 !important;
        width: 100%;
        height: 100vh;
        overflow: auto;
        /*background: linear-gradient(1deg, #FFFFFF 42%, #FFFFFF 42%, #EBEDFE 98%, #EEF0FF 98%);
        border-radius: 8px 8px 8px 8px;
        border-left: 0.5px solid #e6e6e6;*/
        .page-container {
          height: 100%;
          overflow: auto;
          padding-right: 2px;
          .right-page-content {
            min-width: 1250px;
            height: 100%;
          }
        }
      }
      /deep/ .el-menu-item {
        color: $color_title;
      }
      /deep/ .el-submenu__title,
      /deep/ .el-menu-item span,
      /deep/ .el-submenu__title span {
        font-size: 14px !important;
      }
      /deep/ .el-menu-item.is-active,
      /deep/ .el-menu-item:focus {
        background-color: $color_opacity !important;
      }
      /deep/ .el-menu-item.is-active,
      /deep/ .el-submenu.is-active {
        .el-submenu__title:hover {
          background-color: $color_opacity !important;
        }
      }
      /*/deep/ .el-submenu.is-active {
        .el-submenu__title:hover {
          background-color: rgba(255, 255, 255, 0) !important;
        }
      }*/
      /deep/ .el-submenu__title {
        span {
          font-size: 14px !important;
        }
      }
      /deep/ .el-submenu.is-active .el-submenu__title {
        border-bottom-color: $color !important;
      }
      /deep/ .el-menu .el-submenu__title,
      /deep/ .el-menu .el-menu-item {
        height: 40px;
        line-height: 40px;
        border-radius: 6px;
        margin: 10px 20px;
        min-width: auto;
      }
      /deep/ .el-menu {
        border: none;
      }
    }
    .inner-container.is-use-model {
      margin-top: 0;
      height: 100%;
    }
  }
  .outer-container /deep/ {
    .el-submenu.is-active,
    .el-submenu.is-active > .el-submenu__title,
    .el-submenu.is-active > .el-submenu__title i:first-child,
    .el-submenu.is-active > .el-submenu__title .el-submenu__icon-arrow {
      color: $color !important;
    }
  }
}

.header__org_container {
  padding: 12px 15px 0 15px;
  .header__org_wrapper {
    padding-bottom: 8px;
    border-bottom: 1px solid #ebebeb;
  }
  .head-icon {
    width: 26px;
    margin: 0 0 0 10px;
    padding-bottom: 2px;
    display: inline-block;
    vertical-align: bottom;
  }
}

.header__org_active {
  color: $color !important;
}
.header__org_select,
.menu__org_select /deep/ {
  width: calc(100% - 37px);
  .el-input__inner:focus,
  .el-input__inner:hover,
  .el-input.is-focus .el-input__inner {
    border-color: #fff !important; // #dcdfe6
  }
  .el-input__inner {
    background-color: rgba(255, 255, 255, 0);
    border: 1px solid #fff;
    color: $color_title;
    font-weight: bold;
    padding-left: 10px;
  }
  .el-input__inner::placeholder {
    color: rgba(18, 18, 18, 0.7);
  }
  .el-input {
    .el-select__caret {
      color: #aaa;
      font-size: 15px;
    }
  }
}
.menu__org_select /deep/ {
  width: 190px;
  .el-input__inner {
    background-color: rgba(255, 255, 255, 0);
    border: none !important;
    color: $color_title !important;
    font-weight: normal;
    padding-left: 0 !important;
    margin-left: 0 !important;
  }
}
.menu--popover-wrap {
  border-top: 1px solid #ebebeb;
  padding: 4px 0 6px 0;
}
.menu--popover-wrap.wrap-last {
  padding-bottom: 0;
}
.menu--popover-item {
  font-size: 13px;
  color: $color_title;
  height: 34px;
  line-height: 34px;
  cursor: pointer;
  border-radius: 4px;
  padding: 0 8px;
  .menu--popover-item-img {
    height: 16px;
    display: inline-block;
    vertical-align: middle;
    margin-right: 5px;
  }
  .menu--popover-item-name {
    font-size: 13px;
    color: $color_title;
    display: inline-block;
    vertical-align: middle;
  }
  .menu--popover-item-icon {
    width: 16px;
    float: right;
    margin-top: 13px;
  }
  .menu--popover-item-version {
    font-size: 13px;
    float: right;
  }
  .menu--popover-item-version:after {
    display: inline-block;
    content: '';
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: #f59a23;
    margin-bottom: 2px;
  }
}
.menu--popover-item:hover /deep/ {
  background: #f5f7fa !important;
  .el-input .el-input__inner {
    border: none !important;
  }
}
</style>
