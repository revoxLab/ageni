<script lang="ts" setup>
import {
  CustomButton,
  Icon,
  Image,
  mobileMenuOpen,
  themeMode,
} from '@/components'
import {
  checkLogin,
  getAvatar,
  jump,
  logout,
  shortAddress,
  storage,
} from '@/utils'
import {openWallet} from '@/wallet'
import {ElDropdown} from 'element-plus'
import {computed} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import AgentPage from './agent/index.vue'

const route = useRoute()
const {push} = useRouter()
const logged = storage.get('token')
const pathnameWithoutLogin = ['agent']
const pathname = ['agent']
const menuVisible = computed(
  () => route.fullPath.split('/').filter(Boolean).length <= 1
)
const active = computed(
  () =>
    (logged && pathname.find((item) => route.path.startsWith(`/${item}`))) ||
    pathnameWithoutLogin.find((item) => route.path.startsWith(`/${item}`)) ||
    'agent'
)
const nameDict: Record<string, string> = {
  agent: 'Agent Marketplace',
  plugin: 'Plugin Marketplace',
}
</script>

<template>
  <div class="fbh fbjsb w100v oh">
    <div
      v-if="menuVisible"
      class="menu pa fbv fbac"
      :class="{show: mobileMenuOpen}"
    >
      <Teleport to="#modal">
        <div
          v-if="mobileMenuOpen"
          class="pa modal wh100p"
          @click="mobileMenuOpen = false"
        ></div>
      </Teleport>
      <Icon name="logo" @click="jump('https://revox.ai/')" :size="32" />
      <Image class="pt16" name="title" :width="160" :height="12" />
      <div class="fb1 fbv fbac g8 pt36 oa">
        <div
          v-for="name in pathname"
          :class="[
            'fbh fbac g8 menu-item py12 px14 br8',
            {active: active === name},
          ]"
          @click="
            pathnameWithoutLogin.includes(name)
              ? push(`/${name}`)
              : checkLogin(() => push(`/${name}`))
          "
        >
          <Icon :key="`${name}_${themeMode}`" :name="name" :size="16" />
          <p class="f16">
            {{ nameDict[name] }}
          </p>
        </div>
      </div>
      <div class="divider w100p"></div>
      <ElDropdown
        v-if="logged"
        trigger="click"
        placement="top"
        class="dropdown w100p"
      >
        <div class="w100p fbh fbac g8 my18 px14">
          <Icon :name="getAvatar()" :size="32" class="br100" />
          <p class="f14 bold fb1">{{ shortAddress(storage.get('address')) }}</p>
          <Icon name="arrow-down1" :size="12" />
        </div>
        <template #dropdown>
          <div class="fbv f14">
            <div class="fbh fbjsb fbac p16 g100 dropdown-item" @click="logout">
              <p class="pr16">Log out</p>
              <Icon name="exit1" :size="14" />
            </div>
          </div>
        </template>
      </ElDropdown>
      <CustomButton
        v-else
        :width="120"
        :height="36"
        class="center f16 p18"
        @click="openWallet"
      >
        Sign in
      </CustomButton>
    </div>
    <div class="w100p h100v fbv oa">
      <template v-if="logged">
        <AgentPage v-if="active === 'agent'" />
      </template>
      <template v-else>
        <AgentPage v-if="active === 'agent'" />
      </template>
    </div>
  </div>
</template>

<style lang="less" scoped>
.menu {
  z-index: 2;
  width: 240px;
  height: 100vh;
  padding-top: 32px;
  box-sizing: border-box;
  border-right: 1px solid var(--secondary-border-color);
  background: #1b1b1b;
  transition: margin-left 0.5s;
  margin-left: -100%;
  &.show {
    margin-left: 0;
  }
}
.modal {
  top: 0;
  left: 0;
  z-index: 1;
}
.dropdown {
  color: var(--primary-text-color);
}
.dropdown-item {
  color: var(--primary-text-color);
  &:hover {
    background: var(--secondary-background-color);
  }
}
.menu-item {
  width: 180px;
  opacity: 0.6;
  &:hover {
    opacity: 1;
  }
  &.active {
    background: var(--tertiary-background-color);
    opacity: 1;
  }
}
</style>
