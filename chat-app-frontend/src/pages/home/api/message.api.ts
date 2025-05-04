export const getMessage = async (token: string, messageId: string) => {
    const result = await fetch(`${import.meta.env.VITE_DOMAIN_HTTPS}/api/message?messageId=${messageId}`,
        {
            headers: {Authorization: `Bearer ${token}`}
        }
    ).then(data => data.json())

    return result
}