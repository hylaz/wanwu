<template>
  <div>
    <el-dialog
      :visible.sync="dialogVisible"
      width="50%"
      min-width="400px"
      append-to-body
      :close-on-click-modal="false"
      custom-class="auth-box"
    >
      <div>
        <p class="auth-header">
          <span style="font-weight: bold">{{ $t('login.twoFA.title') }}</span>
        </p>
        <div class="auth-form">
          <el-form ref="form" :model="form" :rules="rules" label-position="top">
            <el-form-item
              v-if="!isUpdatePassword"
              :label="$t('login.twoFA.form.oldPassword')"
              class="auth-form-item"
              prop="passwordOld">
              <img class="auth-icon" src="@/assets/imgs/pwd.png" alt=""/>
              <el-input
                :type="isShowPwdOld ? '' : 'password'"
                class="auth-pwd-input"
                v-model.trim="form.passwordOld"
                :placeholder="$t('common.input.placeholder') + $t('login.twoFA.form.oldPassword')"
              />
              <img
                v-if="!isShowPwdOld" class="pwd-icon" src="@/assets/imgs/showPwd.png" alt=""
                @click="isShowPwdOld = true"/>
              <img
                v-else class="pwd-icon" src="@/assets/imgs/hidePwd.png" alt=""
                @click="isShowPwdOld = false"/>
            </el-form-item>
            <el-form-item
              v-if="!isUpdatePassword"
              :label="$t('login.twoFA.form.newPassword')"
              class="auth-form-item"
              prop="password1">
              <img class="auth-icon" src="@/assets/imgs/pwd.png" alt=""/>
              <el-input
                :type="isShowPwd1 ? '' : 'password'"
                class="auth-pwd-input"
                v-model.trim="form.password1"
                :placeholder="$t('reset.pwd1Placeholder')"
              />
              <img
                v-if="!isShowPwd1" class="pwd-icon" src="@/assets/imgs/showPwd.png" alt=""
                @click="isShowPwd1 = true"/>
              <img
                v-else class="pwd-icon" src="@/assets/imgs/hidePwd.png" alt=""
                @click="isShowPwd1 = false"/>
            </el-form-item>
            <el-form-item
              v-if="!isUpdatePassword"
              :label="$t('reset.action2') + $t('login.twoFA.form.newPassword')"
              class="auth-form-item"
              prop="password2">
              <img class="auth-icon" src="@/assets/imgs/pwd.png" alt=""/>
              <el-input
                :type="isShowPwd2 ? '' : 'password'"
                class="auth-pwd-input"
                v-model.trim="form.password2"
                :placeholder="$t('reset.action2') + $t('login.twoFA.form.newPassword')"/>
              <img
                v-if="!isShowPwd2" class="pwd-icon" src="@/assets/imgs/showPwd.png" alt=""
                @click="isShowPwd2 = true"/>
              <img
                v-else class="pwd-icon" src="@/assets/imgs/hidePwd.png" alt=""
                @click="isShowPwd2 = false"/>
            </el-form-item>
            <el-form-item :label="$t('login.twoFA.form.email')" class="auth-form-item" prop="email">
              <img class="auth-icon" src="@/assets/imgs/user.png" alt=""/>
              <el-input
                v-model.trim="form.email"
                :placeholder="$t('common.input.placeholder') + $t('login.twoFA.form.email')" clearable
              />
            </el-form-item>
            <el-form-item :label="$t('login.twoFA.form.code')" class="auth-form-item" prop="code">
              <img class="auth-icon" src="@/assets/imgs/code.png" alt=""/>
              <el-input
                style="width: calc(100% - 90px)"
                v-model.trim="form.code"
                @keyup.enter.native="addByEnterKey"
                :placeholder="$t('common.input.placeholder') + $t('login.twoFA.form.code')" clearable
              />
              <el-button
                style="height: 32px; width: 80px; margin-left: 10px; vertical-align: middle; padding-left: 8px; padding-top: 8px"
                @click="requestEmailCode({email: form.email})"
                :disabled="isCooldown"
              >
                {{ isCooldown ? `${cooldownTime}s` : $t('login.twoFA.action') + $t('login.twoFA.form.code') }}
              </el-button>
            </el-form-item>
          </el-form>
          <div class="auth-bt">
            <p class="primary-bt" :style="`background: ${commonInfo.login.loginButtonColor} !important`"
               @click="doLogin">
              {{ $t('login.twoFA.button') }}
            </p>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import {mapState, mapActions} from 'vuex'
