export const API_URI = import.meta.env.VITE_API_URI;

export const STATUS_OK = "OK";
export const STATUS_ERROR = "Error";
export const PAGE_LIMIT = 15;

export interface BaseResponse {
    status: string;
    error?: string;
};