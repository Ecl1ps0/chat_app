export const uploadAudio = async (token: string, formData: FormData) => {
    const result = await fetch(`${import.meta.env.VITE_DOMAIN_HTTPS}/api/audio/upload`,
        {
            method: 'POST',
            headers: {'Authorization': `Bearer ${token}`},
            body: formData
        }
    ).then(data => data.json());

    return result;
}