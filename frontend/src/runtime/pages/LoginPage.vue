<template>
    <div class="login-wrapper">
        <div class="login-card">
            <h1>📦 Mini Box</h1>
        </div>
        <p class="subtitle">Entre com suas credenciais</p>
        <form @submit.prevent="handleLogin">
            <label>
                Usuário
                <input v-model="username" type="text" placeholder="admin" autofocus />
            </label>
            <label>
                Senha
                <input v-model="password" type="password" placeholder="admin" />
            </label>
            <p v-if="error" class="error-msg">{{ error }}</p>
            <button type="submit" :disabled="loading">
                {{ loading ? 'Entrando...' : 'Entrar' }}
            </button>
        </form>
    </div>
</template>

<script setup lang="ts">

import { ref } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { authService } from '../services/authService';

const router = useRouter()
const route = useRoute()
const username = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

async function handleLogin() {
    error.value = ''
    loading.value = true

    try {
        await authService.login(username.value, password.value)

        const redirectUrl = (route.query.redirectUrl as string) || '/files'
        router.push(redirectUrl)
    } catch {
        error.value = 'Usuário ou senha incorretos'
    } finally {
        loading.value = false
    }
}

</script>