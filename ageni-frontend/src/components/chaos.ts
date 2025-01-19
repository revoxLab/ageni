import {ref} from 'vue'

export const mobileMenuOpen = ref(!localStorage.getItem('token'))
