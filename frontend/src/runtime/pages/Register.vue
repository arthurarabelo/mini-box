<template>
    <div class="login-wrapper">
        <div class="login-card">
            <h1>📦 Mini Box</h1>
            <p class="subtitle">Informe os seus dados abaixo</p>
            <form @submit.prevent="handleRegister">
                <label>
                    Nome
                    <input v-model="name" type="text" placeholder="Arthur Araujo Rabelo" autofocus />
                </label>
                <label>
                    Email
                    <input v-model="email" type="email" placeholder="arthur@email.com" @blur="emailTouched = true" />
                </label>
                <label>
                    Confirme seu email
                    <input v-model="emailConfirmation" type="email" />
                    <span v-if="emailMatchError" class="error-msg">{{ emailMatchError }}</span>
                </label>
                <label>
                    Usuário
                    <input v-model="username" type="text" placeholder="arthurarabelo123" />
                </label>
                <label>
                    Senha
                    <input v-model="password" type="password" @blur="passwordTouched = true" />
                </label>
                <label>
                    Confirme sua senha
                    <input v-model="passwordConfirmation" type="password" />
                    <span v-if="passwordMatchError" class="error-msg"> {{ passwordMatchError }}</span>
                </label>
                <p v-if="error" class="error-msg">{{ error }}</p>
                <button type="submit" :disabled="loading || !isValid">
                    {{ loading ? 'Fazendo seu cadastro...' : 'Cadastrar' }}
                </button>
                <p class="login-link">
                    <RouterLink to="/login">Faça login</RouterLink> 
                </p>
            </form>
        </div>
    </div>
</template>

<script setup lang="ts">

import { computed, ref } from 'vue';
import { authService } from '../services/authService';
import { useRoute, useRouter } from 'vue-router';

const router = useRouter()
const route = useRoute()
const name = ref('')
const username = ref('')
const email = ref('')
const password = ref('')
const emailConfirmation = ref('')
const passwordConfirmation = ref('')
const error = ref('')
const loading = ref(false)
const emailTouched = ref(false)
const passwordTouched = ref(false)

const isValid = computed(() =>
    name.value.trim() !== '' &&
    username.value.trim() !== '' &&
    email.value.trim() !== '' &&
    password.value.trim() !== ''
)

const emailMatchError = computed(() =>
    emailTouched.value && emailConfirmation.value !== email.value ?
        'Os emails não coincidem' :
        ''
)

const passwordMatchError = computed(() =>
    passwordTouched.value && passwordConfirmation.value !== password.value ?
        'As senhas não coincidem' :
        ''
)

async function handleRegister() {
    error.value = ''

    loading.value = true

    try {
        await authService.register(name.value, username.value, email.value, password.value)

        const redirectUrl = (route.query.redirectUrl as string) || '/login'
        router.push(redirectUrl)
    } catch (e) {
        error.value = e instanceof Error ? e.message : 'Erro ao fazer cadastro'
    } finally {
        loading.value = false
    }
}

</script>