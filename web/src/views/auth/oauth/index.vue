<template>
  <div class="popup-overlay">
    <div class="auth-popup">
      <div class="popup-header">
        <img src="@/assets/imgs/logo_icon.png" alt="Logo" class="logo"/>
        <span class="title">{{ $t('oauth.popup.title') }}</span>
      </div>
      <div class="popup-content">
        <p class="message">
          {{ clientName + $t('oauth.popup.perm') }}
        </p>
        <ul class="permissions-list">
          <li>
            <i class="icon-dot"></i>
            {{ $t('oauth.popup.detail') }}
          </li>
        </ul>
      </div>
      <div class="popup-footer">
        <el-button type="primary" @click="handleCancel">{{ $t('common.button.cancel') }}</el-button>
        <el-button type="success" @click="handleConfirm">{{ $t('common.button.confirm') }}</el-button>
      </div>
    </div>
  </div>
</template>

<script>
import {OAUTH_API} from "@/utils/requestConstants";
import {store} from "@/store";

export default {
  data() {
    return {
      clientName: '',
      token: store.getters['user/token'],
    }
  },
  mounted() {
    const searchParams = new URLSearchParams(window.location.search);
    this.clientName = searchParams.get('client_name') || '';
  },

  methods: {
    handleCancel() {
      window.open("about:blank", "_top")
    },
    handleConfirm() {
      const searchParams = new URLSearchParams(window.location.search);

      const allowedParams = ['client_id', 'redirect_uri', 'scope', 'response_type', 'state'];
      const cleanParams = new URLSearchParams();

      allowedParams.forEach(param => {
        if (searchParams.has(param)) {
          cleanParams.append(param, searchParams.get(param));
        }
      });

      const authorizeUrl = `${OAUTH_API}/oauth/code/authorize?${cleanParams.toString()}&jwt_token=${encodeURIComponent(this.token)}`
      window.location.href = authorizeUrl
    }
  }
};
</script>

<style scoped>

.popup-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  overflow: auto;
}

.auth-popup {
  max-width: 400px;
  width: 90%;
  padding: 24px;
  border-radius: 8px;
  background-color: #fff;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.popup-header {
  display: flex;
  align-items: center;
  justify-content: left;
  margin-bottom: 20px;
}

.logo {
  width: 160px;
  height: 32px;
  margin-right: 8px;
}

.title {
  font-size: 16px;
  font-weight: bold;
  color: #333;
}

.popup-content {
  margin-bottom: 24px;
}

.message {
  font-size: 14px;
  color: #666;
  margin-bottom: 12px;
}

.permissions-list {
  list-style: none;
  padding-left: 0;
  margin: 0;
}

.permissions-list li {
  font-size: 14px;
  color: #333;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
}

.icon-dot {
  width: 8px;
  height: 8px;
  background-color: #007aff;
  border-radius: 50%;
  margin-right: 8px;
}

.popup-footer {
  display: flex;
  justify-content: space-around;
  margin-top: 24px;
}
</style>