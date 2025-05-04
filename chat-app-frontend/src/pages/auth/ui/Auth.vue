<script setup lang="ts">
import Button from '@/components/ui/button/Button.vue'
import {
    Card,
    CardContent,
    CardHeader,
    CardFooter,
    CardTitle,
} from '@/components/ui/card'
import Label from '@/components/ui/label/Label.vue'
import { toTypedSchema } from '@vee-validate/zod'
import { Form, Field, ErrorMessage } from 'vee-validate'
import { z } from 'zod'
import { handleAuth } from '../api/auth.api'
import { ref } from 'vue'
import { useToast } from '@/components/ui/toast'
import { useAuth } from '@/shared/stores/auth.store'
import { router } from '@/pages/router/Router'

const isRegistration = ref(true)
const { toast } = useToast()
const { setToken } = useAuth()

const passwordValidation = new RegExp(
    /^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[#?!@$%^&*-])/
)

const formSchema = toTypedSchema(
    z.object({
        username: z
            .string()
            .min(5, { message: 'Username must contain at least 5 characters' }),
        password: z
            .string()
            .min(8, 'Password must contain at least 8 characters')
            .regex(passwordValidation, {
                message:
                    'Password must contain at least one uppercase letter, one lowercase letter, one number and one special character',
            }),
    })
)

const onSubmit = async (values: any) => {
    try {
        const result = await handleAuth(values, isRegistration.value)
        const data = await result.json()

        if (result && data.access_token) {
            setToken(data.access_token)

            router.push('/')
            return
        }

        if (isRegistration.value && result.status === 200) {
            toast({
                title: 'Registration Successful',
                description: 'You have successfully registered!',
                duration: 5000,
            })

            isRegistration.value = false
            return
        } else {
            toast({
                title: 'User Exception',
                description: 'User with such username is already exist!',
                duration: 5000,
            })
        }
    } catch (e) {
        toast({
            title: (e as Error).name,
            description: `Token handling error: ${(e as Error).message}`,
            duration: 5000,
        })
    }
}
</script>

<template>
    <div class="flex items-center justify-center min-h-screen">
        <Card class="w-[350px]">
            <CardHeader>
                <CardTitle>
                    <template v-if="isRegistration"> Registration </template>
                    <template v-else> Login </template>
                </CardTitle>
            </CardHeader>
            <CardContent>
                <Form :validation-schema="formSchema" @submit="onSubmit">
                    <div class="grid items-center w-full gap-4">
                        <div class="flex flex-col space-y-1.5">
                            <Label for="username">Username</Label>
                            <Field
                                id="username"
                                name="username"
                                type="text"
                                class="border rounded-md h-8"
                            />
                            <ErrorMessage
                                name="username"
                                class="text-red-500"
                            />
                        </div>

                        <div class="flex flex-col space-y-1.5">
                            <Label for="password">Password</Label>
                            <Field
                                id="password"
                                name="password"
                                type="password"
                                class="border rounded-md h-8"
                            />
                            <ErrorMessage
                                name="password"
                                class="text-red-500"
                            />
                        </div>

                        <Button class="w-full" type="submit">Submit</Button>
                    </div>
                </Form>
            </CardContent>
            <CardFooter>
                <Button
                    type="button"
                    variant="outline"
                    @click="isRegistration = !isRegistration"
                >
                    <template v-if="isRegistration"> Login </template>
                    <template v-else> Register </template>
                </Button>
            </CardFooter>
        </Card>
    </div>
</template>
