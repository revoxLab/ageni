<script lang="ts" setup>
import {CustomButton, Icon, Image, themeMode} from '@/components'
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
import {computed, ref} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import AgentPage from './agent/index.vue'
import PluginPage from './plugin/index.vue'
import SpacePage from './space/index.vue'

const route = useRoute()
const {push} = useRouter()
const showMenu = ref(true)
const logged = storage.get('token')
const pathname = ['agent', 'plugin', 'space']
const pathnameWithoutLogin: string[] = ['agent']
const active = computed(
  () =>
    (logged && pathname.find((item) => route.path.startsWith(`/${item}`))) ||
    pathnameWithoutLogin.find((item) => route.path.startsWith(`/${item}`)) ||
    'agent'
)
const nameDict: Record<string, string> = {
  agent: 'Agent Marketplace',
  plugin: 'Plugin Marketplace',
  space: 'Personal Space',
}
</script>

<template>
  <div class="fbh fbjsb w100v">
    <div v-if="showMenu" class="menu fbv pr">
      <Icon
        :size="16"
        name="menu"
        class="show-menu-icon pa"
        @click="showMenu = false"
      />
      <div class="center">
        <Icon name="logo" @click="jump('https://revox.ai/')" :size="32" />
      </div>
      <div class="center">
        <Image class="pt16" name="title" :width="160" :height="12" />
      </div>
      <div class="fbv fbac g8 pt36 oa">
        <template v-for="name in pathname">
          <div
            :class="[
              'fbh fbac g8 menu-item py12 px14 br8 hand',
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
        </template>
      </div>
      <div class="divider my12"></div>
      <div class="fb1"></div>
      <div class="divider w100p"></div>
      <ElDropdown
        v-if="logged"
        trigger="click"
        placement="top"
        class="dropdown"
      >
        <div class="w100p fbh fbac g8 m20 hand">
          <Icon :name="getAvatar()" :size="32" class="br100" />
          <p class="f14 bold">{{ shortAddress(storage.get('address')) }}</p>
          <Icon name="arrow-down1" class="fb1" :size="12" />
        </div>
        <template #dropdown>
          <div class="fbv f14">
            <div
              class="fbh fbjsb fbac p16 g100 hand dropdown-item"
              @click="logout"
            >
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
    <div v-else>
      <Icon
        :size="16"
        name="menu1"
        class="hide-menu-icon p29"
        @click="showMenu = true"
      />
    </div>
    <div class="w100p h100v fbv oa">
      <template v-if="logged">
        <SpacePage v-if="active === 'space'" />
        <AgentPage v-if="active === 'agent'" />
        <PluginPage v-if="active === 'plugin'" />
      </template>
      <template v-else>
        <AgentPage v-if="active === 'agent'" />
      </template>
    </div>
  </div>
</template>

<style lang="less" scoped>
.menu {
  width: 240px;
  height: 100vh;
  flex-shrink: 0;
  padding-top: 32px;
  box-sizing: border-box;
  border-right: 1px solid var(--secondary-border-color);
}
.hide-menu-icon {
  border: 1px solid var(--secondary-border-color);
}
.show-menu-icon {
  top: 40px;
  right: 16px;
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
.divider {
  box-sizing: border-box;
}
</style>
