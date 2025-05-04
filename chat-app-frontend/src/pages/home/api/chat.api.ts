import { ref } from 'vue';
import { IMessage } from '@/entities/message.entity';
import { IMessageDTO } from '@/entities/messageDTO.entity';
import { IUser } from '@/entities/user.entity';

const parseMessageData = (message: { id: string; message: string; images: string[]; audio: string; user_from: string; deleted_for: {[key: string]: number}; created_at: number; updated_at: number; }) => {
    const newMessage: IMessage = {
        id: message.id,
        text_content: message.message,
        image_content: message.images ? message.images?.map(image => `${import.meta.env.VITE_DOMAIN_HTTPS}/api/image?id=${image}`) : [],
        audio_content: message.audio,
        sender: message.user_from,
        deleted_for: message.deleted_for,
        created_at: message.created_at,
        updated_at: message.updated_at,
    };
    return newMessage;
};

export function useChatSocket(currentUserId: string) {
    let websocket: WebSocket | null = null;
    const messages = ref<IMessage[]>([]);
    const selectedUser = ref<IUser | null>(null);

    // Handle incoming WebSocket messages
    const initWebSocketHandlers = (): Promise<void> => {
        return new Promise((resolve, reject) => {
            websocket = new WebSocket(`${import.meta.env.VITE_DOMAIN_WSS}/chat/ws`);

            websocket.onopen = () => {
                console.log('WebSocket connection opened');
                resolve();  // Resolve the promise when the WebSocket is open
            };

            websocket.onmessage = (event) => {
                const message = JSON.parse(event.data);

                const parsedMessage = parseMessageData(message);

                let candidate_index: number;
                if (parsedMessage.deleted_for) {
                    if (Object.keys(parsedMessage.deleted_for).some(key => key === currentUserId)) {
                        messages.value = messages.value.filter(element => element.id !== parsedMessage.id);
                        return; 
                    }
                }

                candidate_index = messages.value.findIndex(element => element.id === parsedMessage.id)
                if (candidate_index !== -1) {
                    messages.value[candidate_index] = parsedMessage;
                    return
                }
                
                messages.value.push(parsedMessage);
            };

            websocket.onclose = () => {
                console.log('WebSocket connection closed');
            };

            websocket.onerror = (error) => {
                console.error('WebSocket error:', error);
                reject(error);  // Reject the promise in case of error
            };
        });
    };

    // Send message to the server
    const sendMessage = (newMessage?: string, imageCodes?: string[], audioCode?: string) => {
        if (websocket?.readyState === WebSocket.OPEN) {
            const messageDTO: IMessageDTO = {
                message: newMessage,
                sender_id: currentUserId,
                images: imageCodes,
                audio: audioCode,
                is_update: false,
            };
            websocket.send(JSON.stringify(messageDTO));
        }
    };

    const updateMessage = (messageId: string, updatedContent: string) => {
        if (websocket?.readyState === WebSocket.OPEN) {
            const messageDTO: IMessageDTO = {
                id: messageId,
                message: updatedContent,
                sender_id: "",
                is_update: true,
            };
            websocket.send(JSON.stringify(messageDTO));
        }
    };

    const deleteMessage = (messageId: string, deleteForBoth: boolean, userId?:string) => {
        let delete_for: string[] = [currentUserId];
        if (deleteForBoth) {
            delete_for.push(userId!)
        } 

        if (websocket?.readyState === WebSocket.OPEN) {
            const messageDTO: IMessageDTO = {
                id: messageId,
                sender_id: "",
                delete_for: delete_for,
                is_update: false,
            };
            websocket.send(JSON.stringify(messageDTO));
        }
    }

    // Select a user to start the conversation
    const selectUser = async (token:string, user: IUser) => {
        selectedUser.value = user;
        messages.value = [];  // Clear the current messages
        const result = await fetch(`${import.meta.env.VITE_DOMAIN_HTTPS}/api/chat/init?user1_id=${currentUserId}&user2_id=${user.id}`,
            {
                method: 'GET',
                headers: {'Authorization': `Bearer ${token}`},
            }
        ).then(data => data.json());

        if (result.chat_messages) {
            const oldMessages: IMessage[] = result.chat_messages.map((msg: any) => parseMessageData(msg));
            messages.value.push(...oldMessages);
        }

        // Wait for WebSocket connection to open, then send the chat initiation message
        await initWebSocketHandlers().then(() => {
            if (websocket?.readyState === WebSocket.OPEN && result.chat_id) {
                websocket.send(result.chat_id);
            }
        }).catch(error => {
            console.error('Error initializing WebSocket:', error);
        });
    };

    // Close WebSocket connection
    const closeConnection = () => {
        websocket?.close();
    };

    // Initialize WebSocket handlers when called
    initWebSocketHandlers();

    return {
        messages,
        selectedUser,
        sendMessage,
        updateMessage,
        deleteMessage,
        selectUser,
        closeConnection,
    };
}
