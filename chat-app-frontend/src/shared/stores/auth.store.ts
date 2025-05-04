import { jwtDecode } from "jwt-decode";
import { defineStore } from "pinia";
import { ref, watch } from "vue";

export const useAuth = defineStore('auth', () => {
    const token = ref<string|null>(localStorage.getItem('auth_token') || null);
    const setToken = (newToken: string) => token.value = newToken;

    const isAuthorized = (): boolean => {
        return !!token.value
    }

    const isExpired = (): boolean => {
        if (!token.value) return true;

        const expTime = jwtDecode(token.value).exp!;
        const currentTime = Math.floor(Date.now() / 1000);
        return expTime < currentTime;
    }

    const logout = (): void => {
        token.value = null;
        localStorage.clear();
    }

    watch(token, (newToken, oldToken) => {
        if (newToken) {
            localStorage.setItem('auth_token', newToken);
        } else {
            localStorage.removeItem('auth_token');
        }

        console.log('Token changed from', oldToken, 'to', newToken);
    });

    return { token, setToken, isAuthorized, isExpired, logout };
}) 