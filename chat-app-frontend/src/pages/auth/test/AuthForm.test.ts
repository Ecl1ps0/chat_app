import { render, screen, fireEvent, waitFor } from '@testing-library/vue'
import '@testing-library/jest-dom'
import AuthForm from '@/pages/auth/ui/Auth.vue'
import { describe, it, expect, vi } from 'vitest'
import { handleAuth } from '@/pages/auth/api/auth.api'

// Mock dependencies
vi.mock('@/pages/auth/api/auth.api', () => ({
    handleAuth: vi.fn().mockResolvedValue({
        json: () => Promise.resolve({ access_token: 'abc123' }),
        status: 200,
    }),
}))
vi.mock('@/shared/stores/auth.store', () => ({
    useAuth: () => ({
        setToken: vi.fn(),
    }),
}))
vi.mock('@/components/ui/toast', () => ({
    useToast: () => ({
        toast: vi.fn(),
    }),
}))
vi.mock('@/pages/router/Router', () => ({
    router: {
        push: vi.fn(),
    },
}))

describe('AuthForm.vue', () => {
    it('renders registration form by default', () => {
        render(AuthForm)
        expect(screen.getByText(/registration/i)).toBeInTheDocument()
        expect(screen.getByLabelText(/username/i)).toBeInTheDocument()
        expect(screen.getByLabelText(/password/i)).toBeInTheDocument()
    })

    it('calls handleAuth on valid submit', async () => {
        render(AuthForm)

        fireEvent.update(screen.getByLabelText(/username/i), 'validUser')
        fireEvent.update(screen.getByLabelText(/password/i), 'Valid123@')

        fireEvent.click(screen.getByText(/submit/i))

        await waitFor(() => {
            expect(handleAuth).toHaveBeenCalled()
        })
    })

    it('toggles between registration and login', async () => {
        render(AuthForm)
        expect(screen.getByText(/registration/i)).toBeInTheDocument()

        fireEvent.click(screen.getByText(/login/i))
        await waitFor(() => {
            expect(screen.getByText(/login/i)).toBeInTheDocument()
        })

        fireEvent.click(screen.getByText(/register/i))
        await waitFor(() => {
            expect(screen.getByText(/registration/i)).toBeInTheDocument()
        })
    })
})
