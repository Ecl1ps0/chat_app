export const getAvailableUsers = async (token: string) => {
    const result = await fetch(`${import.meta.env.VITE_DOMAIN_HTTPS}/api/user/available-users`,
        {
            headers: {Authorization: `Bearer ${token}`}
        }
    ).then(data => data.json());
        
    return result;
}

export const updateUser = async (token: string, formData: FormData) => {
    const result = await fetch(`${import.meta.env.VITE_DOMAIN_HTTPS}/api/update/user?userId=${formData.get("id")?.toString()}`,
        {
            method: 'PUT',
            headers: {'Authorization': `Bearer ${token}`},
            body: formData
        }
    ).then(data => data.json());

    return result;
}