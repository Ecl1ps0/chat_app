export interface IMessage {
    id: string;
    text_content: string;
    image_content: string[];
    audio_content: string;
    sender: string;
    deleted_for: {[key: string]: number}
    created_at: number;
    updated_at: number;
}