import {login2FA2Code, reset} from "@/api/user"
import {urlEncrypt} from "@/utils/crypto";

export default {
  data() {
    let checkPassword2 = (rule, value, callback) => {
      if (this.form.password1 !== this.form.password2) callback(new Error(this.$t('resetPwd.differError')))
      callback()
    }
    let checkPassword1 = (rule, value, callback) => {
      let reg = /^(?=.*[a-zA-Z])(?=.*\d)(?=.*[~!@#$%^&*()_+`\-={}:";'<>?,./]).{8,20}$/
      if (!reg.test(value)) {
        callback(new Error(this.$t('resetPwd.pwdError')))
      } else {
        return callback()
      }
    }
    return {
      dialogVisible: false,
      form: {
        email: '',
        code: '',
        passwordOld: '',
        password1: '',
        password2: ''
      },
      rules: {
        email: [
          {required: true, message: this.$t('common.input.placeholder'), trigger: 'blur'},
          {
            pattern: /^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(.[a-zA-Z0-9_-]+)+$/,
            message: this.$t('common.hint.emailError'),
            trigger: "blur"
          }
        ],
        code: [
          {required: true, message: this.$t('common.input.placeholder'), trigger: 'blur'}
        ],
        passwordOld: [
          {required: true, message: this.$t('common.input.placeholder'), trigger: 'blur'},
        ],
        password1: [
          {required: true, message: this.$t('common.input.placeholder'), trigger: 'blur'},
          {validator: checkPassword1, trigger: "blur"}
        ],
        password2: [
          {required: true, message: this.$t('common.input.placeholder'), trigger: 'blur'},
          {validator: checkPassword1, trigger: "blur"},
          {validator: checkPassword2, trigger: "blur"},
        ],
      },
      isCooldown: false,
      cooldownTime: 60,
      cooldownTimer: '',
      codeSentMessage: '',
      isShowPwdOld: false,
      isShowPwd1: false,
      isShowPwd2: false,
      isEmailCheck: false,
      isUpdatePassword: false,
      codeData: {
        key: '',
        b64: ''
      },
      basePath: this.$basePath,
      params: {
        client_id: '',
        redirect_uri: '',
        scope: '',
        response_type: '',
        state: '',
        client_name: ''
      }
    }
  },
  computed: {
    ...mapState('login', ['commonInfo'])
  },
  methods: {
    ...mapActions('user', ['LoginIn2FA2']),
    showDialog(isEmailCheck, isUpdatePassword, params) {
      this.$store.commit('user/setIs2FA', true)
      this.dialogVisible = true
      this.isEmailCheck = isEmailCheck
      this.isUpdatePassword = isUpdatePassword
      this.params = params
    },
    addByEnterKey(e) {
      if (e.keyCode === 13) {
        this.doLogin()
      }
    },
    doLogin() {
      this.$refs.form.validate(async (valid) => {
        if (!valid) return

        const data = {
          email: this.form.email,
          newPassword: urlEncrypt(this.form.password1),
          oldPassword: urlEncrypt(this.form.passwordOld),
          code: this.form.code
        }

        if (this.isUpdatePassword) {
          delete data.newPassword
          delete data.oldPassword
        }

        await this.LoginIn2FA2({loginInfo: data, params: this.params})
      })

    },
    requestEmailCode(data) {
      this.$refs.form.validateField(['email'], err => {
        if (err) return
        this.codeSentMessage = this.$t('common.hint.codeSent')
        this.isCooldown = true
        this.cooldownTimer = setInterval(() => {
          if (this.cooldownTime > 1) {
            this.cooldownTime--
          } else {
            this.isCooldown = false
            this.cooldownTime = 60
            clearInterval(this.cooldownTimer)
          }
        }, 1000)
        login2FA2Code(data)
      })
    }
  },
  beforeDestroy() {
    clearInterval(this.cooldownTimer)
    this.codeSentMessage = ''
  }
}
</script>

<style lang="scss" scoped>
@import "@/style/auth.scss";

/deep/ .auth-box {
  min-width: 400px;
  background: rgba(244, 247, 255, 0.7);
  border-radius: 4px;
  backdrop-filter: blur(10px);
}
</style>
