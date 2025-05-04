import { IUserDTO } from "@/entities/userDTO.entity";

export const handleAuth = async (values: any, isRegistration: boolean) => {
    const user = values as IUserDTO;
    const url = `${import.meta.env.VITE_DOMAIN_HTTPS}/api/auth/${isRegistration ? "sign-up" : "sign-in"}`

    const result = await fetch(url, {
        method: "POST",
        body: JSON.stringify(user),
        headers: {
            "Content-type": "application/json; charset=UTF-8"
        }
    })

    return result
}
 