import { defineStore } from "pinia";
import { computed, ref } from "vue";

export const useAuthStore = defineStore('auth', () => {
    const accessToken = ref<string | null>(localStorage.getItem('minibox_token'))
    const userContextReady = computed(() => accessToken.value !== null)

    function setToken(token: string | null){
        accessToken.value = token
        if(token) localStorage.setItem('minibox_token', token)
        else localStorage.removeItem('minibox_token')
    }

    return {accessToken, userContextReady, setToken}
})