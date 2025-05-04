export interface IMessageDTO {
    id?: string;
    message?: string;
    sender_id: string;
    images?: string[];
    audio?: string;
    delete_for?: string[];
    is_update: boolean;
}